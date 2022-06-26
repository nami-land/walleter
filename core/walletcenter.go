package core

import (
	"errors"
	"github.com/neco-fun/wallet-center/internal/comm"
	"github.com/neco-fun/wallet-center/internal/model"
	"github.com/neco-fun/wallet-center/internal/pkg"
	"github.com/neco-fun/wallet-center/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"os"
)

type WalletCommand struct {
	GameClient     comm.GameClient //game client, like Neco Fishing, Neco Land..etc.
	AccountId      uint            // User account id. unique
	PublicAddress  string          // public address, allowed to be null
	AssetType      comm.AssetType  // 0: ERC20 token, 1: erc1155 token.
	ERC20Commands  []ERC20Command
	ERC1155Command ERC1155Command
	BusinessModule string
	ActionType     comm.WalletActionType
	FeeCommands    []ERC20Command // charge fee, if len(FeeCommands) > 0, should be deducted from user's account.
}

type ERC20Command struct {
	Token   comm.ERC20Token
	Value   float64
	Decimal uint
}

type ERC1155Command struct {
	Ids    []uint
	Values []uint
}

type walletCenter struct {
	db *gorm.DB
}

func New(db *gorm.DB, feeCharger OfficialAccount) *walletCenter {
	migration(db)
	return &walletCenter{db: db}
}

func migration(db *gorm.DB) {
	_ = db.Table("t_erc20_token_data").AutoMigrate(model.ERC20TokenWallet{})
	_ = db.Table("t_erc1155_token_data").AutoMigrate(model.ERC1155TokenWallet{})
	_ = db.Table("t_wallet").AutoMigrate(model.Wallet{})
	_ = db.Table("t_erc20_wallet_log").AutoMigrate(model.ERC20WalletLog{})
	_ = db.Table("t_erc1155_wallet_log").AutoMigrate(model.ERC1155WalletLog{})
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func (s *walletCenter) HandleWalletCommand(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	if !command.GameClient.IsSupport() {
		return model.Wallet{}, errors.New("game client is invalid")
	}

	switch command.ActionType {
	case comm.Initialize:
		return initWallet(db, command)
	default:
		return updateWallet(db, command)
	}
}

func initWallet(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	err := db.Transaction(func(tx1 *gorm.DB) error {
		// 1. Insert change logs, including ERC20 logs and ERC1155 Log.
		walletLogService := newWalletLogService()
		erc20WalletLog, err := walletLogService.InsertNewERC20WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		erc115WalletLog, err := walletLogService.InsertNewERC1155WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := pkg.ParseCommandToERC20WalletArray(command)
		erc1155Data := pkg.ParseCommandToERC1155Wallet(command)
		wallet := model.Wallet{
			GameClient:       int(command.GameClient),
			AccountId:        command.AccountId,
			PublicAddress:    command.PublicAddress,
			ERC20TokenData:   erc20DataArray,
			ERC1155TokenData: erc1155Data,
			CheckSign:        "",
		}

		// 3. generator a new check sign
		newCheckSign, err := pkg.NewWalletValidator().GenerateNewSignHash(wallet)
		if err != nil {
			return err
		}
		wallet.CheckSign = newCheckSign

		err = model.WalletDAO.CreateWallet(tx1, wallet)
		if err != nil {
			return err
		}

		// 4. change log statuses
		_, err = walletLogService.UpdateERC20WalletLog(tx1, erc20WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		_, err = walletLogService.UpdateERC1155WalletLog(tx1, erc115WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.GameClient, command.AccountId)
}

func updateWallet(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	switch command.AssetType {
	case comm.ERC20AssetType:
		return handleERC20Command(db, command)
	case comm.ERC1155AssetType:
		return handleERC1155Command(db, command)
	default:
		return model.Wallet{}, errors.New("not support current asset type")
	}
}

func handleERC20Command(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		logService := newWalletLogService()
		validator := pkg.NewWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(db, command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		// 1. Verify that the user's current wallet status is normal
		result, err := validator.ValidateWallet(userWallet)
		if err != nil || result == false {
			return err
		}

		// 2.Insert a log message
		erc20Log, err := logService.InsertNewERC20WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. Whether to charge a fee
		userWallet, err = newFeeChargerService().chargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.UpdateERC20WalletLog(tx, erc20Log, comm.Failed, userWallet)
			return err
		}

		// 4. Make changes to user assets
		switch command.ActionType {
		case comm.Deposit:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalDeposit += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = model.WalletDAO.UpdateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case comm.Withdraw:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalWithdraw += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = model.WalletDAO.UpdateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case comm.Income:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalIncome += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = model.WalletDAO.UpdateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case comm.Spend:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalSpend += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = model.WalletDAO.UpdateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case comm.ChargeFee:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalFee += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = model.WalletDAO.UpdateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		default:
			return errors.New("not support action type")
		}

		// 6. Generate new verification information
		newCheckSign, err := validator.GenerateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign
		err = model.WalletDAO.UpdateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().UpdateERC20WalletLog(tx, erc20Log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.GameClient, command.AccountId)
}

func handleERC1155Command(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		logService := newWalletLogService()
		validator := pkg.NewWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(db, command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		// 1. Verify that the user's current wallet status is normal
		result, err := validator.ValidateWallet(userWallet)
		if err != nil || !result {
			return err
		}

		// 2.Insert a log message
		erc1155Log, err := logService.InsertNewERC1155WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. Whether to charge a fee
		userWallet, err = newFeeChargerService().chargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.UpdateERC1155WalletLog(tx, erc1155Log, comm.Failed, userWallet)
			return err
		}

		// 5. Make changes to user assets
		ids := utils.ConvertStringToUIntArray(userWallet.ERC1155TokenData.Ids)
		values := utils.ConvertStringToUIntArray(userWallet.ERC1155TokenData.Values)
		switch command.ActionType {
		case comm.Deposit, comm.Income:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := utils.GetIndexFromUIntArray(ids, id)
				if i == -1 {
					ids = append(ids, id)
					values = append(values, value)
				} else {
					values[i] = values[i] + value
				}

				userWallet.ERC1155TokenData.Ids = utils.ConvertUintArrayToString(ids, ",")
				userWallet.ERC1155TokenData.Values = utils.ConvertUintArrayToString(values, ",")
				err = model.WalletDAO.UpdateERC1155WalletData(tx, userWallet.ERC1155TokenData)
				if err != nil {
					return err
				}
			}
			break
		case comm.Withdraw, comm.Spend:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := utils.GetIndexFromUIntArray(ids, id)
				if i == -1 {
					return errors.New("insufficient nft balance")
				} else {
					if values[i] < value {
						return errors.New("insufficient nft balance")
					}
					values[i] = values[i] - value
				}

				userWallet.ERC1155TokenData.Ids = utils.ConvertUintArrayToString(ids, ",")
				userWallet.ERC1155TokenData.Values = utils.ConvertUintArrayToString(values, ",")
				err = model.WalletDAO.UpdateERC1155WalletData(tx, userWallet.ERC1155TokenData)
				if err != nil {
					return err
				}
			}
			break
		default:
			return errors.New("not support action type")
		}

		// 6. Generate new verification information
		newCheckSign, err := validator.GenerateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign
		err = model.WalletDAO.UpdateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().UpdateERC1155WalletLog(tx, erc1155Log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.GameClient, command.AccountId)
}

type feeChargerService struct{}

func newFeeChargerService() *feeChargerService {
	return &feeChargerService{}
}

func (*feeChargerService) chargeFee(db *gorm.DB, command WalletCommand, userWallet model.Wallet) (model.Wallet, error) {
	feeChargerWallet, err := model.WalletDAO.GetWallet(
		db, command.GameClient, GetFeeChargerAccount(command.GameClient).AccountId,
	)
	if err != nil {
		return userWallet, err
	}

	for _, fee := range command.FeeCommands {
		if fee.Value <= 0 {
			continue
		}
		index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, fee.Token)
		if index < 0 || userERC20TokenWallet.Balance < fee.Value {
			return userWallet, errors.New("insufficient balance for fee")
		}

		userERC20TokenWallet.Balance -= fee.Value
		userERC20TokenWallet.TotalFee += fee.Value
		userWallet.ERC20TokenData[index] = userERC20TokenWallet

		index, feeChargerERC20TokenWallet := getUserERC20TokenWallet(feeChargerWallet.ERC20TokenData, fee.Token)
		feeChargerERC20TokenWallet.Balance += fee.Value
		feeChargerERC20TokenWallet.TotalFee += fee.Value
		feeChargerWallet.ERC20TokenData[index] = feeChargerERC20TokenWallet
		err = model.WalletDAO.UpdateERC20WalletData(db, feeChargerERC20TokenWallet)
		if err != nil {
			return model.Wallet{}, err
		}
	}

	return userWallet, nil
}

func getUserERC20TokenWallet(tokens []model.ERC20TokenWallet, token comm.ERC20Token) (int, model.ERC20TokenWallet) {
	for index, item := range tokens {
		if item.Token == token.String() {
			return index, item
		}
	}
	return -1, model.ERC20TokenWallet{}
}

type walletLogService struct{}

func newWalletLogService() *walletLogService {
	return &walletLogService{}
}

// InsertNewERC20WalletLog Insert new log of ERC20 changes
func (receiver *walletLogService) InsertNewERC20WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	erc20WalletLog := pkg.ParseCommandToERC20WalletLog(command, currentWallet)
	return model.ERC20WalletLogDAO.InsertERC20WalletLog(db, erc20WalletLog)
}

// UpdateERC20WalletLog Change the status of ERC20Log in batches
func (receiver *walletLogService) UpdateERC20WalletLog(
	db *gorm.DB, log model.ERC20WalletLog, status comm.WalletLogStatus, newWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return model.ERC20WalletLogDAO.UpdateERC20WalletLogStatus(db, log)
}

// InsertNewERC1155WalletLog Insert an ERC1155 asset change log
func (receiver *walletLogService) InsertNewERC1155WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet model.Wallet,
) (model.ERC1155WalletLog, error) {
	erc1155WalletData := pkg.ParseCommandToERC1155WalletLog(command, currentWallet)
	return model.ERC1155WalletLogDAO.InsertERC1155WalletLog(db, erc1155WalletData)
}

// UpdateERC1155WalletLog Change the state of the ERC1155 log
func (receiver *walletLogService) UpdateERC1155WalletLog(
	db *gorm.DB, log model.ERC1155WalletLog, status comm.WalletLogStatus, newWallet model.Wallet,
) (model.ERC1155WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return model.ERC1155WalletLogDAO.UpdateERC1155WalletLogStatus(db, log)
}

type OfficialAccount struct {
	AccountId     uint
	PublicAddress string
}

var necoFishingFeeChargerAccount = OfficialAccount{
	AccountId:     1,
	PublicAddress: "0xa98Ff091a5F6975162AEa4E3862165bCf81aB4Ad",
}

func GetFeeChargerAccount(gameClient comm.GameClient) OfficialAccount {
	if gameClient == comm.NecoFishing {
		return necoFishingFeeChargerAccount
	}
	return necoFishingFeeChargerAccount
}

var initNecoFishingFeeChargerAccountCommand = WalletCommand{
	GameClient:    comm.NecoFishing,
	AccountId:     GetFeeChargerAccount(comm.NecoFishing).AccountId,
	PublicAddress: GetFeeChargerAccount(comm.NecoFishing).PublicAddress,
	AssetType:     comm.Other,
	ERC20Commands: []ERC20Command{
		{
			Token:   comm.NFISH,
			Value:   0,
			Decimal: 18,
		}, {
			Token:   comm.BUSD,
			Value:   0,
			Decimal: 18,
		},
	},
	ERC1155Command: ERC1155Command{},
	BusinessModule: "Initialization",
	ActionType:     comm.Initialize,
	FeeCommands:    []ERC20Command{},
}

var InitializedCommands = []WalletCommand{
	initNecoFishingFeeChargerAccountCommand,
}

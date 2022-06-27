package core

import (
	"errors"
	"os"

	"github.com/neco-fun/wallet-center/internal/comm"
	"github.com/neco-fun/wallet-center/internal/model"
	"github.com/neco-fun/wallet-center/internal/utils"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletCommand struct {
	AccountId      uint           // User account id. unique
	PublicAddress  string         // public address, allowed to be null
	AssetType      comm.AssetType // 0: ERC20 token, 1: erc1155 token.
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

var feeChargerAccount *OfficialAccount

func New(db *gorm.DB, feeCharger OfficialAccount) *walletCenter {
	migration(db)
	feeChargerAccount = &feeCharger
	return &walletCenter{db: db}
}

func (s *walletCenter) InitFeeChargerAccount() (model.Wallet, error) {
	if feeChargerAccount == nil {
		panic("Please assign official fee charge account.")
	}

	command := buildInitializedCommandFromAccount(*feeChargerAccount)
	return s.HandleWalletCommand(s.db, command)
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
		erc20WalletLog, err := walletLogService.insertNewERC20WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		erc115WalletLog, err := walletLogService.insertNewERC1155WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := parseCommandToERC20WalletArray(command)
		erc1155Data := parseCommandToERC1155Wallet(command)
		wallet := model.Wallet{
			AccountId:        command.AccountId,
			PublicAddress:    command.PublicAddress,
			ERC20TokenData:   erc20DataArray,
			ERC1155TokenData: erc1155Data,
			CheckSign:        "",
		}

		// 3. generator a new check sign
		newCheckSign, err := newWalletValidator().generateNewSignHash(wallet)
		if err != nil {
			return err
		}
		wallet.CheckSign = newCheckSign

		err = model.WalletDAO.CreateWallet(tx1, wallet)
		if err != nil {
			return err
		}

		// 4. change log statuses
		_, err = walletLogService.updateERC20WalletLog(tx1, erc20WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		_, err = walletLogService.updateERC1155WalletLog(tx1, erc115WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.AccountId)
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
		validator := newWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(db, command.AccountId)
		if err != nil {
			return err
		}

		// 1. Verify that the user's current wallet status is normal
		result, err := validator.validateWallet(userWallet)
		if err != nil || result == false {
			return err
		}

		// 2.Insert a log message
		erc20Log, err := logService.insertNewERC20WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. Whether to charge a fee
		userWallet, err = newFeeChargerService().chargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.updateERC20WalletLog(tx, erc20Log, comm.Failed, userWallet)
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
		newCheckSign, err := validator.generateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign
		err = model.WalletDAO.UpdateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().updateERC20WalletLog(tx, erc20Log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.AccountId)
}

func handleERC1155Command(db *gorm.DB, command WalletCommand) (model.Wallet, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		logService := newWalletLogService()
		validator := newWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(db, command.AccountId)
		if err != nil {
			return err
		}

		// 1. Verify that the user's current wallet status is normal
		result, err := validator.validateWallet(userWallet)
		if err != nil || !result {
			return err
		}

		// 2.Insert a log message
		erc1155Log, err := logService.insertNewERC1155WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. Whether to charge a fee
		userWallet, err = newFeeChargerService().chargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.updateERC1155WalletLog(tx, erc1155Log, comm.Failed, userWallet)
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
		newCheckSign, err := validator.generateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign
		err = model.WalletDAO.UpdateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().updateERC1155WalletLog(tx, erc1155Log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(db, command.AccountId)
}

type feeChargerService struct{}

func newFeeChargerService() *feeChargerService {
	return &feeChargerService{}
}

func (*feeChargerService) chargeFee(db *gorm.DB, command WalletCommand, userWallet model.Wallet) (model.Wallet, error) {
	feeChargerWallet, err := model.WalletDAO.GetWallet(db, feeChargerAccount.AccountId)
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

// insertNewERC20WalletLog Insert new log of ERC20 changes
func (receiver *walletLogService) insertNewERC20WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	erc20WalletLog := parseCommandToERC20WalletLog(command, currentWallet)
	return model.ERC20WalletLogDAO.InsertERC20WalletLog(db, erc20WalletLog)
}

// updateERC20WalletLog Change the status of ERC20Log in batches
func (receiver *walletLogService) updateERC20WalletLog(
	db *gorm.DB, log model.ERC20WalletLog, status comm.WalletLogStatus, newWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return model.ERC20WalletLogDAO.UpdateERC20WalletLogStatus(db, log)
}

// insertNewERC1155WalletLog Insert an ERC1155 asset change log
func (receiver *walletLogService) insertNewERC1155WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet model.Wallet,
) (model.ERC1155WalletLog, error) {
	erc1155WalletData := parseCommandToERC1155WalletLog(command, currentWallet)
	return model.ERC1155WalletLogDAO.InsertERC1155WalletLog(db, erc1155WalletData)
}

// updateERC1155WalletLog Change the state of the ERC1155 log
func (receiver *walletLogService) updateERC1155WalletLog(
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

func buildInitializedCommandFromAccount(account OfficialAccount) WalletCommand {
	return WalletCommand{
		AccountId:     feeChargerAccount.AccountId,
		PublicAddress: feeChargerAccount.PublicAddress,
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
}

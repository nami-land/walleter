package wallet_center

import (
	"context"
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type WalletCommand struct {
	AccountId      uint64    // User account id. unique
	AssetType      AssetType // 0: ERC20 token, 1: erc1155 token.
	ERC20Commands  []ERC20Command
	ERC1155Command ERC1155Command
	BusinessModule string
	ActionType     WalletActionType
	FeeCommands    []ERC20Command // charge fee, if len(FeeCommands) > 0, should be deducted from user's account.
}

type ERC20Command struct {
	Token   ERC20TokenEnum
	Value   float64
	Decimal uint
}

type ERC1155Command struct {
	Ids    []uint
	Values []uint
}

type WalletCenter struct {
	db *gorm.DB
}

var feeChargerAccountId uint64

func New(db *gorm.DB, chargerAccountId uint64) *WalletCenter {
	migration(db)
	feeChargerAccountId = chargerAccountId
	return &WalletCenter{db: db}
}

func (s *WalletCenter) SetFeeChargerAccount() (Wallet, error) {
	if feeChargerAccountId == 0 {
		panic("Please assign official fee charge account.")
	}

	command := buildInitializedCommandFromAccount(feeChargerAccountId)
	return s.HandleWalletCommand(context.Background(), command)
}

func migration(db *gorm.DB) {
	_ = db.AutoMigrate(ERC20TokenWallet{})
	_ = db.AutoMigrate(ERC1155TokenWallet{})
	_ = db.AutoMigrate(Wallet{})
	_ = db.AutoMigrate(ERC20WalletLog{})
	_ = db.AutoMigrate(ERC1155WalletLog{})
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

func (s *WalletCenter) HandleWalletCommand(ctx context.Context, command WalletCommand) (Wallet, error) {
	switch command.ActionType {
	case Initialize:
		return initWallet(s.db.WithContext(ctx), command)
	default:
		return updateWallet(s.db.WithContext(ctx), command)
	}
}

func initWallet(db *gorm.DB, command WalletCommand) (Wallet, error) {
	err := db.Transaction(func(tx1 *gorm.DB) error {
		// 1. Insert change logs, including ERC20 logs and ERC1155 Log.
		walletLogService := newWalletLogService()
		erc20WalletLog, err := walletLogService.insertNewERC20WalletLog(tx1, command, Wallet{})
		if err != nil {
			return err
		}

		erc115WalletLog, err := walletLogService.insertNewERC1155WalletLog(tx1, command, Wallet{})
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := parseCommandToERC20WalletArray(command)
		erc1155Data := parseCommandToERC1155Wallet(command)
		wallet := Wallet{
			AccountId:        command.AccountId,
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

		err = walletDAO.createWallet(tx1, wallet)
		if err != nil {
			return err
		}

		// 4. change log statuses
		_, err = walletLogService.updateERC20WalletLog(tx1, erc20WalletLog, Done, wallet)
		if err != nil {
			return err
		}

		_, err = walletLogService.updateERC1155WalletLog(tx1, erc115WalletLog, Done, wallet)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return Wallet{}, err
	}
	return walletDAO.getWallet(db, command.AccountId)
}

func updateWallet(db *gorm.DB, command WalletCommand) (Wallet, error) {
	switch command.AssetType {
	case ERC20AssetType:
		return handleERC20Command(db, command)
	case ERC1155AssetType:
		return handleERC1155Command(db, command)
	default:
		return Wallet{}, errors.New("not support current asset type")
	}
}

func handleERC20Command(db *gorm.DB, command WalletCommand) (Wallet, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		logService := newWalletLogService()
		validator := newWalletValidator()

		userWallet, err := walletDAO.getWallet(db, command.AccountId)
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
			_, err = logService.updateERC20WalletLog(tx, erc20Log, Failed, userWallet)
			return err
		}

		// 4. Make changes to user assets
		switch command.ActionType {
		case Deposit:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalDeposit += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = walletDAO.updateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case Withdraw:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalWithdraw += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = walletDAO.updateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case Income:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalIncome += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = walletDAO.updateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case Spend:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalSpend += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = walletDAO.updateERC20WalletData(tx, userERC20TokenWallet)
				if err != nil {
					return err
				}
			}
			break
		case ChargeFee:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalFee += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
				err = walletDAO.updateERC20WalletData(tx, userERC20TokenWallet)
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
		err = walletDAO.updateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().updateERC20WalletLog(tx, erc20Log, Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Wallet{}, err
	}
	return walletDAO.getWallet(db, command.AccountId)
}

func handleERC1155Command(db *gorm.DB, command WalletCommand) (Wallet, error) {
	err := db.Transaction(func(tx *gorm.DB) error {
		logService := newWalletLogService()
		validator := newWalletValidator()

		userWallet, err := walletDAO.getWallet(db, command.AccountId)
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
			_, err = logService.updateERC1155WalletLog(tx, erc1155Log, Failed, userWallet)
			return err
		}

		// 5. Make changes to user assets
		ids := convertStringToUIntArray(userWallet.ERC1155TokenData.Ids)
		values := convertStringToUIntArray(userWallet.ERC1155TokenData.Values)
		switch command.ActionType {
		case Deposit, Income:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := getIndexFromUIntArray(ids, id)
				if i == -1 {
					ids = append(ids, id)
					values = append(values, value)
				} else {
					values[i] = values[i] + value
				}

				userWallet.ERC1155TokenData.Ids = convertUintArrayToString(ids, ",")
				userWallet.ERC1155TokenData.Values = convertUintArrayToString(values, ",")
				err = walletDAO.updateERC1155WalletData(tx, userWallet.ERC1155TokenData)
				if err != nil {
					return err
				}
			}
			break
		case Withdraw, Spend:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := getIndexFromUIntArray(ids, id)
				if i == -1 {
					return errors.New("insufficient nft balance")
				} else {
					if values[i] < value {
						return errors.New("insufficient nft balance")
					}
					values[i] = values[i] - value
				}

				userWallet.ERC1155TokenData.Ids = convertUintArrayToString(ids, ",")
				userWallet.ERC1155TokenData.Values = convertUintArrayToString(values, ",")
				err = walletDAO.updateERC1155WalletData(tx, userWallet.ERC1155TokenData)
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
		err = walletDAO.updateWalletCheckSign(tx, userWallet)
		if err != nil {
			return err
		}

		// 8. Update log information
		_, err = newWalletLogService().updateERC1155WalletLog(tx, erc1155Log, Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return Wallet{}, err
	}
	return walletDAO.getWallet(db, command.AccountId)
}

type feeChargerService struct{}

func newFeeChargerService() *feeChargerService {
	return &feeChargerService{}
}

func (*feeChargerService) chargeFee(db *gorm.DB, command WalletCommand, userWallet Wallet) (Wallet, error) {
	feeChargerWallet, err := walletDAO.getWallet(db, feeChargerAccountId)
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
		err = walletDAO.updateERC20WalletData(db, feeChargerERC20TokenWallet)
		if err != nil {
			return Wallet{}, err
		}
	}

	return userWallet, nil
}

func getUserERC20TokenWallet(tokens []ERC20TokenWallet, token ERC20TokenEnum) (int, ERC20TokenWallet) {
	for index, item := range tokens {
		if item.Token == token.String() {
			return index, item
		}
	}
	return -1, ERC20TokenWallet{}
}

type walletLogService struct{}

func newWalletLogService() *walletLogService {
	return &walletLogService{}
}

// insertNewERC20WalletLog Insert new log of ERC20 changes
func (receiver *walletLogService) insertNewERC20WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet Wallet,
) (ERC20WalletLog, error) {
	erc20WalletLog := parseCommandToERC20WalletLog(command, currentWallet)
	return erc20LogDAO.insertERC20WalletLog(db, erc20WalletLog)
}

// updateERC20WalletLog Change the status of ERC20Log in batches
func (receiver *walletLogService) updateERC20WalletLog(
	db *gorm.DB, log ERC20WalletLog, status WalletLogStatus, newWallet Wallet,
) (ERC20WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return erc20LogDAO.updateERC20WalletLogStatus(db, log)
}

// insertNewERC1155WalletLog Insert an ERC1155 asset change log
func (receiver *walletLogService) insertNewERC1155WalletLog(
	db *gorm.DB, command WalletCommand, currentWallet Wallet,
) (ERC1155WalletLog, error) {
	erc1155WalletData := parseCommandToERC1155WalletLog(command, currentWallet)
	return erc1155LogDAO.insertERC1155WalletLog(db, erc1155WalletData)
}

// updateERC1155WalletLog Change the state of the ERC1155 log
func (receiver *walletLogService) updateERC1155WalletLog(
	db *gorm.DB, log ERC1155WalletLog, status WalletLogStatus, newWallet Wallet,
) (ERC1155WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return erc1155LogDAO.updateERC1155WalletLogStatus(db, log)
}

func buildInitializedCommandFromAccount(accountId uint64) WalletCommand {
	var erc20Commands []ERC20Command
	for _, item := range supportedERC20Tokens {
		erc20Commands = append(erc20Commands, ERC20Command{
			Token:   ERC20TokenEnum(item.Index),
			Value:   0,
			Decimal: item.Decimal,
		})
	}

	return WalletCommand{
		AccountId:      accountId,
		AssetType:      Other,
		ERC20Commands:  erc20Commands,
		ERC1155Command: ERC1155Command{},
		BusinessModule: "Initialization",
		ActionType:     Initialize,
		FeeCommands:    []ERC20Command{},
	}
}

package walleter

import (
	"errors"
	"os"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// WalletCommand is an object of the operation for withdrawals or depositions.
// Anytime when users deposit, withdraw or spend an asset, we should
// construct a WalletCommand object to operate the wallet in database.
type WalletCommand struct {
	// User account id. unique
	AccountId uint64

	// 0: ERC20 token, 1: erc1155 token. 2. other type.
	AssetType AssetType

	// Action of this command. initialize, withdraw, deposit, etc.
	ActionType WalletActionType

	// ERC20 command, if we want to operate ERC20 asset, this should not be nil. otherwise this must be nil.
	ERC20Commands []ERC20Command

	// ERC20 command, if we want to operate ERC1155 asset, this should not be nil. otherwise this must be nil.
	ERC1155Command ERC1155Command

	// Fee charging command, if len(FeeCommands) > 0, assets should be deducted from user's account.
	FeeCommands []ERC20Command

	// A name about which part sent this command
	BusinessModule string
}

type ERC20Command struct {
	Token   ERC20TokenEnum
	Value   float64
	Decimal uint64
}

type ERC1155Command struct {
	Ids    []uint64
	Values []uint64
}

// Walleter the library entry object.
type Walleter struct {
	db *gorm.DB
}

var feeChargerAccountId uint64

func New(db *gorm.DB, chargerAccountId uint64) *Walleter {
	migration(db)
	feeChargerAccountId = chargerAccountId
	walleter := Walleter{db: db}
	_, err := walleter.setFeeChargerAccount()
	if err != nil {
		panic("initialize fee charger account failed")
	}
	return &walleter
}

func (s *Walleter) HandleWalletCommand(db *gorm.DB, command WalletCommand) (Wallet, error) {
	switch command.ActionType {
	case Initialize:
		return initWallet(db, command)
	default:
		return updateWallet(db, command)
	}
}

func (s *Walleter) GetWalletByAccountId(accountId uint64) (Wallet, error) {
	return walletDAO.getWallet(s.db, accountId)
}

// initialize fee charger account in database.
func (s *Walleter) setFeeChargerAccount() (Wallet, error) {
	if feeChargerAccountId == 0 {
		panic("Please assign official fee charge account.")
	}

	wallet, err := s.GetWalletByAccountId(feeChargerAccountId)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		command := buildInitializedCommandFromAccount(feeChargerAccountId)
		return s.HandleWalletCommand(s.db, command)
	}
	return wallet, nil
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

		// 3. generate a new check sign
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
	}
	return Wallet{}, errors.New("not support current asset type")
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

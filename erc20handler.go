package walleter

import (
	"gorm.io/gorm"
)

// This function doesn't contain Transaction, so if we should wrap this function in a Transaction out of this function.
func handleERC20Command(db *gorm.DB, command WalletCommand) (Wallet, error) {
	logService := newWalletLogService()
	validator := newWalletValidator()

	userWallet, err := walletDAO.getWallet(db, command.AccountId)
	if err != nil {
		return Wallet{}, err
	}

	// 1. Verify that the user's current wallet status is normal
	result, err := validator.validateWallet(userWallet)
	if err != nil || !result {
		return Wallet{}, err
	}

	// 2.Insert a log message
	erc20Log, err := logService.insertNewERC20WalletLog(db, command, userWallet)
	if err != nil {
		return Wallet{}, err
	}

	// 3. Whether to charge a fee
	for _, fee := range command.FeeCommands {
		if fee.Value <= 0 {
			continue
		}
		userWallet, err = newFeeChargerService().chargeFee(db, fee, userWallet)
		if err != nil {
			logService.updateERC20WalletLog(db, erc20Log, Failed, userWallet)
			return Wallet{}, err
		}
	}

	// 4. Make changes to user assets
	switch command.ActionType {
	case Deposit:
		for _, token := range command.ERC20Commands {
			index, userERC20TokenWallet := getUserSpecifiedERC20TokenWallet(userWallet, token.Token)
			userERC20TokenWallet.Balance += token.Value
			userERC20TokenWallet.TotalDeposit += token.Value
			userWallet.ERC20TokenData[index] = userERC20TokenWallet
			err = walletDAO.updateERC20WalletData(db, userERC20TokenWallet)
			if err != nil {
				return Wallet{}, err
			}
		}
	case Withdraw:
		for _, token := range command.ERC20Commands {
			index, userERC20TokenWallet := getUserSpecifiedERC20TokenWallet(userWallet, token.Token)
			if userERC20TokenWallet.Balance < token.Value {
				return Wallet{}, ErrNoEnoughERC20Balance
			}
			userERC20TokenWallet.Balance -= token.Value
			userERC20TokenWallet.TotalWithdraw += token.Value
			userWallet.ERC20TokenData[index] = userERC20TokenWallet
			err = walletDAO.updateERC20WalletData(db, userERC20TokenWallet)
			if err != nil {
				return Wallet{}, err
			}
		}
	case Income:
		for _, token := range command.ERC20Commands {
			index, userERC20TokenWallet := getUserSpecifiedERC20TokenWallet(userWallet, token.Token)
			userERC20TokenWallet.Balance += token.Value
			userERC20TokenWallet.TotalIncome += token.Value
			userWallet.ERC20TokenData[index] = userERC20TokenWallet
			err = walletDAO.updateERC20WalletData(db, userERC20TokenWallet)
			if err != nil {
				return Wallet{}, err
			}
		}
	case Spend, ChargeFee:
		for _, token := range command.ERC20Commands {
			userWallet, err := newFeeChargerService().chargeFee(db, token, userWallet)
			if err != nil {
				return userWallet, err
			}
		}
	default:
		return Wallet{}, ErrActionTypeNotSupport
	}

	// 6. Generate new verification information
	newCheckSign, err := validator.generateNewSignHash(userWallet)
	if err != nil {
		return Wallet{}, err
	}
	userWallet.CheckSign = newCheckSign
	err = walletDAO.updateWalletCheckSign(db, userWallet)
	if err != nil {
		return Wallet{}, err
	}

	// 8. Update log information
	_, err = newWalletLogService().updateERC20WalletLog(db, erc20Log, Done, userWallet)
	if err != nil {
		return Wallet{}, err
	}
	return walletDAO.getWallet(db, command.AccountId)
}

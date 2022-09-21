package walleter

import (
	"gorm.io/gorm"
)

func handleERC1155Command(db *gorm.DB, command WalletCommand) (Wallet, error) {
	if len(command.ERC1155Command.Values) != len(command.ERC1155Command.Ids) {
		return Wallet{}, ErrIncorrectERC1155Param
	}
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
	erc1155Log, err := logService.insertNewERC1155WalletLog(db, command, userWallet)
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
			logService.updateERC1155WalletLog(db, erc1155Log, Failed, userWallet)
			return Wallet{}, err
		}
	}

	// 5. Make changes to user assets
	ids := convertStringToUIntArray(userWallet.ERC1155TokenData.Ids)
	values := convertStringToUIntArray(userWallet.ERC1155TokenData.Values)
	switch command.ActionType {
	case Deposit, Income:
		for index, id := range command.ERC1155Command.Ids {
			value := command.ERC1155Command.Values[index]
			i := indexOfArray(ids, id)
			if i == -1 {
				ids = append(ids, id)
				values = append(values, value)
			} else {
				values[i] = values[i] + value
			}

			userWallet.ERC1155TokenData.Ids = convertArrayToString(ids, ",")
			userWallet.ERC1155TokenData.Values = convertArrayToString(values, ",")
			err = walletDAO.updateERC1155WalletData(db, userWallet.ERC1155TokenData)
			if err != nil {
				return Wallet{}, err
			}
		}
	case Withdraw, Spend:
		for index, id := range command.ERC1155Command.Ids {
			value := command.ERC1155Command.Values[index]
			i := indexOfArray(ids, id)
			if i == -1 {
				return Wallet{}, ErrNoEnoughNFT
			} else {
				if values[i] < value {
					return Wallet{}, ErrNoEnoughNFT
				}
				values[i] = values[i] - value
			}

			userWallet.ERC1155TokenData.Ids = convertArrayToString(ids, ",")
			userWallet.ERC1155TokenData.Values = convertArrayToString(values, ",")
			err = walletDAO.updateERC1155WalletData(db, userWallet.ERC1155TokenData)
			if err != nil {
				return Wallet{}, err
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
	_, err = newWalletLogService().updateERC1155WalletLog(db, erc1155Log, Done, userWallet)
	if err != nil {
		return Wallet{}, err
	}
	return walletDAO.getWallet(db, command.AccountId)
}

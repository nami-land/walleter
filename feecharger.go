package walleter

import (
	"gorm.io/gorm"
)

// fee charging service.
type feeChargerService struct{}

func newFeeChargerService() *feeChargerService {
	return &feeChargerService{}
}

func (*feeChargerService) chargeFee(db *gorm.DB, token ERC20Command, userWallet Wallet) (Wallet, error) {
	// get fee charger account.
	feeChargerWallet, err := walletDAO.getWallet(db, feeChargerAccountId)
	if err != nil {
		return userWallet, err
	}

	index, userERC20TokenWallet := getUserSpecifiedERC20TokenWallet(userWallet, token.Token)
	if index == -1 || userERC20TokenWallet.Balance < token.Value {
		return userWallet, ErrNoEnoughBalanceForFee
	}

	userERC20TokenWallet.Balance -= token.Value
	userERC20TokenWallet.TotalFee += token.Value
	userWallet.ERC20TokenData[index] = userERC20TokenWallet
	err = walletDAO.updateERC20WalletData(db, userWallet.ERC20TokenData[index])
	if err != nil {
		return userWallet, err
	}

	index, feeChargerERC20TokenWallet := getUserSpecifiedERC20TokenWallet(feeChargerWallet, token.Token)
	if index == -1 {
		return feeChargerWallet, ErrCannotFindERC20Wallet
	}
	feeChargerERC20TokenWallet.Balance += token.Value
	feeChargerERC20TokenWallet.TotalIncome += token.Value
	feeChargerWallet.ERC20TokenData[index] = feeChargerERC20TokenWallet

	// update fee charger account
	err = walletDAO.updateERC20WalletData(db, feeChargerERC20TokenWallet)
	return userWallet, err
}

// get user's specified erc20 wallet, like BUSD, NFISH wallet.
func getUserSpecifiedERC20TokenWallet(wallet Wallet, tokenType ERC20TokenEnum) (int, ERC20TokenWallet) {
	for index, item := range wallet.ERC20TokenData {
		if item.Token == tokenType.String() {
			return index, item
		}
	}
	return -1, ERC20TokenWallet{}
}

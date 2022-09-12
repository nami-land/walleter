package walleter

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
)

// walletValidator will check if the sign code of wallet is correct.
type walletValidator struct{}

func newWalletValidator() *walletValidator {
	return &walletValidator{}
}

func (receiver walletValidator) validateWallet(wallet Wallet) (bool, error) {
	checkSign := wallet.CheckSign
	md5Value, err := receiver.generateNewSignHash(wallet)
	if err != nil {
		return false, err
	}

	if md5Value != checkSign {
		return false, IncorrectCheckSignError
	}
	return true, nil
}

func (receiver walletValidator) generateNewSignHash(w Wallet) (string, error) {
	var newERC20TokenData []ERC20TokenWallet
	for _, token := range w.ERC20TokenData {
		erc20Data := ERC20TokenWallet{
			Model:         gorm.Model{},
			AccountId:     token.AccountId,
			Token:         token.Token,
			Balance:       token.Balance,
			Decimal:       token.Decimal,
			TotalIncome:   token.TotalIncome,
			TotalSpend:    token.TotalSpend,
			TotalDeposit:  token.TotalDeposit,
			TotalWithdraw: token.TotalWithdraw,
			TotalFee:      token.TotalFee,
		}
		newERC20TokenData = append(newERC20TokenData, erc20Data)
	}

	erc1155Data := ERC1155TokenWallet{
		Model:     gorm.Model{},
		AccountId: w.ERC1155TokenData.AccountId,
		Ids:       w.ERC1155TokenData.Ids,
		Values:    w.ERC1155TokenData.Values,
	}

	tempWallet := Wallet{
		Model:            gorm.Model{},
		AccountId:        w.AccountId,
		ERC20TokenData:   newERC20TokenData,
		ERC1155TokenData: erc1155Data,
		CheckSign:        "",
	}

	b, err := json.Marshal(tempWallet)
	if err != nil {
		return "", err
	}

	md5Value := md5Value(string(b))
	return md5Value, nil
}

func md5Value(str string) string {
	data := []byte(str)
	has := md5.Sum(data)
	md5str := fmt.Sprintf("%x", has)
	return md5str
}

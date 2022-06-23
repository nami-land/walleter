package pkg

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"neco-wallet-center/internal/model"
)

type walletValidator struct{}

func NewWalletValidator() *walletValidator {
	return &walletValidator{}
}

func (receiver walletValidator) ValidateWallet(wallet model.Wallet) (bool, error) {
	checkSign := wallet.CheckSign
	md5Value, err := receiver.GenerateNewSignHash(wallet)
	if err != nil {
		return false, err
	}
	if md5Value != checkSign {
		return false, errors.New("check sign is invalid")
	}
	return true, nil
}

func (receiver walletValidator) GenerateNewSignHash(wallet model.Wallet) (string, error) {
	var newERC20TokenData []model.ERC20TokenWallet
	for _, token := range wallet.ERC20TokenData {
		erc20Data := model.ERC20TokenWallet{
			Model:         gorm.Model{},
			GameClient:    token.GameClient,
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

	erc1155Data := model.ERC1155TokenWallet{
		Model:      gorm.Model{},
		GameClient: wallet.ERC1155TokenData.GameClient,
		AccountId:  wallet.ERC1155TokenData.AccountId,
		Ids:        wallet.ERC1155TokenData.Ids,
		Values:     wallet.ERC1155TokenData.Values,
	}

	tempWallet := model.Wallet{
		Model:            gorm.Model{},
		GameClient:       wallet.GameClient,
		AccountId:        wallet.AccountId,
		PublicAddress:    wallet.PublicAddress,
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
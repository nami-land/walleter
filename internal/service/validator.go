package service

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"neco-wallet-center/internal/model"
)

type walletValidator struct{}

func NewWalletValidator() *walletValidator {
	return &walletValidator{}
}

func (receiver walletValidator) ValidateWallet(wallet model.Wallet) (bool, error) {
	checkSign := wallet.CheckSign
	wallet.CheckSign = ""
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
	b, err := json.Marshal(wallet)
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

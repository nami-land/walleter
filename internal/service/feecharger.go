package service

import (
	"errors"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/model/initial"
)

type feeChargerService struct{}

func NewFeeChargerService() *feeChargerService {
	return &feeChargerService{}
}

func (receiver *feeChargerService) ChargeFee(
	db *gorm.DB, command model.WalletCommand, userWallet model.Wallet,
) (model.Wallet, error) {
	feeChargerWallet, err := model.WalletDAO.GetWallet(
		db, command.GameClient, initial.GetFeeChargerAccount(command.GameClient).AccountId,
	)
	if err != nil {
		return userWallet, err
	}

	if len(command.FeeCommands) > 0 {
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
		}
	}

	err = model.WalletDAO.UpdateWallet(db, feeChargerWallet)
	if err != nil {
		return userWallet, err
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

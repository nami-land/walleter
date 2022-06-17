package service

import (
	"context"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
)

type erc20TokenWalletService struct{}

func NewERC20TokenWalletService() *erc20TokenWalletService {
	return &erc20TokenWalletService{}
}

func (receiver erc20TokenWalletService) handleERC20Deposit(
	ctx context.Context, command model.WalletCommand,
) error {
	err := model.GetDb(ctx).Transaction(func(tx *gorm.DB) error {
		userWallet, err := model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		log, err := NewWalletLogService().InsertNewERC20WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		userWallet, err = NewFeeChargerService().ChargeFee(tx, command, userWallet)
		if err != nil {
			return err
		}

		for _, token := range command.ERC20Commands {
			index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
			userERC20TokenWallet.Balance += token.Value
			userERC20TokenWallet.TotalDeposit += token.Value
			userWallet.ERC20TokenData[index] = userERC20TokenWallet
		}

		err = model.WalletDAO.UpdateWallet(tx, userWallet)
		if err != nil {
			return err
		}

		_, err = NewWalletLogService().UpdateERC20WalletLog(tx, log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserERC20TokenWallet(tokens []model.ERC20TokenWallet, token comm.ERC20Token) (int, model.ERC20TokenWallet) {
	for index, item := range tokens {
		if item.Token == token.String() {
			return index, item
		}
	}
	return -1, model.ERC20TokenWallet{}
}

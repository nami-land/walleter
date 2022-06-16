package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
)

type walletCenterService struct{}

func NewWalletCenterService() *walletCenterService {
	return &walletCenterService{}
}

func (s walletCenterService) HandleWalletCommand(ctx context.Context, command model.WalletCommand) error {
	if !command.GameClient.IsSupport() {
		return errors.New("game client is invalid")
	}

	switch command.ActionType {
	case comm.Initialize:
		return initWallet(ctx, command)
	default:
		return updateWallet(ctx, command)
	}
}

func initWallet(ctx context.Context, command model.WalletCommand) error {
	err := model.GetDb(ctx).Transaction(func(tx1 *gorm.DB) error {
		// 1. Insert change logs, including ERC20 logs and ERC1155 Log.
		erc20WalletLogDataArray := generateERC20WalletLogArray(command, model.Wallet{})
		for index, logData := range erc20WalletLogDataArray {
			newLog, err := model.ERC20WalletLogDAO.InsertERC20WalletLog(tx1, &logData)
			if err != nil {
				return err
			}
			erc20WalletLogDataArray[index] = *newLog
		}

		erc115WalletLog := generateERC1155WalletLog(command, model.Wallet{})
		nftWalletLog, err := model.ERC1155WalletLogDAO.InsertERC1155WalletLog(tx1, &erc115WalletLog)
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := generateERC20DataArray(command)
		erc1155Data := model.ERC1155TokenData{
			Model:      gorm.Model{},
			GameClient: int(command.GameClient),
			AccountId:  command.AccountId,
		}
		wallet := model.Wallet{
			Model:            gorm.Model{},
			GameClient:       int(command.GameClient),
			AccountId:        command.AccountId,
			PublicAddress:    command.PublicAddress,
			ERC20TokenData:   erc20DataArray,
			ERC1155TokenData: erc1155Data,
			CheckSign:        "",
		}
		err = model.WalletDAO.CreateWallet(ctx, wallet)
		if err != nil {
			return err
		}

		// 3. change log statuses
		for _, logData := range erc20WalletLogDataArray {
			logData.SettledWallet = wallet
			logData.Status = comm.Done.String()
			_, err = model.ERC20WalletLogDAO.UpdateERC20WalletLogStatus(tx1, &logData)
			if err != nil {
				return err
			}
		}

		nftWalletLog.SettledWallet = wallet
		nftWalletLog.Status = comm.Done.String()
		_, err = model.ERC1155WalletLogDAO.UpdateERC1155WalletLogStatus(tx1, nftWalletLog)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func updateWallet(ctx context.Context, command model.WalletCommand) error {
	return nil
}

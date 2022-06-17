package service

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/pkg"
)

type walletCenterService struct{}

func NewWalletCenterService() *walletCenterService {
	return &walletCenterService{}
}

func (s *walletCenterService) HandleWalletCommand(ctx context.Context, command model.WalletCommand) error {
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
		walletLogService := NewWalletLogService()
		erc20WalletLog, err := walletLogService.InsertNewERC20WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		erc115WalletLog, err := walletLogService.InsertNewERC1155WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := pkg.ParseCommandToERC20WalletArray(command)
		erc1155Data := pkg.ParseCommandToERC1155Wallet(command)
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
		_, err = walletLogService.UpdateERC20WalletLog(tx1, erc20WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		_, err = walletLogService.UpdateERC1155WalletLog(tx1, erc115WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		return nil
	})
	return err
}

func updateWallet(ctx context.Context, command model.WalletCommand) error {
	userWallet, err := model.WalletDAO.GetWallet(ctx, command.GameClient, command.AccountId)
	if err != nil {
		return err
	}

	// 1. charge fee
	if len(command.FeeCommands) > 0 {
		for _, fee := range command.FeeCommands {
			if fee.Value <= 0 {
				continue
			}
			userTokenData := getUserERC20TokenData(userWallet.ERC20TokenData, fee.Token)
			if userTokenData == nil {
				return errors.New("insufficient balance for fee")
			}

			userTokenData.Balance -= fee.Value
			userTokenData.TotalFee += fee.Value
		}
	}

	switch command.AssetType {
	case comm.ERC20AssetType:
		return handleERC20Command(ctx, command)
	case comm.ERC1155AssetType:
		return handleERC1155Command(ctx, command)
	default:
		return errors.New("not support current asset type")
	}
}

func handleERC20Command(ctx context.Context, command model.WalletCommand) error {
	switch command.ActionType {
	case comm.Deposit:
		fmt.Println("erc20 deposit")
	case comm.Withdraw:
		fmt.Println("erc20 withdraw")
	case comm.Income:
		fmt.Println("erc20 income")
	case comm.Spend:
		fmt.Println("erc20 spend")
	case comm.ChargeFee:
		fmt.Println("erc20 charge fee")
	default:
		return errors.New("not support action type")
	}
	return nil
}

func handleERC1155Command(ctx context.Context, command model.WalletCommand) error {
	switch command.ActionType {
	case comm.Deposit:
		fmt.Println("erc1155 deposit")
	case comm.Withdraw:
		fmt.Println("erc1155 withdraw")
	case comm.Income:
		fmt.Println("erc1155 income")
	case comm.Spend:
		fmt.Println("erc1155 spend")
	case comm.ChargeFee:
		fmt.Println("erc1155 charge fee")
	default:
		return errors.New("not support action type")
	}
	return nil
}

func getUserERC20TokenData(tokens []model.ERC20TokenWallet, token comm.ERC20Token) *model.ERC20TokenWallet {
	for _, item := range tokens {
		if item.Token == token.String() {
			return &item
		}
	}
	return nil
}

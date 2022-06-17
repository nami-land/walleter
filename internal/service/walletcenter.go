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
		err = model.WalletDAO.CreateWallet(tx1, wallet)
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
	err := model.GetDb(ctx).Transaction(func(tx *gorm.DB) error {
		logService := NewWalletLogService()
		userWallet, err := model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		// 1.插入一条log信息
		log, err := logService.InsertNewERC20WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 2. 收取手续费
		userWallet, err = NewFeeChargerService().ChargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.UpdateERC20WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 3. 对用户资产进行变更
		switch command.ActionType {
		case comm.Deposit:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalDeposit += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Withdraw:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalWithdraw += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Income:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalIncome += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Spend:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalSpend += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.ChargeFee:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalFee += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		default:
			return errors.New("not support action type")
		}

		// 4. 更新用户资产
		err = model.WalletDAO.UpdateWallet(tx, userWallet)
		if err != nil {
			_, err = logService.UpdateERC20WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 5. 更新log信息
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

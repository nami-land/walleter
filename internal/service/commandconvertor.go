package service

import (
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/utils"
)

func generateERC20DataArray(command model.WalletCommand) []model.ERC20TokenData {
	var result []model.ERC20TokenData
	for _, item := range command.ERC20Commands {
		data := model.ERC20TokenData{
			Model:      gorm.Model{},
			GameClient: int(command.GameClient),
			AccountId:  command.AccountId,
			TokenType:  item.Token.String(),
		}
		result = append(result, data)
	}
	return result
}

func generateERC20WalletLogArray(command model.WalletCommand, wallet model.Wallet) []model.ERC20WalletLog {
	var result []model.ERC20WalletLog
	var fees []model.WalletFeeItem
	for _, fee := range command.FeeCommands {
		feeData := model.WalletFeeItem{
			TokenType: fee.Token.String(),
			Amount:    fee.Value,
			Decimal:   fee.Decimal,
		}
		fees = append(fees, feeData)
	}

	for _, item := range command.ERC20Commands {
		walletLog := model.ERC20WalletLog{
			Model:          gorm.Model{},
			GameClient:     int(command.GameClient),
			AccountId:      command.AccountId,
			BusinessModule: command.BusinessModule,
			ActionType:     command.ActionType,
			TokenType:      item.Token.String(),
			Value:          item.Value,
			Fee:            fees,
			Status:         comm.Pending.String(),
			OriginalWallet: wallet,
			SettledWallet:  model.Wallet{},
		}
		result = append(result, walletLog)
	}
	return result
}

func generateERC1155WalletLog(command model.WalletCommand, wallet model.Wallet) model.ERC1155WalletLog {
	var fees []model.WalletFeeItem
	for _, fee := range command.FeeCommands {
		feeData := model.WalletFeeItem{
			TokenType: fee.Token.String(),
			Amount:    fee.Value,
			Decimal:   fee.Decimal,
		}
		fees = append(fees, feeData)
	}

	return model.ERC1155WalletLog{
		Model:          gorm.Model{},
		GameClient:     int(command.GameClient),
		AccountId:      command.AccountId,
		BusinessModule: command.BusinessModule,
		ActionType:     command.ActionType,
		Ids:            utils.ConvertUintArrayToString(command.ERC1155Command.Ids),
		Values:         utils.ConvertUintArrayToString(command.ERC1155Command.Values),
		Fee:            fees,
		Status:         comm.Pending.String(),
		OriginalWallet: wallet,
		SettledWallet:  model.Wallet{},
	}
}

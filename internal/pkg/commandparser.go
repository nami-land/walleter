package pkg

import (
	"github.com/neco-fun/wallet-center/core"
	"github.com/neco-fun/wallet-center/internal/comm"
	"github.com/neco-fun/wallet-center/internal/model"
	"github.com/neco-fun/wallet-center/internal/utils"
	"gorm.io/gorm"
)

func ParseCommandToERC20WalletArray(command core.WalletCommand) []model.ERC20TokenWallet {
	var result []model.ERC20TokenWallet
	for _, item := range command.ERC20Commands {
		data := model.ERC20TokenWallet{
			Model:     gorm.Model{},
			AccountId: command.AccountId,
			Token:     item.Token.String(),
			Decimal:   item.Decimal,
		}
		result = append(result, data)
	}
	return result
}

func ParseCommandToERC1155Wallet(command core.WalletCommand) model.ERC1155TokenWallet {
	return model.ERC1155TokenWallet{
		Model:     gorm.Model{},
		AccountId: command.AccountId,
		Ids:       utils.ConvertUintArrayToString(command.ERC1155Command.Ids, ","),
		Values:    utils.ConvertUintArrayToString(command.ERC1155Command.Values, ","),
	}
}

func ParseCommandToERC20WalletLog(command core.WalletCommand, wallet model.Wallet) model.ERC20WalletLog {
	fees := ParseERC20Commands(command.FeeCommands)
	gonnaChangedTokens := ParseERC20Commands(command.ERC20Commands)

	return model.ERC20WalletLog{
		Model:          gorm.Model{},
		AccountId:      command.AccountId,
		BusinessModule: command.BusinessModule,
		ActionType:     command.ActionType.String(),
		Tokens:         model.ERC20TokenCollection{Items: gonnaChangedTokens},
		Fees:           model.ERC20TokenCollection{Items: fees},
		Status:         comm.Pending.String(),
		OriginalWallet: wallet,
		SettledWallet:  model.Wallet{},
	}
}

func ParseCommandToERC1155WalletLog(command core.WalletCommand, wallet model.Wallet) model.ERC1155WalletLog {
	fees := ParseERC20Commands(command.FeeCommands)

	return model.ERC1155WalletLog{
		Model:          gorm.Model{},
		AccountId:      command.AccountId,
		BusinessModule: command.BusinessModule,
		ActionType:     command.ActionType.String(),
		Ids:            utils.ConvertUintArrayToString(command.ERC1155Command.Ids, ","),
		Values:         utils.ConvertUintArrayToString(command.ERC1155Command.Values, ","),
		Fees:           model.ERC20TokenCollection{Items: fees},
		Status:         comm.Pending.String(),
		OriginalWallet: wallet,
		SettledWallet:  model.Wallet{},
	}
}

func ParseERC20Commands(commands []core.ERC20Command) []model.ERC20TokenData {
	var result []model.ERC20TokenData
	for _, data := range commands {
		tokenData := model.ERC20TokenData{
			TokenType: data.Token.String(),
			Amount:    data.Value,
			Decimal:   data.Decimal,
		}
		result = append(result, tokenData)
	}
	return result
}

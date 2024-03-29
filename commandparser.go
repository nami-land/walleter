package walleter

import (
	"gorm.io/gorm"
)

// parseCommandToERC20WalletArray convert erc20command array to erc20TokenWallet Array.
func parseCommandToERC20WalletArray(command WalletCommand) []ERC20TokenWallet {
	var result []ERC20TokenWallet
	for _, item := range command.ERC20Commands {
		data := ERC20TokenWallet{
			Model:     gorm.Model{},
			AccountId: command.AccountId,
			Token:     item.Token.String(),
			Decimal:   item.Decimal,
		}
		result = append(result, data)
	}
	return result
}

func parseCommandToERC1155Wallet(command WalletCommand) ERC1155TokenWallet {
	return ERC1155TokenWallet{
		Model:     gorm.Model{},
		AccountId: command.AccountId,
		Ids:       convertArrayToString(command.ERC1155Command.Ids, ","),
		Values:    convertArrayToString(command.ERC1155Command.Values, ","),
	}
}

func parseCommandToERC20WalletLog(command WalletCommand, w Wallet) ERC20WalletLog {
	fees := parseERC20Commands(command.FeeCommands)
	gonnaChangedTokens := parseERC20Commands(command.ERC20Commands)

	return ERC20WalletLog{
		Model:          gorm.Model{},
		AccountId:      command.AccountId,
		BusinessModule: command.BusinessModule,
		ActionType:     command.ActionType.String(),
		Tokens:         erc20TokenCollection{Items: gonnaChangedTokens},
		Fees:           erc20TokenCollection{Items: fees},
		Status:         Pending.String(),
		OriginalWallet: w,
		Source:         command.CommandSource.String(),
		SettledWallet:  Wallet{},
	}
}

func parseCommandToERC1155WalletLog(command WalletCommand, w Wallet) ERC1155WalletLog {
	fees := parseERC20Commands(command.FeeCommands)

	return ERC1155WalletLog{
		Model:          gorm.Model{},
		AccountId:      command.AccountId,
		BusinessModule: command.BusinessModule,
		ActionType:     command.ActionType.String(),
		Ids:            convertArrayToString(command.ERC1155Command.Ids, ","),
		Values:         convertArrayToString(command.ERC1155Command.Values, ","),
		Fees:           erc20TokenCollection{Items: fees},
		Status:         Pending.String(),
		Source:         command.CommandSource.String(),
		OriginalWallet: w,
		SettledWallet:  Wallet{},
	}
}

func parseERC20Commands(commands []ERC20Command) []erc20TokenData {
	var result []erc20TokenData
	for _, data := range commands {
		tokenData := erc20TokenData{
			TokenType: data.Token.String(),
			Amount:    data.Value,
			Decimal:   data.Decimal,
		}
		result = append(result, tokenData)
	}
	return result
}

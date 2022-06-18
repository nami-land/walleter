package pkg

import (
	"neco-wallet-center/api/pb"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/utils"
)

type commandBuilder struct{}

func NewCommandBuilder() *commandBuilder {
	return &commandBuilder{}
}

func (receiver *commandBuilder) BuildCommandFromRequest(request *pb.UpdateUserWalletRequest) model.WalletCommand {
	var erc20Commands []model.ERC20Command
	for _, token := range request.ERC20TokenData {
		command := model.ERC20Command{
			Token:   comm.ERC20Token(token.Token),
			Value:   float64(token.Balance),
			Decimal: uint(token.Decimal),
		}
		erc20Commands = append(erc20Commands, command)
	}

	var feeCommands []model.ERC20Command
	for _, token := range request.FeeData {
		command := model.ERC20Command{
			Token:   comm.ERC20Token(token.Token),
			Value:   float64(token.Balance),
			Decimal: uint(token.Decimal),
		}
		feeCommands = append(feeCommands, command)
	}

	erc1155Command := model.ERC1155Command{
		Ids:    utils.ConvertUInt32ArrayToUIntArray(request.ERC1155TokenData.Ids),
		Values: utils.ConvertUInt32ArrayToUIntArray(request.ERC1155TokenData.Values),
	}

	return model.WalletCommand{
		GameClient:     comm.GameClient(request.GameClient),
		AccountId:      uint(request.AccountId),
		PublicAddress:  "",
		AssetType:      comm.AssetType(request.AssetType),
		ERC20Commands:  erc20Commands,
		ERC1155Command: erc1155Command,
		BusinessModule: request.BusinessModule,
		ActionType:     comm.WalletActionType(pb.WalletActionType_Initialize),
		FeeCommands:    feeCommands,
	}
}

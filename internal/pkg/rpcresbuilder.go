package pkg

import (
	"neco-wallet-center/api/pb"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/utils"
)

type rpcResponseBuilder struct{}

func NewRPCResponseBuilder() *rpcResponseBuilder {
	return &rpcResponseBuilder{}
}

func (receiver rpcResponseBuilder) BuilderRPCResponseWallet(wallet model.Wallet) pb.UserWallet {
	return pb.UserWallet{
		Id:            uint64(wallet.ID),
		GameClient:    pb.GameClient(wallet.GameClient),
		AccountId:     uint64(wallet.AccountId),
		PublicAddress: wallet.PublicAddress,
		ERC20Tokens:   convertERC20Data(wallet.ERC20TokenData),
		ERC1155Token:  convertERC1155Data(wallet.ERC1155TokenData),
	}
}

func convertERC20Data(array []model.ERC20TokenWallet) []*pb.ERC20TokenWallet {
	var result []*pb.ERC20TokenWallet
	for _, item := range array {
		tokenData := pb.ERC20TokenWallet{
			Token:   pb.ERC20Token(comm.GetERC20TokenType(item.Token)),
			Balance: float32(item.Balance),
			Decimal: uint64(item.Decimal),
		}
		result = append(result, &tokenData)
	}
	return result
}

func convertERC1155Data(data model.ERC1155TokenWallet) *pb.ERC1155TokenWallet {
	return &pb.ERC1155TokenWallet{
		Ids:    utils.CovertUIntArrayToUInt64Array(utils.ConvertStringToUIntArray(data.Ids)),
		Values: utils.CovertUIntArrayToUInt64Array(utils.ConvertStringToUIntArray(data.Values)),
	}
}

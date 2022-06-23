package initial

import (
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
)

type OfficialAccount struct {
	AccountId     uint
	PublicAddress string
}

var necoFishingFeeChargerAccount = OfficialAccount{
	AccountId:     1,
	PublicAddress: "0xa98Ff091a5F6975162AEa4E3862165bCf81aB4Ad",
}

func GetFeeChargerAccount(gameClient comm.GameClient) OfficialAccount {
	if gameClient == comm.NecoFishing {
		return necoFishingFeeChargerAccount
	}
	return necoFishingFeeChargerAccount
}

var initNecoFishingFeeChargerAccountCommand = model.WalletCommand{
	GameClient:    comm.NecoFishing,
	AccountId:     GetFeeChargerAccount(comm.NecoFishing).AccountId,
	PublicAddress: GetFeeChargerAccount(comm.NecoFishing).PublicAddress,
	AssetType:     comm.Other,
	ERC20Commands: []model.ERC20Command{
		{
			Token:   comm.NFISH,
			Value:   0,
			Decimal: 18,
		}, {
			Token:   comm.BUSD,
			Value:   0,
			Decimal: 18,
		},
	},
	ERC1155Command: model.ERC1155Command{},
	BusinessModule: "Initialization",
	ActionType:     comm.Initialize,
	FeeCommands:    []model.ERC20Command{},
}

var InitializedCommands = []model.WalletCommand{
	initNecoFishingFeeChargerAccountCommand,
}

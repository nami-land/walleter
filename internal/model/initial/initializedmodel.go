package initial

import (
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
)

var InitNecoFishingFeeChargerAccountCommand = model.WalletCommand{
	GameClient:    comm.NecoFishing,
	AccountId:     comm.GetFeeChargerAccount(comm.NecoFishing).AccountId,
	PublicAddress: comm.GetFeeChargerAccount(comm.NecoFishing).PublicAddress,
	AssetType:     comm.Other,
	ERC20Commands: []model.ERC20Command{model.ERC20Command{
		Token:   comm.NFISH,
		Value:   0,
		Decimal: 18,
	}, model.ERC20Command{
		Token:   comm.BUSD,
		Value:   0,
		Decimal: 18,
	}},
	ERC1155Command: model.ERC1155Command{},
	BusinessModule: "Initialization",
	ActionType:     comm.Initialize,
	FeeCommands:    []model.ERC20Command{},
}

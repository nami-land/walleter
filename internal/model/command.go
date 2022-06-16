package model

import "neco-wallet-center/internal/comm"

type WalletCommand struct {
	GameClient     comm.GameClient
	AccountId      uint
	PublicAddress  string
	AssetType      comm.AssetType
	ERC20Commands  []ERC20Command
	ERC1155Command ERC1155Command
	BusinessModule string
	ActionType     string
	FeeCommands    []ERC20Command
}

type ERC20Command struct {
	Token   comm.ERC20Token
	Value   float64
	Decimal uint
}

type ERC1155Command struct {
	Ids    []uint
	Values []uint
}

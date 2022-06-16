package model

import "neco-wallet-center/internal/comm"

type WalletCommand struct {
	GameClient     comm.GameClient //game client, like Neco Fishing, Neco Land..etc.
	AccountId      uint            // User account id. unique
	PublicAddress  string          // public address, allowed to be null
	AssetType      comm.AssetType  // 0: ERC20 token, 1: erc1155 token.
	ERC20Commands  []ERC20Command
	ERC1155Command ERC1155Command
	BusinessModule string
	ActionType     comm.WalletActionType
	FeeCommands    []ERC20Command // charge fee, if len(FeeCommands) > 0, should be deducted from user's account.
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

package walleter

import "errors"

var (
	ErrIncorrectAssetType    = errors.New("incorrect asset type in command")
	ErrIncorrectERC1155Param = errors.New("incorrect erc1155 parameters")
	ErrIncorrectCheckSign    = errors.New("check sign is invalid")
	ErrNoEnoughNFT           = errors.New("insufficient nft balance")
	ErrNoEnoughERC20Balance  = errors.New("insufficient balance")
	ErrNoEnoughBalanceForFee = errors.New("insufficient balance for fee")
	ErrAssetTypeNotSupport   = errors.New("not support current asset type")
	ErrActionTypeNotSupport  = errors.New("not support action type")
	ErrCannotFindERC20Wallet = errors.New("cannot find erc20 wallet")
)

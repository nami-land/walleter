package walleter

import "errors"

var (
	IncorrectAssetTypeError    = errors.New("incorrect asset type in command")
	IncorrectERC1155ParamError = errors.New("incorrect erc1155 parameters")
	IncorrectCheckSignError    = errors.New("check sign is invalid")
	NoEnoughNFTError           = errors.New("insufficient nft balance")
	NoEnoughERC20BalanceError  = errors.New("insufficient balance")
	NoEnoughBalanceForFeeError = errors.New("insufficient balance for fee")
	AssetTypeNotSupportError   = errors.New("not support current asset type")
	ActionTypeNotSupportError  = errors.New("not support action type")
	CannotFindERC20WalletError = errors.New("cannot find erc20 wallet")
)

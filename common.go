package wallet_center

type AssetType int

const (
	ERC20AssetType   = 0
	ERC1155AssetType = 1
	Other            = 2
)

// ERC20TokenEnum ERC20 token name currently supported.
type ERC20TokenEnum int

const (
	NFISH ERC20TokenEnum = 0
	BUSD  ERC20TokenEnum = 1
)

func (t ERC20TokenEnum) String() string {
	switch t {
	case NFISH:
		return "NFISH"
	case BUSD:
		return "BUSD"
	}
	return "unknown"
}

type ERC20Token struct {
	Index   uint64
	Symbol  string
	Decimal uint64
}

var supportedERC20Tokens = []ERC20Token{
	{
		Index:   0,
		Symbol:  NFISH.String(),
		Decimal: 18,
	},
	{
		Index:   1,
		Symbol:  BUSD.String(),
		Decimal: 18,
	},
}

//func GetERC20TokenEnum(tokenName string) (ERC20TokenEnum, error) {
//	switch tokenName {
//	case "NFISH":
//		return NFISH, nil
//	case "BUSD":
//		return BUSD, nil
//	}
//	return NFISH, errors.New("incorrect token tokenName")
//}

// WalletActionType Actions for wallet command, we will change user's assets in wallet according to wallet action type.
type WalletActionType int

const (
	// Initialize initializing a wallet for a user when there is no wallet account of this user.
	Initialize WalletActionType = 0

	// Income will perform the addition operation to user's wallet in game database.
	Income WalletActionType = 1

	// Spend will perform subtraction to user's wallet in game database.
	Spend WalletActionType = 2

	// Deposit when uses deposit assets from Blockchain, will perform addition operation in game database.
	Deposit WalletActionType = 3

	// Withdraw when users withdraw assets from game, will perform subtraction operation in game database.
	Withdraw WalletActionType = 4

	// ChargeFee will perform subtraction operation in game database.
	ChargeFee WalletActionType = 5
)

func (t WalletActionType) String() string {
	switch t {
	case Initialize:
		return "Initialize"
	case Income:
		return "TotalIncome"
	case Spend:
		return "Spend"
	case Deposit:
		return "Deposit"
	case Withdraw:
		return "Withdraw"
	case ChargeFee:
		return "ChargeFee"
	}
	return "unknown"
}

// WalletLogStatus the results of performing wallet commands.
type WalletLogStatus int

const (
	Pending WalletLogStatus = 0
	Done    WalletLogStatus = 1
	Failed  WalletLogStatus = 2
)

func (s WalletLogStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Done:
		return "Done"
	case Failed:
		return "Failed"
	}
	return "unknown"
}

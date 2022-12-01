package walleter

type AssetType int

const (
	ERC20AssetType   = 0
	ERC1155AssetType = 1
	Other            = 2
)

// ERC20TokenEnum ERC20 token name currently supported.
type ERC20TokenEnum int

const (
	ETH ERC20TokenEnum = iota + 1
	BNB
	USDT
	USDC
	BUSD
	NAMIX
	FISHX
)

var supportedERC20Tokens = []ERC20Token{
	{
		Index:   1,
		Symbol:  ETH.String(),
		Decimal: 18,
	},
	{
		Index:   2,
		Symbol:  BNB.String(),
		Decimal: 18,
	},
	{
		Index:   3,
		Symbol:  USDT.String(),
		Decimal: 6,
	},
	{
		Index:   4,
		Symbol:  USDC.String(),
		Decimal: 6,
	},
	{
		Index:   5,
		Symbol:  BUSD.String(),
		Decimal: 18,
	},
	{
		Index:   6,
		Symbol:  NAMIX.String(),
		Decimal: 18,
	},
	{
		Index:   7,
		Symbol:  FISHX.String(),
		Decimal: 18,
	},
}

func (t ERC20TokenEnum) String() string {
	switch t {
	case ETH:
		return "ETH"
	case BNB:
		return "BNB"
	case USDC:
		return "USDC"
	case USDT:
		return "USDT"
	case BUSD:
		return "BUSD"
	case NAMIX:
		return "NAMIX"
	case FISHX:
		return "FISHX"
	}
	return "unknown"
}

type ERC20Token struct {
	Index   uint64
	Symbol  string
	Decimal uint64
}

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
		return "initialize"
	case Income:
		return "income"
	case Spend:
		return "spend"
	case Deposit:
		return "deposit"
	case Withdraw:
		return "withdraw"
	case ChargeFee:
		return "fee"
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
		return "pending"
	case Done:
		return "done"
	case Failed:
		return "failed"
	}
	return "unknown"
}

type CommandSourceType int

const (
	InGame        CommandSourceType = 0
	Ethereum      CommandSourceType = 1
	GoerliTestnet CommandSourceType = 2
	BSC           CommandSourceType = 3
	BSCTestnet    CommandSourceType = 4
)

func (s CommandSourceType) String() string {
	switch s {
	case InGame:
		return "game"
	case Ethereum:
		return "ethereum"
	case GoerliTestnet:
		return "goerli_testnet"
	case BSC:
		return "bsc"
	case BSCTestnet:
		return "bsc_testnet"
	}
	return "unknown"
}

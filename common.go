package wallet_center

type AssetType int

const (
	ERC20AssetType   = 0
	ERC1155AssetType = 1
	Other            = 2
)

type ERC20Token struct {
	Index   uint64
	Symbol  string
	Decimal uint
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
	default:
		return "NFISH"
	}
}

func GetERC20TokenType(name string) ERC20TokenEnum {
	switch name {
	case "NFISH":
		return NFISH
	case "BUSD":
		return BUSD
	default:
		return NFISH
	}
}

type WalletActionType int

const (
	Initialize WalletActionType = 0
	Income     WalletActionType = 1
	Spend      WalletActionType = 2
	Deposit    WalletActionType = 3
	Withdraw   WalletActionType = 4
	ChargeFee  WalletActionType = 5
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
	default:
		return "TotalIncome"
	}
}

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
	default:
		return "Pending"
	}
}

package comm

type AssetType int

const (
	ERC20AssetType   = 0
	ERC1155AssetType = 1
	Other            = 2
)

type ERC20Token int

const (
	NFISH ERC20Token = 0
	BUSD  ERC20Token = 1
)

func (t ERC20Token) String() string {
	switch t {
	case NFISH:
		return "NFISH"
	case BUSD:
		return "BUSD"
	default:
		return "NFISH"
	}
}

func GetERC20TokenType(name string) ERC20Token {
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

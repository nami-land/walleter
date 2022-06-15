package comm

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

type WalletActionType int

const (
	Initialize WalletActionType = 0
	Income     WalletActionType = 1
	Spend      WalletActionType = 2
	Deposit    WalletActionType = 3
	Withdraw   WalletActionType = 4
)

func (t WalletActionType) String() string {
	switch t {
	case Initialize:
		return "Initialize"
	case Income:
		return "Income"
	case Spend:
		return "Spend"
	case Deposit:
		return "Deposit"
	case Withdraw:
		return "Withdraw"
	default:
		return "Income"
	}
}

type GameClient int

const (
	NecoFishing = 0
)

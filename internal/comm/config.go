package comm

type OfficialAccount struct {
	AccountId     uint
	PublicAddress string
}

var necoFishingFeeChargerAccount = OfficialAccount{
	AccountId:     1,
	PublicAddress: "0xa98Ff091a5F6975162AEa4E3862165bCf81aB4Ad",
}

func GetFeeChargerAccount(gameClient GameClient) OfficialAccount {
	if gameClient == NecoFishing {
		return necoFishingFeeChargerAccount
	}
	return OfficialAccount{}
}

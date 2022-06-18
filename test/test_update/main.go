package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"neco-wallet-center/api/pb"
	"time"
)

func main() {
	ctx, fn := context.WithTimeout(context.Background(), 30*time.Second)
	defer fn()
	var targetUrl = "localhost:8081"
	conn, err := grpc.DialContext(ctx, targetUrl, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		log.Errorf("init NecoBlockchainAssetChain fialed: %x", err)
		panic(err)
	}
	client := pb.NewNecoWalletCenterClient(conn)
	reply, err := client.UpdateUserWallet(ctx, &pb.UpdateUserWalletRequest{
		AccountId:      11,
		GameClient:     pb.GameClient_NecoFishing,
		BusinessModule: "Update",
		AssetType:      pb.AssetType_ERC1155AssetType,
		ActionType:     pb.WalletActionType_Deposit,
		ERC20TokenData: []*pb.ERC20TokenWallet{
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_NFISH,
				Balance: 10,
				Decimal: 18,
			},
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_BUSD,
				Balance: 10,
				Decimal: 18,
			},
		},
		ERC1155TokenData: &pb.ERC1155TokenWallet{
			Ids:    []uint64{10001},
			Values: []uint64{1},
		},
		FeeData: []*pb.ERC20TokenWallet{
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_NFISH,
				Balance: 1,
				Decimal: 18,
			},
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_BUSD,
				Balance: 1,
				Decimal: 18,
			},
		},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%x", reply)
}

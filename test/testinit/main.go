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
		AccountId:      12,
		GameClient:     pb.GameClient_NecoFishing,
		BusinessModule: "Initialization",
		AssetType:      pb.AssetType_Other,
		ActionType:     pb.WalletActionType_Initialize,
		ERC20TokenData: []*pb.ERC20TokenWallet{
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_NFISH,
				Balance: 0,
				Decimal: 18,
			},
			&pb.ERC20TokenWallet{
				Token:   pb.ERC20Token_BUSD,
				Balance: 0,
				Decimal: 18,
			},
		},
		ERC1155TokenData: &pb.ERC1155TokenWallet{
			Ids:    []uint64{},
			Values: []uint64{},
		},
		FeeData: []*pb.ERC20TokenWallet{},
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%x", reply)
}

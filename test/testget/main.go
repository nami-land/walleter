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
	reply, err := client.GetUserWallet(ctx, &pb.GetUserWalletRequest{
		AccountId:  1,
		GameClient: 0,
	})
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Printf("%v", reply.ERC20Tokens[1].Token)
}

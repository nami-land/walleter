package server

import (
	"context"
	"neco-wallet-center/api/pb"
)

type walletCenterServer struct {
	pb.UnimplementedNecoWalletCenterServer
}

func (w walletCenterServer) InitUserWallet(ctx context.Context, request *pb.InitUserWalletRequest) (*pb.UserWallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w walletCenterServer) UpdateUserWallet(ctx context.Context, request *pb.UpdateUserWalletRequest) (*pb.UserWallet, error) {
	//TODO implement me
	panic("implement me")
}

func (w walletCenterServer) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.UserWallet, error) {
	//TODO implement me
	panic("implement me")
}

var _ pb.NecoWalletCenterServer = walletCenterServer{}

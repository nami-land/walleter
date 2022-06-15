package server

import (
	"context"
	"neco-wallet-center/api/pb"
)

type walletCenterServer struct {
	pb.UnimplementedNecoWalletCenterServer
}

func (w walletCenterServer) UpdateUserERC20Wallet(ctx context.Context, request *pb.UpdateUserERC20WalletRequest) (*pb.UpdateUserERC20WalletResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w walletCenterServer) UpdateUserERC1155Wallet(ctx context.Context, request *pb.UpdateUserERC1155WalletRequest) (*pb.UpdateUserERC1155WalletResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w walletCenterServer) GetUserERC20Wallet(ctx context.Context, request *pb.GetUserERC20WalletRequest) (*pb.GetUserERC20WalletResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (w walletCenterServer) GetUserERC1155Wallet(ctx context.Context, request *pb.GetUserERC1155WalletRequest) (*pb.GetUserERC1155WalletResponse, error) {
	//TODO implement me
	panic("implement me")
}

var _ pb.NecoWalletCenterServer = walletCenterServer{}

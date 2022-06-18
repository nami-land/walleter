package server

import (
	"context"
	"neco-wallet-center/api/pb"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/pkg"
	"neco-wallet-center/internal/service"
)

type walletCenterServer struct {
	pb.UnimplementedNecoWalletCenterServer
}

func (w walletCenterServer) UpdateUserWallet(ctx context.Context, request *pb.UpdateUserWalletRequest) (*pb.UserWallet, error) {
	command := pkg.NewCommandBuilder().BuildCommandFromRequest(request)
	wallet, err := service.NewWalletCenterService().HandleWalletCommand(ctx, command)
	if err != nil {
		return nil, err
	}
	userWallet := pkg.NewRPCResponseBuilder().BuilderRPCResponseWallet(wallet)
	return &userWallet, nil
}

func (w walletCenterServer) GetUserWallet(ctx context.Context, request *pb.GetUserWalletRequest) (*pb.UserWallet, error) {
	wallet, err := model.WalletDAO.GetWallet(model.GetDb(ctx), comm.GameClient(request.GameClient), uint(request.AccountId))
	if err != nil {
		return nil, err
	}

	userWallet := pkg.NewRPCResponseBuilder().BuilderRPCResponseWallet(wallet)
	return &userWallet, nil
}

var _ pb.NecoWalletCenterServer = walletCenterServer{}

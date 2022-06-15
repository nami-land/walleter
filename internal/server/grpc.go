package server

import (
	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"neco-wallet-center/api/pb"
)

func NewGrpcServer() *grpc.Server {
	customFunc := func(p interface{}) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}
	opts := []grpcRecovery.Option{
		grpcRecovery.WithRecoveryHandler(customFunc),
	}

	srv := grpc.NewServer(
		grpcmiddleware.WithUnaryServerChain(grpcRecovery.UnaryServerInterceptor(opts...)),
	)
	pb.RegisterNecoWalletCenterServer(srv, &walletCenterServer{})
	return srv
}

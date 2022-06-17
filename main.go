package main

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/utils"
	"os"
)

func main() {
	config := &utils.Config{}
	config, err := utils.GetConfig("config.dev.yaml")
	if err != nil {
		return
	}

	db, err := model.InitDB(config)
	if err != nil {
		_ = fmt.Errorf("connect error")
		return
	}
	migration(db)

	model.WalletDAO.InitWallet(context.Background(), comm.NecoFishing, 1, "0x53e310e72592591263AF6b07EaC18e6bB08eD5bb")

	//l, err := net.Listen("tcp", ":8081")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//log.Infoln("start gRPC server")
	//grpcServer := server.NewGrpcServer()
	//err = grpcServer.Serve(l)
	//if err != nil {
	//	log.Fatal("Launch gRPC server failed.")
	//}
}

func migration(db *gorm.DB) {
	_ = db.Table("t_erc20_token_data_0").AutoMigrate(model.ERC20TokenWallet{})
	_ = db.Table("t_erc1155_token_data_0").AutoMigrate(model.ERC1155TokenWallet{})
	_ = db.Table("t_wallet_0").AutoMigrate(model.Wallet{})
	_ = db.Table("t_erc20_wallet_log_0").AutoMigrate(model.ERC20WalletLog{})
	_ = db.Table("t_erc1155_wallet_log_0").AutoMigrate(model.ERC1155WalletLog{})
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

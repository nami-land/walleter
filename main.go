package main

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/utils"
	"os"
)

func main() {
	fmt.Println("hello World")

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
}

func migration(db *gorm.DB) {
	_ = db.AutoMigrate(model.ERC20TokenData{})
	_ = db.AutoMigrate(model.ERC1155TokenData{})
	_ = db.AutoMigrate(model.FishingWallet{})
	_ = db.AutoMigrate(model.ERC20WalletLog{})
	_ = db.AutoMigrate(model.ERC1155WalletLog{})
}

func init() {
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
}

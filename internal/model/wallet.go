package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
)

type Wallet struct {
	gorm.Model       `swagger-ignore:"true"`
	GameClient       int              `json:"game_client"`
	AccountId        uint             `json:"account_id" gorm:"unique; not null"` // 玩家账户ID
	PublicAddress    string           `json:"address" gorm:"unique;not null"`     // 玩家的钱包地址
	ERC20TokenData   []ERC20TokenData `json:"erc_20_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	ERC1155TokenData ERC1155TokenData `json:"erc_1155_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	CheckSign        string           `json:"check_sign" gorm:"type:varchar(128);not null;comment:'安全签名'"`
}

func (wallet Wallet) TableName() string {
	return fmt.Sprintf("t_wallet_%d", wallet.GameClient)
}

func (wallet Wallet) Value() (driver.Value, error) {
	b, err := json.Marshal(wallet)
	return string(b), err
}

func (wallet *Wallet) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), wallet)
}

type ERC20TokenData struct {
	gorm.Model    `swagger-ignore:"true"`
	GameClient    int     `json:"game_client"`
	AccountId     uint    `json:"account_id"`     //往家账户的ID
	TokenType     string  `json:"token_type"`     //代币类型 NFISH, BUSD
	TokenBalance  float64 `json:"token_balance"`  // 玩家当前代币的余额
	TokenIncome   float64 `json:"token_income"`   // 玩家通过玩游戏的总收入
	TokenSpend    float64 `json:"token_spend"`    // 玩家通过玩游戏的总花费
	TokenDeposit  float64 `json:"token_deposit"`  // 玩家通过质押的总额度
	TokenWithdraw float64 `json:"token_withdraw"` // 玩家提取代币的总金额
	TokenFee      float64 `json:"token_fee"`      // 玩家使用当前代币付的总手续费
}

func (s ERC20TokenData) TableName() string {
	return fmt.Sprintf("t_erc20_token_data_%d", s.GameClient)
}

type ERC1155TokenData struct {
	gorm.Model `swagger-ignore:"true"`
	GameClient int    `json:"game_client"`
	AccountId  uint   `json:"account_id"` //往家账户的ID
	Ids        string `json:"ids"`        //玩家拥有的NFT所有的id
	Values     string `json:"values"`     // 玩家拥有的NFT的数量
}

func (s ERC1155TokenData) TableName() string {
	return fmt.Sprintf("t_erc1155_token_data_%d", s.GameClient)
}

//func (s ERC1155TokenData) Value() (driver.Value, error) {
//	b, err := json.Marshal(s)
//	return string(b), err
//}
//
//func (s *ERC1155TokenData) Scan(input interface{}) error {
//	return json.Unmarshal(input.([]byte), s)
//}

type walletDA0 struct{}

var WalletDAO = &walletDA0{}

func (dao *walletDA0) getWalletTableName(gameClient comm.GameClient) string {
	return fmt.Sprintf("t_wallet_%d", gameClient)
}

func (dao walletDA0) InitWallet(ctx context.Context, gameClient comm.GameClient, accountId uint, publicAddress string) bool {
	err := getDb(ctx).Transaction(func(tx1 *gorm.DB) error {
		nfishData := ERC20TokenData{
			Model:      gorm.Model{},
			GameClient: int(gameClient),
			AccountId:  accountId,
			TokenType:  comm.NFISH.String(),
		}
		nfishDataWalletLog := ERC20WalletLog{
			Model:          gorm.Model{},
			GameClient:     int(gameClient),
			AccountId:      accountId,
			BusinessModule: "Initialization",
			ActionType:     comm.Initialize.String(),
			TokenType:      comm.NFISH.String(),
			Value:          0,
			Fee:            0,
			Status:         comm.Pending.String(),
			OriginalWallet: Wallet{},
			SettledWallet:  Wallet{},
		}
		nfishWalletLog, err := ERC20WalletLogDAO.InsertERC20WalletLog(tx1, &nfishDataWalletLog)
		if err != nil {
			return err
		}

		busdData := ERC20TokenData{
			Model:      gorm.Model{},
			GameClient: int(gameClient),
			AccountId:  accountId,
			TokenType:  comm.BUSD.String(),
		}
		busdDataWalletLog := ERC20WalletLog{
			Model:          gorm.Model{},
			GameClient:     int(gameClient),
			AccountId:      accountId,
			BusinessModule: "Initialization",
			ActionType:     comm.Initialize.String(),
			TokenType:      comm.BUSD.String(),
			Value:          0,
			Fee:            0,
			Status:         comm.Pending.String(),
			OriginalWallet: Wallet{},
			SettledWallet:  Wallet{},
		}
		busdWalletLog, err := ERC20WalletLogDAO.InsertERC20WalletLog(tx1, &busdDataWalletLog)
		if err != nil {
			return err
		}

		nftData := ERC1155TokenData{
			Model:      gorm.Model{},
			GameClient: int(gameClient),
			AccountId:  accountId,
		}
		nftLog := ERC1155WalletLog{
			Model:          gorm.Model{},
			GameClient:     int(gameClient),
			AccountId:      accountId,
			BusinessModule: "Initialization",
			ActionType:     comm.Initialize.String(),
			Ids:            "",
			Values:         0,
			Fee:            0,
			Status:         comm.Pending.String(),
			OriginalWallet: Wallet{},
			SettledWallet:  Wallet{},
		}
		nftWalletLog, err := ERC1155WalletLogDAO.InsertERC1155WalletLog(tx1, &nftLog)
		if err != nil {
			return err
		}

		wallet := Wallet{
			Model:            gorm.Model{},
			GameClient:       int(gameClient),
			AccountId:        accountId,
			PublicAddress:    publicAddress,
			ERC20TokenData:   []ERC20TokenData{nfishData, busdData},
			ERC1155TokenData: nftData,
			CheckSign:        "",
		}
		if err := tx1.Create(&wallet).Error; err != nil {
			log.Errorf("%v", err)
			return err
		}

		nfishWalletLog.SettledWallet = wallet
		nfishWalletLog.Status = comm.Done.String()
		_, err = ERC20WalletLogDAO.UpdateERC20WalletLogStatus(tx1, nfishWalletLog)
		if err != nil {
			return err
		}
		busdWalletLog.SettledWallet = wallet
		busdWalletLog.Status = comm.Done.String()
		_, err = ERC20WalletLogDAO.UpdateERC20WalletLogStatus(tx1, busdWalletLog)
		if err != nil {
			return err
		}

		nftWalletLog.SettledWallet = wallet
		nftWalletLog.Status = comm.Done.String()
		_, err = ERC1155WalletLogDAO.UpdateERC1155WalletLogStatus(tx1, nftWalletLog)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return false
	}
	return true
}

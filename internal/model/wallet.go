package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
)

type Wallet struct {
	gorm.Model       `swagger-ignore:"true"`
	GameClient       int                `json:"game_client"`
	AccountId        uint               `json:"account_id" gorm:"unique;not null"` // 玩家账户ID
	PublicAddress    string             `json:"address"`                           // 玩家的钱包地址
	ERC20TokenData   []ERC20TokenWallet `json:"erc_20_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	ERC1155TokenData ERC1155TokenWallet `json:"erc_1155_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	CheckSign        string             `json:"check_sign" gorm:"type:varchar(128);not null;comment:'安全签名'"`
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

type ERC20TokenWallet struct {
	gorm.Model    `swagger-ignore:"true"`
	GameClient    int     `json:"game_client"`
	AccountId     uint    `json:"account_id"` //往家账户的ID
	Token         string  `json:"token"`      //代币类型 NFISH, BUSD
	Balance       float64 `json:"balance"`    // 玩家当前代币的余额
	Decimal       uint    `json:"decimal"`
	TotalIncome   float64 `json:"total_income"`   // 玩家通过玩游戏的总收入
	TotalSpend    float64 `json:"total_spend"`    // 玩家通过玩游戏的总花费
	TotalDeposit  float64 `json:"total_deposit"`  // 玩家通过质押的总额度
	TotalWithdraw float64 `json:"total_withdraw"` // 玩家提取代币的总金额
	TotalFee      float64 `json:"total_fee"`      // 玩家使用当前代币付的总手续费
}

func (s ERC20TokenWallet) TableName() string {
	return fmt.Sprintf("t_erc20_token_data_%d", s.GameClient)
}

type ERC1155TokenWallet struct {
	gorm.Model `swagger-ignore:"true"`
	GameClient int    `json:"game_client"`
	AccountId  uint   `json:"account_id"` //往家账户的ID
	Ids        string `json:"ids"`        //玩家拥有的NFT所有的id
	Values     string `json:"values"`     // 玩家拥有的NFT的数量
}

func (s ERC1155TokenWallet) TableName() string {
	return fmt.Sprintf("t_erc1155_token_data_%d", s.GameClient)
}

type walletDA0 struct{}

var WalletDAO = &walletDA0{}

func (dao walletDA0) GetWallet(db *gorm.DB, gameClient comm.GameClient, accountId uint) (Wallet, error) {
	var wallet Wallet
	if err := db.
		Preload("ERC20TokenData").
		Preload("ERC1155TokenData").
		Where("game_client = ? AND account_id = ?", gameClient, accountId).
		First(&wallet).Error; err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

func (dao walletDA0) UpdateWallet(db *gorm.DB, newWallet Wallet) error {
	if err := db.Save(&newWallet).Error; err != nil {
		return err
	}
	return nil
}

func (dao walletDA0) CreateWallet(db *gorm.DB, wallet Wallet) error {
	//var mysqlErr *mysql.MySQLError
	if err := db.Create(&wallet).Error; err != nil {
		//if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
		//	return errors.New("record is already existed")
		//}
		return err
	}
	return nil
}

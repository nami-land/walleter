package model

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"gorm.io/gorm"
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
	AccountId     uint    `json:"account_id"`    //往家账户的ID
	TokenType     string  `json:"token_type"`    //代币类型 NFISH, BUSD
	TokenBalance  float64 `json:"token_balance"` // 玩家当前代币的余额
	Decimal       uint    `json:"decimal"`
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

type walletDA0 struct{}

var WalletDAO = &walletDA0{}

func (dao walletDA0) GetWallet(ctx context.Context, accountId string) (*Wallet, error) {
	var wallet Wallet
	if err := GetDb(ctx).Where("account_id = ?", accountId).First(&wallet).Error; err != nil {
		return nil, err
	}
	return &wallet, nil
}

func (dao walletDA0) UpdateWallet(ctx context.Context, newWallet Wallet) error {
	if err := GetDb(ctx).Save(&newWallet).Error; err != nil {
		return err
	}
	return nil
}

func (dao walletDA0) CreateWallet(ctx context.Context, wallet Wallet) error {
	if err := GetDb(ctx).Create(&wallet).Error; err != nil {
		return err
	}
	return nil
}

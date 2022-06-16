package model

import (
	"fmt"
	"gorm.io/gorm"
)

// ERC20WalletLog 钱包流水日志
type ERC20WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	GameClient     int     `json:"game_client"`
	AccountId      uint    `json:"account_id"` //往家账户的ID
	BusinessModule string  `json:"business_module" gorm:"type:varchar(64);not null;comment:'业务模块'"`
	ActionType     string  `json:"action_type" gorm:"type:varchar(64);not null;comment:'操作类型'"`
	TokenType      string  `json:"token_type;comment:'变更的代币数据'"`
	Value          float64 `json:"value"` // 代币金额
	Fee            float64 `json:"fee"`   // 手续费
	Status         string  `json:"status" gorm:"type:varchar(64);not null;comment:处理状态"`
	OriginalWallet Wallet  `json:"original_wallet" gorm:"type:json;not null;comment:'变更前的钱包数据'"`
	SettledWallet  Wallet  `json:"settled_wallet" gorm:"type:json;not null;comment:'变更后的钱包数据'"`
}

func (s ERC20WalletLog) TableName() string {
	return fmt.Sprintf("t_erc20_wallet_log_%d", s.GameClient)
}

// ERC1155WalletLog 钱包流水日志
type ERC1155WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	GameClient     int     `json:"game_client"`
	AccountId      uint    `json:"account_id"` //往家账户的ID
	BusinessModule string  `json:"business_module" gorm:"type:varchar(64);not null;comment:'业务模块'"`
	ActionType     string  `json:"action_type" gorm:"type:varchar(64);not null;comment:'操作类型'"`
	Ids            string  `json:"ids;comment:'变更的NFT IDs'"`
	Values         float64 `json:"values"` // 变更的NFT数量
	Fee            float64 `json:"fee"`    // 手续费
	Status         string  `json:"status" gorm:"type:varchar(64);not null;comment:处理状态"`
	OriginalWallet Wallet  `json:"original_wallet" gorm:"type:json;not null;comment:'变更前的钱包数据'"`
	SettledWallet  Wallet  `json:"settled_wallet" gorm:"type:json;not null;comment:'变更后的钱包数据'"`
}

func (s ERC1155WalletLog) TableName() string {
	return fmt.Sprintf("t_erc1155_wallet_log_%d", s.GameClient)
}

type erc20WalletLogDAO struct{}

var ERC20WalletLogDAO = &erc20WalletLogDAO{}

func (s erc20WalletLogDAO) InsertERC20WalletLog(db *gorm.DB, erc20Log *ERC20WalletLog) (*ERC20WalletLog, error) {
	err := db.Create(erc20Log).Error
	if err != nil {
		return nil, nil
	}
	return erc20Log, err
}

func (s erc20WalletLogDAO) UpdateERC20WalletLogStatus(db *gorm.DB, newLog *ERC20WalletLog) (*ERC20WalletLog, error) {
	err := db.Save(&newLog).Error
	if err != nil {
		return nil, err
	}
	return newLog, nil
}

type erc1155WalletLogDAO struct{}

var ERC1155WalletLogDAO = &erc1155WalletLogDAO{}

func (s erc1155WalletLogDAO) InsertERC1155WalletLog(db *gorm.DB, erc1155Log *ERC1155WalletLog) (*ERC1155WalletLog, error) {
	err := db.Create(erc1155Log).Error
	if err != nil {
		return nil, nil
	}
	return erc1155Log, err
}

func (s erc1155WalletLogDAO) UpdateERC1155WalletLogStatus(db *gorm.DB, newLog *ERC1155WalletLog) (*ERC1155WalletLog, error) {
	err := db.Save(&newLog).Error
	if err != nil {
		return nil, err
	}
	return newLog, nil
}

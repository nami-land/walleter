package wallet_center

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

type Wallet struct {
	gorm.Model       `swagger-ignore:"true"`
	AccountId        uint64             `json:"account_id" gorm:"unique;not null"` // 玩家账户ID
	ERC20TokenData   []ERC20TokenWallet `json:"erc_20_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	ERC1155TokenData ERC1155TokenWallet `json:"erc_1155_token_data" gorm:"foreignKey:AccountId;references:AccountId"`
	CheckSign        string             `json:"check_sign" gorm:"type:varchar(128);not null;comment:'安全签名'"`
}

func (w Wallet) Value() (driver.Value, error) {
	b, err := json.Marshal(w)
	return string(b), err
}

func (w *Wallet) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), w)
}

type ERC20TokenWallet struct {
	gorm.Model    `swagger-ignore:"true"`
	AccountId     uint64  `json:"account_id"`
	Token         string  `json:"token" gorm:"type:varchar(20)"`
	Balance       float64 `json:"balance"`
	Decimal       uint    `json:"decimal"`
	TotalIncome   float64 `json:"total_income"`
	TotalSpend    float64 `json:"total_spend"`
	TotalDeposit  float64 `json:"total_deposit"`
	TotalWithdraw float64 `json:"total_withdraw"`
	TotalFee      float64 `json:"total_fee"`
}

type ERC1155TokenWallet struct {
	gorm.Model `swagger-ignore:"true"`
	AccountId  uint64 `json:"account_id"`
	Ids        string `json:"ids"`
	Values     string `json:"values"`
}

type walletDA0 struct{}

var walletDAO = &walletDA0{}

func (dao walletDA0) getWallet(db *gorm.DB, accountId uint64) (Wallet, error) {
	var w Wallet
	if err := db.
		Preload("ERC20TokenData").
		Preload("ERC1155TokenData").
		Where("account_id = ?", accountId).
		First(&w).Error; err != nil {
		return Wallet{}, err
	}
	return w, nil
}

func (dao walletDA0) updateWallet(db *gorm.DB, newWallet Wallet) error {
	if err := db.Save(&newWallet).Error; err != nil {
		return err
	}
	return nil
}

func (dao walletDA0) createWallet(db *gorm.DB, wallet Wallet) error {
	var mysqlErr *mysql.MySQLError
	if err := db.Create(&wallet).Error; err != nil {
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return nil
		}
		return err
	}
	return nil
}

func (dao walletDA0) updateWalletCheckSign(db *gorm.DB, newWallet Wallet) error {
	return db.Model(&newWallet).
		Where("account_id = ?", newWallet.AccountId).
		Update("check_sign", newWallet.CheckSign).
		Error
}

func (dao walletDA0) updateERC20WalletData(db *gorm.DB, newERC20Data ERC20TokenWallet) error {
	return db.Save(&newERC20Data).Error
}

func (dao walletDA0) updateERC1155WalletData(db *gorm.DB, newERC1155Data ERC1155TokenWallet) error {
	return db.Save(&newERC1155Data).Error
}

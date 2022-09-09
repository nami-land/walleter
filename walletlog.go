package walleter

import (
	"database/sql/driver"
	"encoding/json"

	"gorm.io/gorm"
)

// ERC20WalletLog Wallet flow log
type ERC20WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	AccountId      uint64               `json:"account_id"`
	BusinessModule string               `json:"business_module" gorm:"type:varchar(64);not null;"`
	ActionType     string               `json:"action_type" gorm:"type:varchar(64);not null;"`
	Tokens         erc20TokenCollection `json:"tokens" gorm:"type:json;not null"`
	Fees           erc20TokenCollection `json:"fees" gorm:"type:json;"`
	Status         string               `json:"status" gorm:"type:varchar(64);not null;"`
	OriginalWallet Wallet               `json:"original_wallet" gorm:"type:json;not null;"`
	SettledWallet  Wallet               `json:"settled_wallet" gorm:"type:json;not null;"`
}

// ERC1155WalletLog Wallet flow log
type ERC1155WalletLog struct {
	gorm.Model     `swagger-ignore:"true"`
	AccountId      uint64               `json:"account_id"`
	BusinessModule string               `json:"business_module" gorm:"type:varchar(64);not null;"`
	ActionType     string               `json:"action_type" gorm:"type:varchar(64);not null;"`
	Ids            string               `json:"ids"`
	Values         string               `json:"values"`
	Fees           erc20TokenCollection `json:"fees" gorm:"type:json;"`
	Status         string               `json:"status" gorm:"type:varchar(10);not null;"`
	OriginalWallet Wallet               `json:"original_wallet" gorm:"type:json;not null;"`
	SettledWallet  Wallet               `json:"settled_wallet" gorm:"type:json;"`
}

type erc20WalletLogDAO struct{}

var erc20LogDAO = &erc20WalletLogDAO{}

func (s erc20WalletLogDAO) insertERC20WalletLog(db *gorm.DB, erc20Log ERC20WalletLog) (ERC20WalletLog, error) {
	err := db.Create(&erc20Log).Error
	if err != nil {
		return ERC20WalletLog{}, err
	}
	return erc20Log, nil
}

func (s erc20WalletLogDAO) updateERC20WalletLogStatus(db *gorm.DB, newLog ERC20WalletLog) (ERC20WalletLog, error) {
	err := db.Save(&newLog).Error
	if err != nil {
		return ERC20WalletLog{}, err
	}
	return newLog, nil
}

type erc1155WalletLogDAO struct{}

var erc1155LogDAO = &erc1155WalletLogDAO{}

func (s erc1155WalletLogDAO) insertERC1155WalletLog(db *gorm.DB, erc1155Log ERC1155WalletLog) (ERC1155WalletLog, error) {
	err := db.Create(&erc1155Log).Error
	if err != nil {
		return ERC1155WalletLog{}, err
	}
	return erc1155Log, nil
}

func (s erc1155WalletLogDAO) updateERC1155WalletLogStatus(db *gorm.DB, newLog ERC1155WalletLog) (ERC1155WalletLog, error) {
	err := db.Save(&newLog).Error
	if err != nil {
		return ERC1155WalletLog{}, err
	}
	return newLog, nil
}

type erc20TokenCollection struct {
	Items []erc20TokenData `json:"items"`
}

func (item erc20TokenCollection) Value() (driver.Value, error) {
	b, err := json.Marshal(item)
	return string(b), err
}

func (item *erc20TokenCollection) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), item)
}

type erc20TokenData struct {
	TokenType string  `json:"token_type" gorm:"type:varchar(20)"`
	Amount    float64 `json:"amount"`
	Decimal   uint64  `json:"decimal"`
}

func (item erc20TokenData) Value() (driver.Value, error) {
	b, err := json.Marshal(item)
	return string(b), err
}

func (item *erc20TokenData) Scan(input interface{}) error {
	return json.Unmarshal(input.([]byte), item)
}

// /----------------------------
// Wallet log service
type walletLogService struct{}

func newWalletLogService() *walletLogService {
	return &walletLogService{}
}

// insertNewERC20WalletLog Insert new log of ERC20 changes
func (receiver *walletLogService) insertNewERC20WalletLog(db *gorm.DB, command WalletCommand, currentWallet Wallet) (ERC20WalletLog, error) {
	erc20WalletLog := parseCommandToERC20WalletLog(command, currentWallet)
	return erc20LogDAO.insertERC20WalletLog(db, erc20WalletLog)
}

// updateERC20WalletLog Change the status of ERC20Log in batches
func (receiver *walletLogService) updateERC20WalletLog(db *gorm.DB, log ERC20WalletLog, status WalletLogStatus, newWallet Wallet) (ERC20WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return erc20LogDAO.updateERC20WalletLogStatus(db, log)
}

// insertNewERC1155WalletLog Insert an ERC1155 asset change log
func (receiver *walletLogService) insertNewERC1155WalletLog(db *gorm.DB, command WalletCommand, currentWallet Wallet) (ERC1155WalletLog, error) {
	erc1155WalletData := parseCommandToERC1155WalletLog(command, currentWallet)
	return erc1155LogDAO.insertERC1155WalletLog(db, erc1155WalletData)
}

// updateERC1155WalletLog Change the state of the ERC1155 log
func (receiver *walletLogService) updateERC1155WalletLog(db *gorm.DB, log ERC1155WalletLog, status WalletLogStatus, newWallet Wallet) (ERC1155WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return erc1155LogDAO.updateERC1155WalletLogStatus(db, log)
}

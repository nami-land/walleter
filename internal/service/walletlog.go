package service

import (
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/pkg"
)

type walletLogService struct{}

func NewWalletLogService() *walletLogService {
	return &walletLogService{}
}

// InsertNewERC20WalletLog 插入新的ERC20变更的log
func (receiver *walletLogService) InsertNewERC20WalletLog(
	db *gorm.DB, command model.WalletCommand, currentWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	erc20WalletLog := pkg.ParseCommandToERC20WalletLog(command, currentWallet)
	return model.ERC20WalletLogDAO.InsertERC20WalletLog(db, erc20WalletLog)
}

// UpdateERC20WalletLog 批量变更ERC20Log的状态
func (receiver *walletLogService) UpdateERC20WalletLog(
	db *gorm.DB, log model.ERC20WalletLog, status comm.WalletLogStatus, newWallet model.Wallet,
) (model.ERC20WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return model.ERC20WalletLogDAO.UpdateERC20WalletLogStatus(db, log)
}

// InsertNewERC1155WalletLog 插入一条ERC1155资产变更log
func (receiver *walletLogService) InsertNewERC1155WalletLog(
	db *gorm.DB, command model.WalletCommand, currentWallet model.Wallet,
) (model.ERC1155WalletLog, error) {
	erc1155WalletData := pkg.ParseCommandToERC1155WalletLog(command, currentWallet)
	return model.ERC1155WalletLogDAO.InsertERC1155WalletLog(db, erc1155WalletData)
}

// UpdateERC1155WalletLog 变更ERC1155 log的状态
func (receiver *walletLogService) UpdateERC1155WalletLog(
	db *gorm.DB, log model.ERC1155WalletLog, status comm.WalletLogStatus, newWallet model.Wallet,
) (model.ERC1155WalletLog, error) {
	log.Status = status.String()
	log.SettledWallet = newWallet
	return model.ERC1155WalletLogDAO.UpdateERC1155WalletLogStatus(db, log)
}

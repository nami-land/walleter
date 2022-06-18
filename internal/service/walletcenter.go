package service

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"neco-wallet-center/internal/comm"
	"neco-wallet-center/internal/model"
	"neco-wallet-center/internal/pkg"
	"neco-wallet-center/internal/utils"
)

type walletCenterService struct{}

func NewWalletCenterService() *walletCenterService {
	return &walletCenterService{}
}

func (s *walletCenterService) HandleWalletCommand(ctx context.Context, command model.WalletCommand) (model.Wallet, error) {
	if !command.GameClient.IsSupport() {
		return model.Wallet{}, errors.New("game client is invalid")
	}

	switch command.ActionType {
	case comm.Initialize:
		return initWallet(ctx, command)
	default:
		return updateWallet(ctx, command)
	}
}

func initWallet(ctx context.Context, command model.WalletCommand) (model.Wallet, error) {
	err := model.GetDb(ctx).Transaction(func(tx1 *gorm.DB) error {
		// 1. Insert change logs, including ERC20 logs and ERC1155 Log.
		walletLogService := NewWalletLogService()
		erc20WalletLog, err := walletLogService.InsertNewERC20WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		erc115WalletLog, err := walletLogService.InsertNewERC1155WalletLog(tx1, command, model.Wallet{})
		if err != nil {
			return err
		}

		// 2. initialize user's wallet data.
		erc20DataArray := pkg.ParseCommandToERC20WalletArray(command)
		erc1155Data := pkg.ParseCommandToERC1155Wallet(command)
		wallet := model.Wallet{
			Model:            gorm.Model{},
			GameClient:       int(command.GameClient),
			AccountId:        command.AccountId,
			PublicAddress:    command.PublicAddress,
			ERC20TokenData:   erc20DataArray,
			ERC1155TokenData: erc1155Data,
			CheckSign:        "",
		}

		// 3. generator a new check sign
		newCheckSign, err := pkg.NewWalletValidator().GenerateNewSignHash(wallet)
		if err != nil {
			return err
		}
		wallet.CheckSign = newCheckSign

		err = model.WalletDAO.CreateWallet(tx1, wallet)
		if err != nil {
			return err
		}

		// 4. change log statuses
		_, err = walletLogService.UpdateERC20WalletLog(tx1, erc20WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		_, err = walletLogService.UpdateERC1155WalletLog(tx1, erc115WalletLog, comm.Done, wallet)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
}

func updateWallet(ctx context.Context, command model.WalletCommand) (model.Wallet, error) {
	switch command.AssetType {
	case comm.ERC20AssetType:
		return handleERC20Command(ctx, command)
	case comm.ERC1155AssetType:
		return handleERC1155Command(ctx, command)
	default:
		return model.Wallet{}, errors.New("not support current asset type")
	}
}

func handleERC20Command(ctx context.Context, command model.WalletCommand) (model.Wallet, error) {
	err := model.GetDb(ctx).Transaction(func(tx *gorm.DB) error {
		logService := NewWalletLogService()
		validator := pkg.NewWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		// 1. 验证用户当前钱包状态是否正常
		result, err := validator.ValidateWallet(userWallet)
		if err != nil || result == false {
			return err
		}

		// 2.插入一条log信息
		log, err := logService.InsertNewERC20WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. 是否收取手续费
		userWallet, err = NewFeeChargerService().ChargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.UpdateERC20WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 4. 对用户资产进行变更
		switch command.ActionType {
		case comm.Deposit:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalDeposit += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Withdraw:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalWithdraw += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Income:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance += token.Value
				userERC20TokenWallet.TotalIncome += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.Spend:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalSpend += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		case comm.ChargeFee:
			for _, token := range command.ERC20Commands {
				index, userERC20TokenWallet := getUserERC20TokenWallet(userWallet.ERC20TokenData, token.Token)
				userERC20TokenWallet.Balance -= token.Value
				userERC20TokenWallet.TotalFee += token.Value
				userWallet.ERC20TokenData[index] = userERC20TokenWallet
			}
			break
		default:
			return errors.New("not support action type")
		}

		// 6. 生成新的验证信息
		newCheckSign, err := validator.GenerateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign

		// 7. 更新用户资产
		err = model.WalletDAO.UpdateWallet(tx, userWallet)
		if err != nil {
			_, err = logService.UpdateERC20WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 8. 更新log信息
		_, err = NewWalletLogService().UpdateERC20WalletLog(tx, log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
}

func handleERC1155Command(ctx context.Context, command model.WalletCommand) (model.Wallet, error) {
	err := model.GetDb(ctx).Transaction(func(tx *gorm.DB) error {
		logService := NewWalletLogService()
		validator := pkg.NewWalletValidator()

		userWallet, err := model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
		if err != nil {
			return err
		}

		// 1. 验证用户当前钱包状态是否正常
		result, err := validator.ValidateWallet(userWallet)
		if err != nil || result == false {
			return err
		}

		// 2.插入一条log信息
		log, err := logService.InsertNewERC1155WalletLog(tx, command, userWallet)
		if err != nil {
			return err
		}

		// 3. 是否收取手续费
		userWallet, err = NewFeeChargerService().ChargeFee(tx, command, userWallet)
		if err != nil {
			_, err = logService.UpdateERC1155WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 5. 对用户资产进行变更
		ids := utils.ConvertStringToUIntArray(userWallet.ERC1155TokenData.Ids)
		values := utils.ConvertStringToUIntArray(userWallet.ERC1155TokenData.Values)
		switch command.ActionType {
		case comm.Deposit, comm.Income:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := utils.GetIndexFromUIntArray(ids, id)
				if i == -1 {
					ids = append(ids, id)
					values = append(values, value)
				} else {
					values[i] = values[i] + value
				}

				userWallet.ERC1155TokenData.Ids = utils.ConvertUintArrayToString(ids)
				userWallet.ERC1155TokenData.Values = utils.ConvertUintArrayToString(values)
			}
			break
		case comm.Withdraw, comm.Spend:
			for index, id := range command.ERC1155Command.Ids {
				value := command.ERC1155Command.Values[index]
				i := utils.GetIndexFromUIntArray(ids, id)
				if i == -1 {
					return errors.New("insufficient nft balance")
				} else {
					if values[i] < value {
						return errors.New("insufficient nft balance")
					}
					values[i] = values[i] - value
				}

				userWallet.ERC1155TokenData.Ids = utils.ConvertUintArrayToString(ids)
				userWallet.ERC1155TokenData.Values = utils.ConvertUintArrayToString(values)
			}
			break
		default:
			return errors.New("not support action type")
		}

		// 6. 生成新的验证信息
		newCheckSign, err := validator.GenerateNewSignHash(userWallet)
		if err != nil {
			return err
		}
		userWallet.CheckSign = newCheckSign

		// 7. 更新用户资产
		err = model.WalletDAO.UpdateWallet(tx, userWallet)
		if err != nil {
			_, err = logService.UpdateERC1155WalletLog(tx, log, comm.Failed, userWallet)
			return err
		}

		// 8. 更新log信息
		_, err = NewWalletLogService().UpdateERC1155WalletLog(tx, log, comm.Done, userWallet)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return model.Wallet{}, err
	}
	return model.WalletDAO.GetWallet(model.GetDb(ctx), command.GameClient, command.AccountId)
}

package main

import (
	"fmt"

	"github.com/neco-fun/walleter"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"root123..",
		"localhost",
		"3306",
		"walleter_test_db",
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("connect database failed")
	}

	// init a walleter instance
	w := walleter.New(db, 1)

	feeChargerWallet, err := w.GetWalletByAccountId(1)
	if err != nil {
		logrus.Fatalln(err)
	}
	fmt.Println(feeChargerWallet)

	var userId uint64 = 10
	// create a new user wallet account
	userInitCommand := walleter.NewInitWalletCommand(userId)
	userWallet, err := w.HandleWalletCommand(db, userInitCommand)
	if err != nil {
		logrus.Fatalln(err)
	}
	fmt.Println(userWallet)

	// Testing income operation
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	_, err := w.HandleWalletCommand(
	//		db,
	//		walleter.NewERC20WalletCommand(
	//			userId,
	//			walleter.Income,
	//			"Testing",
	//			walleter.InGame,
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 100.0,
	//				walleter.BUSD:  5,
	//			},
	//			map[walleter.ERC20TokenEnum]float64{},
	//		),
	//	)
	//	return err
	//})
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//
	//userWallet, err = w.GetWalletByAccountId(userId)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//fmt.Println(userWallet)

	// Testing spend operation
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	_, err := w.HandleWalletCommand(
	//		db,
	//		walleter.NewERC20WalletCommand(
	//			userId,
	//			walleter.Spend,
	//			"Testing",
	//			walleter.InGame,
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 100.0,
	//				walleter.BUSD:  5,
	//			},
	//			map[walleter.ERC20TokenEnum]float64{},
	//		),
	//	)
	//	return err
	//})
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//
	//userWallet, err = w.GetWalletByAccountId(userId)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//fmt.Println(userWallet)

	// Testing deposit operation
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	_, err := w.HandleWalletCommand(
	//		db,
	//		walleter.NewERC20WalletCommand(
	//			userId,
	//			walleter.Deposit,
	//			"Testing",
	//			walleter.BSC,
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 100.0,
	//				walleter.BUSD:  5,
	//			},
	//			map[walleter.ERC20TokenEnum]float64{},
	//		),
	//	)
	//	return err
	//})
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//
	//userWallet, err = w.GetWalletByAccountId(userId)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//fmt.Println(userWallet)

	// Testing withdraw operation
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	_, err := w.HandleWalletCommand(
	//		db,
	//		walleter.NewERC20WalletCommand(
	//			userId,
	//			walleter.Withdraw,
	//			"Testing",
	//			walleter.BSC,
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 90.0,
	//				walleter.BUSD:  4,
	//			},
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 10.0,
	//				walleter.BUSD:  1,
	//			},
	//		),
	//	)
	//	return err
	//})
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//
	//userWallet, err = w.GetWalletByAccountId(userId)
	//if err != nil {
	//	logrus.Fatalln(err)
	//}
	//fmt.Println(userWallet)

	// Testing income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC20WalletCommand(
				userId,
				walleter.Income,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 10.0,
					walleter.BUSD:  5.0,
				},
				map[walleter.ERC20TokenEnum]float64{},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	//// Testing charge fee operation
	//err = db.Transaction(func(tx *gorm.DB) error {
	//	_, err := w.HandleWalletCommand(
	//		db,
	//		walleter.NewERC20WalletCommand(
	//			userId,
	//			walleter.ChargeFee,
	//			"Testing",
	//			walleter.InGame,
	//			map[walleter.ERC20TokenEnum]float64{
	//				walleter.NFISH: 10.0,
	//				walleter.BUSD:  5.0,
	//			},
	//			map[walleter.ERC20TokenEnum]float64{},
	//		),
	//	)
	//	return err
	//})
	//if err != nil {
	//	logrus.Fatalln(err)
	//}

	// Testing erc1155 income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				userId,
				walleter.Income,
				"Testing",
				walleter.InGame,
				[]uint64{10001, 10002, 10003},
				[]uint64{1, 2, 3},
				map[walleter.ERC20TokenEnum]float64{},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	// Testing erc1155 spend operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				userId,
				walleter.Spend,
				"Testing",
				walleter.InGame,
				[]uint64{10001, 10002, 10003},
				[]uint64{1, 2, 3},
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 10.0,
					walleter.BUSD:  5.0,
				},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	// Testing erc1155 income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				userId,
				walleter.Deposit,
				"Testing",
				walleter.BSC,
				[]uint64{10001, 10002, 10003},
				[]uint64{1, 2, 3},
				map[walleter.ERC20TokenEnum]float64{},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	// Testing erc1155 income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				userId,
				walleter.Withdraw,
				"Testing",
				walleter.BSC,
				[]uint64{10001, 10002, 10003},
				[]uint64{1, 2, 3},
				map[walleter.ERC20TokenEnum]float64{},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}
}

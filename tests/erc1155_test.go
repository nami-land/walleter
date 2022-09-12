package main

import (
	"fmt"
	"github.com/neco-fun/walleter"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestERC1155Income(t *testing.T) {
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

	// Testing income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
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

	userWallet, err := w.GetWalletByAccountId(testUserId)
	if err != nil {
		logrus.Fatalln(err)
	}

	if userWallet.ERC1155TokenData.Values != "1,2,3" {
		t.Fatalf("%s testing failed", "TestERC1155Income")
	}
}

func TestERC1155Spend(t *testing.T) {
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

	// Testing income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
				walleter.Spend,
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

	userWallet, err := w.GetWalletByAccountId(testUserId)
	if err != nil {
		logrus.Fatalln(err)
	}

	if userWallet.ERC1155TokenData.Values != "0,0,0" {
		t.Fatalf("%s testing failed", "TestERC1155Spend")
	}
}

func TestERC1155Deposit(t *testing.T) {
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

	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
				walleter.Deposit,
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

	userWallet, err := w.GetWalletByAccountId(testUserId)
	if err != nil {
		logrus.Fatalln(err)
	}

	if userWallet.ERC1155TokenData.Values != "1,2,3" {
		t.Fatalf("%s testing failed", "TestERC1155Income")
	}
}

func TestERC1155Withdraw(t *testing.T) {
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

	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
				walleter.Withdraw,
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

	userWallet, err := w.GetWalletByAccountId(testUserId)
	if err != nil {
		logrus.Fatalln(err)
	}

	if userWallet.ERC1155TokenData.Values != "0,0,0" {
		t.Fatalf("%s testing failed", "TestERC1155Spend")
	}
}

func TestERC1155FeeCharge(t *testing.T) {
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

	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.Deposit,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 10.0,
					walleter.BUSD:  10.0,
				},
				map[walleter.ERC20TokenEnum]float64{},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
				walleter.Deposit,
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

	// Testing income operation
	err = db.Transaction(func(tx *gorm.DB) error {
		_, err := w.HandleWalletCommand(
			db,
			walleter.NewERC1155WalletCommand(
				testUserId,
				walleter.Withdraw,
				"Testing",
				walleter.InGame,
				[]uint64{10001, 10002, 10003},
				[]uint64{1, 2, 3},
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 10.0,
					walleter.BUSD:  10.0,
				},
			),
		)
		return err
	})
	if err != nil {
		logrus.Fatalln(err)
	}

	userWallet, err := w.GetWalletByAccountId(testUserId)
	if err != nil {
		logrus.Fatalln(err)
	}

	for _, erc20 := range userWallet.ERC20TokenData {
		if erc20.Balance != 0.0 {
			t.Fatalf("%s failed", "TestERC1155FeeCharge")
		}
	}

	if userWallet.ERC1155TokenData.Values != "0,0,0" {
		t.Fatalf("%s testing failed", "TestERC1155Spend")
	}
}

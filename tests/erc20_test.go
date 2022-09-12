package main

import (
	"fmt"
	"github.com/neco-fun/walleter"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

var testUserId uint64 = 10

func TestERC20Income(t *testing.T) {
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
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.Income,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 110.0,
					walleter.BUSD:  10,
				},
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

	for _, erc20 := range userWallet.ERC20TokenData {
		if erc20.Token == walleter.NFISH.String() {
			if erc20.Balance != 110.0 {
				t.Fatalf("%s failed", "TestERC20Income")
			}
		} else if erc20.Token == walleter.BUSD.String() {
			if erc20.Balance != 10.0 {
				t.Fatalf("%s failed", "TestERC20Income")
			}
		}
	}
}

func TestERC20Spend(t *testing.T) {
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
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.Spend,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 110.0,
					walleter.BUSD:  10,
				},
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

	for _, erc20 := range userWallet.ERC20TokenData {
		if erc20.Balance != 0.0 {
			t.Fatalf("%s failed", "TestERC20Spend")
		}
	}
}

func TestERC20Deposit(t *testing.T) {
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
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.Deposit,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 110.0,
					walleter.BUSD:  10,
				},
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

	for _, erc20 := range userWallet.ERC20TokenData {
		if erc20.Token == walleter.NFISH.String() {
			if erc20.Balance != 110.0 {
				t.Fatalf("%s failed", "TestERC20Deposit")
			}
		} else if erc20.Token == walleter.BUSD.String() {
			if erc20.Balance != 10.0 {
				t.Fatalf("%s failed", "TestERC20Deposit")
			}
		}
	}
}

func TestERC20Withdraw(t *testing.T) {
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
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.Withdraw,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 110.0,
					walleter.BUSD:  10,
				},
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

	for _, erc20 := range userWallet.ERC20TokenData {
		if erc20.Balance != 0.0 {
			t.Fatalf("%s failed", "TestERC20Withdraw")
		}
	}
}

func TestERC20FeeCharge(t *testing.T) {
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
					walleter.NFISH: 110.0,
					walleter.BUSD:  10,
				},
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
			walleter.NewERC20WalletCommand(
				testUserId,
				walleter.ChargeFee,
				"Testing",
				walleter.InGame,
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 100.0,
					walleter.BUSD:  9,
				},
				map[walleter.ERC20TokenEnum]float64{
					walleter.NFISH: 10.0,
					walleter.BUSD:  1,
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
			t.Fatalf("%s failed", "TestERC20FeeCharge")
		}
	}
}

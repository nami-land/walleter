package model

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"neco-wallet-center/internal/utils"
)

var db *gorm.DB

func InitDB(config *utils.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Database.Username,
		config.Database.Password,
		config.Database.Addr,
		config.Database.Port,
		config.Database.Database)
	var err error
	db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	return db, nil
}

func getDb(ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

func mustNoErr(res *gorm.DB) {
	if res.Error != nil {
		panic(res.Error)
	}
}

func ifExist(res *gorm.DB) bool {
	if errors.Is(res.Error, gorm.ErrRecordNotFound) {
		return false
	}
	mustNoErr(res)
	return true
}

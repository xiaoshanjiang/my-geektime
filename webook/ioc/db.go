package ioc

import (
	"github.com/xiaoshanjiang/my-geektime/webook/config"
	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}
	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}

package dao

import (
	"gorm.io/gorm"

	"github.com/xiaoshanjiang/my-geektime/webook/internal/repository/dao/article"
)

func InitTables(db *gorm.DB) error {
	return db.AutoMigrate(&User{}, &article.PublishedArticle{})
}

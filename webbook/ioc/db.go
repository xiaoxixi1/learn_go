package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"project_go/webbook/config"
	"project_go/webbook/internal/repository/dao"
)

func InitDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic("failed to connect database")
	}
	dao.InitTables(db)
	return db
}

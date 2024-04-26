package wire

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func intiDb() *gorm.DB {
	db, err := gorm.Open(mysql.Open("xxxx"))
	if err != nil {
		panic(err)
	}
	return db
}

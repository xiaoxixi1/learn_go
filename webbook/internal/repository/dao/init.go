package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) {
	// 但是这不是最佳实践。因为一般创建数据库时要DBA审核的，并且把init放在dao里面不利于dao层迁移
	// 比如从gorm切换成别的框架就不好迁移
	db.AutoMigrate(&User{})
}

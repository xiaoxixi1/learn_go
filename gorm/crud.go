package main

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm" // 安装gorm本体
)

type Product struct {
	gorm.Model
	Code  string
	Price uint
}

func main() {
	// db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	//  dsn := "username:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// 初始化DB实例
	db, err := gorm.Open(mysql.Open("root:Jike1504240602*@tcp(localhost:3306)/learn_go"))
	if err != nil {
		panic("failed to connect database")
	}
	// 使用debug模式可以打印所有的执行的sql语句
	db = db.Debug()
	// 迁移 schema
	// 初始化表结构（可选）
	db.AutoMigrate(&Product{})

	// Create
	db.Create(&Product{Code: "D42", Price: 100})

	// Read
	var product Product
	db.First(&product, 1)                 // 根据整型主键查找
	db.First(&product, "code = ?", "D42") // 查找 code 字段值为 D42 的记录

	// Update - 将 product 的 price 更新为 200
	db.Model(&product).Update("Price", 200)
	// Update - 更新多个字段
	db.Model(&product).Updates(Product{Price: 200, Code: "F42"}) // 仅更新非零值字段
	db.Model(&product).Updates(map[string]interface{}{"Price": 200, "Code": "F42"})

	// Delete - 删除 product
	db.Delete(&product, 1)
}

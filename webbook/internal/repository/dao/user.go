package dao

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{db: db}
}
func (ud *UserDao) Insert(cxt context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CTime = now
	user.UTime = now
	return ud.db.WithContext(cxt).Create(user).Error
}

/*
*

			domain.User是业务概念，不一定和数据库中的表和列完全一致
		    而dao.User则是直接映射到表里的
	        比如有些字段在数据库中是JSON格式存储的，但是在domain里面会被转化成结构体
*/
type User struct {
	Id       string `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	// 创建时间，使用UTC 0的毫秒数，时区的转换一般统一让前端转换，或者留到要传给前端时转换
	CTime int64
	UTime int64
}

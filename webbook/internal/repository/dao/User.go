package dao

import (
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

func (ud *UserDao) Insert(cxt context.Context, user User) error {
	return ud.db.WithContext(cxt).Create(user).Error
}

type User struct {
	Id       string `gorm:"primaryKey"`
	Email    string
	Password string
	CTime    int64
	UTime    int64
}

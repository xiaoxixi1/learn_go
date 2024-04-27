package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

// 预定义邮箱冲突的错误
var (
	UserDuplicateError = errors.New("用户冲突")
	RecordNotFoundErr  = gorm.ErrRecordNotFound
)

type UserDao interface {
	Insert(cxt context.Context, user *User) error
	QueryByEmail(cxt context.Context, email string) (User, error)
	Update(cxt context.Context, user User) error
	QueryById(cxt context.Context, userid int64) (User, error)
	QueryByPhone(cxt context.Context, phone string) (User, error)
}

type GormUserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) UserDao {
	return &GormUserDao{db: db}
}
func (ud *GormUserDao) Insert(cxt context.Context, user *User) error {
	now := time.Now().UnixMilli()
	user.CTime = now
	user.UTime = now
	err := ud.db.WithContext(cxt).Create(user).Error
	if msg, ok := err.(*mysql.MySQLError); ok {
		const duplicateError uint64 = 1062
		if msg.Number == 1062 {
			return UserDuplicateError
		}

	}
	return err
}

func (ud *GormUserDao) QueryByEmail(cxt context.Context, email string) (User, error) {
	var u User
	err := ud.db.WithContext(cxt).Where("email=?", email).First(&u).Error
	return u, err
}

func (ud *GormUserDao) Update(cxt context.Context, user User) error {
	return ud.db.WithContext(cxt).Model(&user).Updates(user).Error
}

func (ud *GormUserDao) QueryById(cxt context.Context, userid int64) (User, error) {
	var result User
	err := ud.db.WithContext(cxt).Model(User{Id: userid}).First(&result).Error
	return result, err
}

func (ud *GormUserDao) QueryByPhone(cxt context.Context, phone string) (User, error) {
	var result User
	err := ud.db.WithContext(cxt).Where("phone=?", phone).First(&result).Error
	return result, err
}

/*
*

			domain.User是业务概念，不一定和数据库中的表和列完全一致
		    而dao.User则是直接映射到表里的
	        比如有些字段在数据库中是JSON格式存储的，但是在domain里面会被转化成结构体
*/
type User struct {
	Id              int64          `gorm:"primaryKey,autoIncrement"`
	Email           sql.NullString `gorm:"unique"`
	Phone           sql.NullString `gorm:"unique"`
	Name            string
	Birthday        int64
	PersonalProfile string
	Password        string
	// 创建时间，使用UTC 0的毫秒数，时区的转换一般统一让前端转换，或者留到要传给前端时转换
	CTime int64
	UTime int64
}

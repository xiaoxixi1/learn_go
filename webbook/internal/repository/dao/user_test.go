package dao

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestGormUserDao_Insert(t *testing.T) {
	testCases := []struct {
		name string
		mock func(t *testing.T) *sql.DB

		ctx  context.Context
		user User

		wantErr error
	}{
		{
			name: "插入成功",
			mock: func(t *testing.T) *sql.DB {
				db, sqlMock, err := sqlmock.New()
				assert.NoError(t, err)
				sqlResult := sqlmock.NewResult(1, 1)
				// 正则匹配sql语句
				sqlMock.ExpectExec("INSERT INTO .*").WillReturnResult(sqlResult)
				return db
			},
			ctx: context.Background(),
			user: User{
				Name: "测试",
			},
		},
		{
			name: "邮箱冲突",
			mock: func(t *testing.T) *sql.DB {
				db, sqlMock, err := sqlmock.New()
				assert.NoError(t, err)
				// 正则匹配sql语句
				sqlMock.ExpectExec("INSERT INTO .*").WillReturnError(&mysqlDriver.MySQLError{Number: 1062})
				return db
			},
			ctx: context.Background(),
			user: User{
				Name: "测试",
			},
			wantErr: UserDuplicateError,
		},
		{
			name: "数据库错误",
			mock: func(t *testing.T) *sql.DB {
				db, sqlMock, err := sqlmock.New()
				assert.NoError(t, err)
				// 正则匹配sql语句
				sqlMock.ExpectExec("INSERT INTO .*").WillReturnError(errors.New("数据库错误"))
				return db
			},
			ctx: context.Background(),
			user: User{
				Name: "测试",
			},
			wantErr: errors.New("数据库错误"),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sqlDb := tc.mock(t)
			db, err := gorm.Open(mysql.New(mysql.Config{
				Conn:                      sqlDb,
				SkipInitializeWithVersion: true,
			}), &gorm.Config{
				DisableAutomaticPing:   true,
				SkipDefaultTransaction: true,
			})
			assert.NoError(t, err)
			userDao := NewUserDao(db)
			err = userDao.Insert(tc.ctx, &tc.user)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

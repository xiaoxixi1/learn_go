package domain

import "time"

type User struct {
	Id              int64
	Email           string
	Password        string
	Name            string //昵称
	Birthday        time.Time
	PersonalProfile string // 个人简介
	Phone           string
	CTime           time.Time
	UTime           time.Time
}

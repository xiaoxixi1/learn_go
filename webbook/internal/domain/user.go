package domain

type User struct {
	Id              int64
	Email           string
	Password        string
	Name            string //昵称
	Birthday        string
	PersonalProfile string // 个人简介
}

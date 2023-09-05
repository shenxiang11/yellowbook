package domain

import "time"

type User struct {
	Id         uint64
	Email      string
	Phone      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
	// 用户信息，为注册后用户补充
	Profile *Profile
}

type Profile struct {
	UserId       uint64
	Nickname     string
	Birthday     string
	Introduction string
	CreateTime   time.Time
	UpdateTime   time.Time
}

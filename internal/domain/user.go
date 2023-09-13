package domain

import "time"

type User struct {
	Id         uint64
	Email      string
	Phone      string
	Password   string
	CreateTime time.Time
	UpdateTime time.Time
	GithubId   uint64
	Profile    *Profile
}

type Profile struct {
	UserId       uint64
	Nickname     string
	Birthday     string
	Introduction string
	CreateTime   time.Time
	UpdateTime   time.Time
}

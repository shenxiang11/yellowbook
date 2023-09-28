package domain

import (
	"github.com/shenxiang11/yellowbook-proto/proto"
	"time"
)

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
	Avatar       string
	Gender       proto.Gender
	CreateTime   time.Time
	UpdateTime   time.Time
}

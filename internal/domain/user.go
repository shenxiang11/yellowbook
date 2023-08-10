package domain

type User struct {
	Id       uint64
	Email    string
	Password string
	// 用户信息，为注册后用户补充
	Profile *Profile
}

type Profile struct {
	UserId       uint64
	Nickname     string
	Birthday     string
	Introduction string
}

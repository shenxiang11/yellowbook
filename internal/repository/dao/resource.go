package dao

type Resource struct {
	Id         uint64 `gorm:"primaryKey,autoIncrement"`
	Url        string `gorm:"unique"`
	Mimetype   string
	CreateTime int64
	UpdateTime int64
}

package ioc

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"yellowbook/config"
	"yellowbook/internal/repository/dao"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Conf.DB.DSN))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}

	return db
}

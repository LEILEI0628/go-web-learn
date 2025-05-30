package wire

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitTestDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("dsn"))
	if err != nil {
		panic(err)
	}
	return db
}

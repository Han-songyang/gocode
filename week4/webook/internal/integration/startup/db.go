package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:3308)/webook"), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	return db
}

package mdb

import (
	"github.com/peterouob/golang_template/tools"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitMysql() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:123456@tcp(127.0.0.1:3306)/orm?charset=utf8"), &gorm.Config{})
	tools.HandelError("error in init mysql", err)
	err = db.AutoMigrate(UserModel{})
	tools.HandelError("error in auto migrate", err)
	return db
}

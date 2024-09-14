package dao

import (
	"fmt"
	"ginProject/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// 因为要和数据库建立连接，所以我们要把链接保存到var中
var (
	Db  *gorm.DB
	err error
)

func Init() {
	Db, err = gorm.Open(mysql.Open(config.Mysqldb), &gorm.Config{})
	if err != nil {
		fmt.Println("-----mysql connect error:", err)
	}
	if Db.Error != nil {
		fmt.Println("-----database err:", err)
	}
	//Db.DB().SetMaxIdleConns(10)
	//Db.DB().SetMaxOpenConns(100)
	//Db.DB().SetConnMaxLifetime(time.Hour)

}

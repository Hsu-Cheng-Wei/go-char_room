package databases

import (
	"chatRoom/enviroment"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var MysqlOrm gorm.DB

func init() {
	db, err := gorm.Open(mysql.Open(enviroment.MysqlCon), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	MysqlOrm = *db
}

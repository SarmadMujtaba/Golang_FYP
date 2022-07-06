package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

var Conn *gorm.DB

func Connection() {
	const (
		host     = "localhost"
		port     = 3306
		user     = "root"
		password = "DummySQL786"
		dbname   = "db"
	)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True", user, password, host, port, dbname)

	var err error
	Conn, err = gorm.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection Established...")

}

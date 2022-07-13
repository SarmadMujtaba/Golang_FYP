package db

import (
	"PostJson/structures"
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
	Conn.Model(&structures.Organizations{}).AddIndex("org_id", "org_id")

	Conn.Model(&structures.Memberships{}).AddForeignKey("u_id", "users(id)", "RESTRICT", "RESTRICT")
	Conn.Model(&structures.Memberships{}).AddForeignKey("org_id", "organizations(org_id)", "RESTRICT", "RESTRICT")

	fmt.Println("Connection Established...")

}

package db

import (
	"database/sql"
	"fmt"
)

var Conn *sql.DB

func Connection() {
	const (
		host     = "localhost"
		port     = 3306
		user     = "root"
		password = "DummySQL786"
		dbname   = "db"
	)
	connString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s", user, password, host, port, dbname)

	var err error
	Conn, err = sql.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connection Established...")

}

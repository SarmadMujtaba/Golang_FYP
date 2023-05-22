package db

import (
	"PostJson/structures"
	"fmt"
	"os"

	"github.com/jinzhu/gorm"
)

var Conn *gorm.DB

func Connection() {
	// var cats []structures.Category
	// var cat structures.Category
	const (
		host   = "localhost"
		port   = 3306
		user   = "root"
		dbname = "go_db"
	)
	password := os.Getenv("DB_PASSWORD")

	// for running without docker compose
	// connString := fmt.Sprintf("%s:%s@tcp(%s:%v)/%s?charset=utf8&parseTime=True", user, password, host, port, dbname)

	connString := fmt.Sprintf("%s:%s@tcp(db)/%s?charset=utf8&parseTime=True", user, password, dbname)
	fmt.Println(connString)
	var err error
	Conn, err = gorm.Open("mysql", connString)
	if err != nil {
		panic(err)
	}

	// Migrating structures to Mysql tables
	Conn.AutoMigrate(&structures.Users{}, &structures.Organizations{}, &structures.Memberships{})
	// Conn.AutoMigrate(&structures.Category{}, &structures.Jobs{}, &structures.RequiredSkills{})
	Conn.AutoMigrate(&structures.Experience{}, &structures.Skills{}, &structures.Profile{}, structures.Applications{})
	Conn.AutoMigrate(&structures.Invites{}, &structures.Jobs{}, &structures.RequiredSkills{})

	Conn.Model(&structures.Organizations{}).AddIndex("org_id", "org_id")

	// Adding foreign Keys

	Conn.Model(&structures.Memberships{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.Memberships{}).AddForeignKey("org_id", "organizations(org_id)", "CASCADE", "CASCADE")

	Conn.Model(&structures.Jobs{}).AddForeignKey("org_id", "organizations(org_id)", "CASCADE", "CASCADE")
	// Conn.Model(&structures.Jobs{}).AddForeignKey("cat_id", "categories(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.RequiredSkills{}).AddForeignKey("job_id", "jobs(id)", "CASCADE", "CASCADE")

	Conn.Model(&structures.Experience{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.Skills{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.Profile{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")

	Conn.Model(&structures.Applications{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.Applications{}).AddForeignKey("job_id", "jobs(id)", "CASCADE", "CASCADE")

	Conn.Model(&structures.Invites{}).AddForeignKey("u_id", "users(id)", "CASCADE", "CASCADE")
	Conn.Model(&structures.Invites{}).AddForeignKey("org_id", "organizations(org_id)", "CASCADE", "CASCADE")

	// Populating categories for the First time execution
	// Conn.Find(&cats)
	// if len(cats) == 0 {
	// 	cat.ID = "1"
	// 	cat.Type = "Full-Time"
	// 	Conn.Create(&cat)
	// 	cat.ID = "2"
	// 	cat.Type = "Part-Time"
	// 	Conn.Create(&cat)
	// 	cat.ID = "3"
	// 	cat.Type = "Internship"
	// 	Conn.Create(&cat)
	// }

	fmt.Println("Connection Established...")
}

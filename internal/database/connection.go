package database

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

var DATABASE *sqlx.DB

func InitDB() error {
	if DATABASE != nil {
		return nil
	}
	userName := os.Getenv("DB_USER")
	host := os.Getenv("DB_HOST")
	passWord := os.Getenv("DB_PASSWORD")
	name := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", userName, passWord, host, port, name)
	database, err := sqlx.Connect("mysql", connectionString)
	if err != nil {
		return err
	}
	DATABASE = database
	return nil
}

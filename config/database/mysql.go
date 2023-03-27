package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type MysqlConfig struct {
	DatabaseName string
	Username string
	Password string
	Host string
	Port int
}

func NewMysqlDB(config MysqlConfig) *sql.DB {
	connection := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true", config.Username, config.Password, config.Host, config.Port, config.DatabaseName)
	db, err := sql.Open("mysql", connection)
	if err != nil {
		log.Fatal(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}
	
	return db
}
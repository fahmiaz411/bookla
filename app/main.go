package main

import (
	"database/sql"
	"os"
	"strconv"
	"time"

	"github.com/fahmiaz411/bookla/config"
	"github.com/fahmiaz411/bookla/config/database"
	"github.com/fahmiaz411/bookla/helper/constant"
	_authHandler "github.com/fahmiaz411/bookla/modules/auth/v1/delivery"
	_authRepo "github.com/fahmiaz411/bookla/modules/auth/v1/repository"
	_authUsecase "github.com/fahmiaz411/bookla/modules/auth/v1/usecase"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
)

func DBInit(db *sql.DB) {
	db.Exec(`
		CREATE TABLE activities (
			activity_id int auto_increment primary key,
			title varchar(200),
			email varchar(200),
			created_at datetime default current_timestamp,
			updated_at datetime default current_timestamp
		)
	`)

	db.Exec(`
		CREATE TABLE todos (
			todo_id int auto_increment primary key,
			activity_group_id int not null,
			title varchar(200),
			is_active bool default 1,
			priority enum('very-high') default 'very-high',
			created_at datetime default current_timestamp,
			updated_at datetime default current_timestamp,
			constraint fk_ag_id foreign key (activity_group_id) references activities(activity_id)
		)
	`)
}

func main() {

	config.New()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		WriteBufferSize: 4096,
		// IdleTimeout: 1000000,
		EnablePrintRoutes: true,
	})

	mysqlPort, _ := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if mysqlPort == constant.ZeroValue {
		mysqlPort = 3306
	}
	
	// db := database.NewMysqlDB(database.MysqlConfig{
	// 	DatabaseName: os.Getenv("MYSQL_DBNAME"),
	// 	Username: os.Getenv("MYSQL_USER"),
	// 	Password: os.Getenv("MYSQL_PASSWORD"),
	// 	Host: os.Getenv("MYSQL_HOST"),
	// 	Port: mysqlPort,
	// })

	// Dev

	db := database.NewMysqlDB(database.MysqlConfig{
		DatabaseName: "bookla",
		Username: "root",
		Password: "1234",
		Host: "localhost",
		Port: 3306,
	})

	// DBInit(db)

	route := app.Group("/api")

	timeout := time.Duration(1 * time.Minute)

	authRepo := _authRepo.NewRepository(db)
	authUsecase := _authUsecase.NewUsecase(authRepo, timeout)
	_authHandler.NewRESTHandler(route, authUsecase)

	// go writer.StartScheduling()

	app.Listen(":8000")
}
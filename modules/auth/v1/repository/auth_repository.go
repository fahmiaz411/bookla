package repository

import (
	"database/sql"

	"github.com/fahmiaz411/bookla/modules/auth/v1/interfaces"
	"github.com/fahmiaz411/bookla/modules/auth/v1/repository/mysql"
)

type Repository struct {
	MySQL interfaces.AuthRepoMysql
}

// NewRepository constructor
func NewRepository(mysqlConn *sql.DB) *Repository {
	return &Repository{
		MySQL: mysql.NewMysqlRepository(mysqlConn),
	}
}
package usersdb

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/SerhiiKhyzhko/bookstore_users-api/v2/config"
)

func NewClient(cfg config.DbConfig) (*sql.DB, error) {
	conectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
	cfg.DbUser,
	cfg.DbPassword,
	cfg.DbHost,
	cfg.DbPort,
	cfg.DbName,
	)// ?charset=utf8 is optional command
	client, err := sql.Open("mysql", conectionString)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(); err != nil {
		return nil, err
	}

	log.Println("database successfuly configured")
	return client, nil
}
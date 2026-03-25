package usersdb

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/SerhiiKhyzhko/bookstore_users-api/config"
	"github.com/joho/godotenv"
)

func NewClient(cfg config.DbConfig) (*sql.DB, error) {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

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
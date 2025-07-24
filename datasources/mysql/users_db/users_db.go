package usersdb

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/go-sql-driver/mysql"
)

var Client *sql.DB

func init() {
		err := godotenv.Load()
		if err != nil {
			panic(err)
		}

	conectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8",
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"),
	os.Getenv("DB_NAME"),
	)// ?charset=utf8 is optional command
	Client, err = sql.Open("mysql", conectionString)
	if err != nil {
		panic(err)
	}

	if err = Client.Ping(); err != nil {
		panic(err)
	}

	log.Println("database successfuly configured")
}
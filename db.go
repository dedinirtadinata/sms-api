package main

import (
	"database/sql"
	"log"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func InitDB() {
	var err error

	dsn := "host=postgres-gammu user=smsuser password=smspassword dbname=smsdb sslmode=disable"

	for i := 1; i <= 20; i++ {
		log.Println("Connecting to DB attempt", i)

		db, err = sql.Open("pgx", dsn)
		if err == nil {
			err = db.Ping()
		}

		if err == nil {
			log.Println("PostgreSQL connected")
			return
		}

		log.Println("DB not ready:", err)
		time.Sleep(3 * time.Second)
	}

	log.Fatal("FATAL: cannot connect to database")
}

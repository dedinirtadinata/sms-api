package main

import (
	"database/sql"

	_ "github.com/jackc/pgx/v5/stdlib"
)

var db *sql.DB

func InitDB() {
	conn := "host=postgres-gammu user=smsuser password=smspassword dbname=smsdb sslmode=disable"
	db, _ = sql.Open("postgres", conn)
}

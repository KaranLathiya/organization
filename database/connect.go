package database

import (
	"database/sql"
	"log"
	"organization/config"

	_ "github.com/lib/pq"
)

func Connect() *sql.DB {
	connection_string := "postgresql://" + config.ConfigVal.Database.DBUser + "@localhost:26257/" + config.ConfigVal.Database.DBName + "?sslmode=disable"
	db, err := sql.Open("postgres", connection_string)
	if err != nil {
		log.Fatalf("Database Connection err: %v", err)
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging the database:", err)
		panic(err)
	}
	return db
}

package utils

import (
	"database/sql"
	"fmt"
)

var (
	dbHost     = DotEnv("DB_HOST")
	dbPort     = DotEnv("DB_PORT")
	dbUser     = DotEnv("DB_USER")
	dbPassword = DotEnv("DB_PASSWORD")
	dbName     = DotEnv("DB_NAME")
	sslMode    = DotEnv("SSL_MODE")
)

// db connection
var dataSourceName = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

func ConnectDB() *sql.DB {
	db, err := sql.Open("postgres", dataSourceName)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Successfully connected")
	}
	return db
}

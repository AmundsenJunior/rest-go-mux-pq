package main

import (
	"os"
)

var (
	dbUser = os.Getenv("DB_USER")
	dbPassword = os.Getenv("DB_PASSWORD")
	dbName = os.Getenv("DB_NAME")
	dbHost = os.Getenv("DB_HOST")
	dbPort = os.Getenv("DB_PORT")
	dbSSLMode = os.Getenv("DB_SSLMODE")
)

func main() {
	a := App{}

	// use db connection credentials stored in env vars
	a.Initialize(dbUser, dbPassword, dbName, dbHost, dbPort, dbSSLMode)

	a.Run(":8080")
}

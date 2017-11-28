package main

import (
	"os"
)

func main() {
	a := App{}

	// use db connection creds stored in env vars
	a.Initialize(
		os.Getenv("APP_DB_USERNAME"),
		os.Getenv("APP_DB_PASSWORD"),
		os.Getenv("APP_DB_NAME"),
		os.Getenv("APP_DB_SSLMODE"),
	)

	a.Run(":8080")
}

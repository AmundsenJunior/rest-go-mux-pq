package main

import (
	"database/sql"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
)

//struct to hold refs of router and database
type App struct {
	Router *mux.Router
	DB     *sql.DB
	Logger *log.Logger
}

// create database connection, set up routing and logging
func (a *App) Initialize(user, password, dbname, host, port, sslmode string) {
	a.Logger = log.New(os.Stdout, "", log.LstdFlags)

	dsn := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=%s", user, password, dbname, host, port, sslmode)

	var err error
	a.DB, err = sql.Open("postgres", dsn)
	if err != nil {
		a.Logger.Fatal(err)
	}

	err = a.DB.Ping()
	if err != nil {
		a.Logger.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// run application
// TODO: add logging entry for host/address of application
func (a *App) Run(addr string) {
	loggedRouter := a.createLoggingRouter(a.Logger.Writer())
	a.Logger.Fatal(http.ListenAndServe(addr, loggedRouter))

	defer a.DB.Close()
}

// initialize routes into router that call methods on requests
func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/health", a.healthStatus).Methods("GET")
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/product/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")
}



package main

import (
	"net/http"
	"fmt"
	"github.com/gorilla/mux"
	"strconv"
	"database/sql"
	"encoding/json"
)

// send a payload of JSON content
func (a *App) respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

// send a JSON error message
func (a *App) respondWithError(w http.ResponseWriter, code int, message string) {
	a.respondWithJSON(w, code, map[string]string{"error": message})

	a.Logger.Printf("App error: code %d, message %s", code, message)
}

// provide health status check
func (a *App) healthStatus(w http.ResponseWriter, r *http.Request) {
	dbStatus := "OK"
	err := a.DB.Ping()
	if err != nil {
		dbStatus = fmt.Sprintf("DB access error: %s", err)

	}

	healthStatus := struct{DbStatus string `json:"dbStatus"`}{dbStatus}
	a.respondWithJSON(w, http.StatusOK, healthStatus)
}

// handle get product request
func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}

	p := product{ID: id}
	if err := p.getProduct(a.DB); err != nil {
		switch err {
		case sql.ErrNoRows:
			msg := fmt.Sprintf("Product not found. Error: %s", err.Error())
			a.respondWithError(w, http.StatusNotFound, msg)
		default:
			a.respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	a.respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, _ := strconv.Atoi(r.FormValue("count"))
	start, _ := strconv.Atoi(r.FormValue("start"))

	if count > 10 || count < 1 {
		count = 10
	}
	if start < 0 {
		start = 0
	}

	products, err := getProducts(a.DB, start, count)
	if err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, products)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	var p product
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	if err := p.createProduct(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	var p product

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	p.ID = id

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&p); err != nil {
		msg := fmt.Sprintf("Invalid request payload. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	defer r.Body.Close()

	if err := p.updateProduct(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	a.respondWithJSON(w, http.StatusOK, p)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	var p product

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		msg := fmt.Sprintf("Invalid product ID. Error: %s", err.Error())
		a.respondWithError(w, http.StatusBadRequest, msg)
		return
	}
	p.ID = id

	if err := p.deleteProduct(a.DB); err != nil {
		a.respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	a.respondWithJSON(w, http.StatusOK, p)
}
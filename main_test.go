package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

// table creation query for database setup
const (
	tableCreationQuery = `CREATE TABLE IF NOT EXISTS products (
		id SERIAL,
		name TEXT NOT NULL,
		price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
		CONSTRAINT products_pkey PRIMARY KEY (id)
	)`
	dataDeletionQuery = "DELETE FROM products"
	sequenceRestartQuery = "ALTER SEQUENCE products_id_seq RESTART WITH 1"
	addProductQuery = "INSERT INTO products(name, price) VALUES($1, $2)"
	getProductByIDQuery = "SELECT name, price FROM products WHERE id = $1"
)

// global var that represents application we want to test
var a App

// create table if it doesn't exist
func ensureTableExists() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec(dataDeletionQuery)
	a.DB.Exec(sequenceRestartQuery)
}

// execute an HTTP request and return the response
func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

// check HTTP response code actual against expected
func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code is %d. Got %d", expected, actual)
	}
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 1; i <= count; i++ {
		a.DB.Exec(addProductQuery, "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}

func getProductByID(id int) (string, float64) {
	var name string
	var price float64

	err := a.DB.QueryRow(getProductByIDQuery, id).Scan(&name, &price)
	if err != nil {
		log.Fatal(err)
	}

	return name, price
}

func TestMain(m *testing.M) {
	a = App{}

	// assume that env vars are present for db credentials
	a.Initialize(
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"),
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_SSLMODE"),
	)

	// create table if it doesn't exist
	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

// test that health status works
func TestHealthStatus(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	dbStatus := m["dbStatus"]
	if dbStatus != "OK" {
		t.Errorf("Expected an OK status for DB. Got %v", dbStatus)
	}}

// test that an empty table returns empty products
func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	// check that response has an OK status code
	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

// GET a product that doesn't exist and verify a not found response code
func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/product/11", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)

	errorMsg := "Product not found. Error: sql: no rows in result set"
	if m["error"] != errorMsg {
		t.Errorf("Expected the 'error' key of the response to be set to %s. Got '%s'", errorMsg, m["error"])
	}
}

// POST a new product create and verify the product exists correctly in the response body
func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name": "test product","price": 11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	// response body should contain the JSON of the product just created
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshalling converts ints to floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

// add a product directly to db and test that its GET returns OK
func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

// directly create a number of products in the db, and test that GET returns OK with correct number of products
// TODO: correct test check step to verify product IDs 1-10 exist, instead of just 10 items
func TestGetProducts(t *testing.T) {
	clearTable()
	addProducts(10)

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m []map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if len(m) != 10 {
		t.Errorf("Expected 10 products to return. Got %v", len(m))
	}
}

// add a product directly to db and test that a PUT successfully applies to alter product's values
func TestUpdateProduct(t *testing.T) {
	clearTable()
	originalID := 1
	addProducts(originalID)

	originalName, originalPrice := getProductByID(originalID)

	payload := []byte(`{"name": "test product - updated name", "price": 11.22}`)

	req, _ := http.NewRequest("PUT", "/product/1", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	// comparing against JSON-decoded value, which is float64 by default
	if m["id"] != float64(originalID) {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalID, m["id"])
	}

	if m["name"] == originalName {
		t.Errorf("Exected the name to change from '%v' to '%v'. Got '%v'", originalName, "test product - updated name", m["name"])
	}

	if m["price"] == originalPrice {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalPrice, 11.22, m["price"])
	}
}

// add product directly to db and test that deleting a product from application is successful
// WARNING: requires GET functionality to verify DELETE test passed
// TODO: replace GET function calls with direct DB queries, pre- and post-test
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	// verify that product exists
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// delete product from application
	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	// verify that product is deleted
	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}


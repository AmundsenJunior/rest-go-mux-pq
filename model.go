package main

import (
	"database/sql"
)

// the product type is the fields that cover each item in our database
// struct fields include encoded JSON key names
type product struct {
	ID    int     `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

const (
	getProductQuery    = "SELECT name, price FROM products WHERE id=$1"
	getProductsQuery   = "SELECT id, name, price FROM products LIMIT $1 OFFSET $2"
	updateProductQuery = "UPDATE products SET name=$1, price=$2 WHERE id=$3"
	deleteProductQuery = "DELETE FROM products WHERE id=$1"
	createProductQuery = "INSERT INTO products(name, price) VALUES($1, $2) RETURNING id"
)

// the set of functions for CRUD operations on our model
func (p *product) getProduct(db *sql.DB) error {
	return db.QueryRow(getProductQuery, p.ID).Scan(&p.Name, &p.Price)
}

func (p *product) updateProduct(db *sql.DB) error {
	_, err := db.Exec(updateProductQuery, p.Name, p.Price, p.ID)
	return err
}

func (p *product) deleteProduct(db *sql.DB) error {
	_, err := db.Exec(deleteProductQuery, p.ID)
	return err
}

func (p *product) createProduct(db *sql.DB) error {
	err := db.QueryRow(createProductQuery, p.Name, p.Price).Scan(&p.ID)
	if err != nil {
		return err
	}

	return nil
}

// separate function to get multiple products
func getProducts(db *sql.DB, start, count int) ([]product, error) {
	rows, err := db.Query(getProductsQuery, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []product{}

	for rows.Next() {
		var p product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}


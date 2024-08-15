package main

import (
	"database/sql"

	"github.com/google/uuid"

	_ "github.com/go-sql-driver/mysql"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(id, name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {
	db, error := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")

	if error != nil {
		panic(error)
	}

	defer db.Close()

	createTable(db)

	product := NewProduct("1", "Product 1", 100.0)
	error = insertProduct(db, product)

	if error != nil {
		panic(error)
	}

}

func insertProduct(db *sql.DB, product *Product) error {

	stmt, err := db.Prepare("INSERT INTO products (id, name, price) VALUES (?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, error := stmt.Exec(product.ID, product.Name, product.Price)

	if error != nil {
		return error
	}

	return nil

}

func createTable(db *sql.DB) error {

	_, error := db.Exec(`CREATE TABLE products (
		id VARCHAR(36) PRIMARY KEY,
		name VARCHAR(100),
		price DECIMAL(10, 2)
	)`)

	if error != nil {
		return error
	}

	return nil

}

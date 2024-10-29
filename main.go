package main

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

type Product struct {
	ID    string
	Name  string
	Price float64
}

func NewProduct(name string, price float64) *Product {
	return &Product{
		ID:    uuid.New().String(),
		Name:  name,
		Price: price,
	}
}

func main() {

	//ctx := context.Background()

	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/goexpert")

	if err != nil {
		panic(err)
	}

	defer db.Close()

	product := NewProduct("Product 1", 100.0)

	err = insertProduc(db, product)

	if err != nil {
		panic(err)
	}

	product.Price = 200.0

	err = updateProduct(db, product)
	if err != nil {
		panic(err)
	}

	//p, err := selectProduct(ctx, db, product.ID)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Print("Product ID: ", p.ID, " Name: ", p.Name, " Price: ", p.Price, "\n")

	products, err := selectAllProducts(db)

	if err != nil {
		panic(err)
	}

	for _, p := range products {
		fmt.Printf("Product ID: %s Name: %s Price: %.2f\n", p.ID, p.Name, p.Price)
	}

	err = deleteProduct(db, product.ID)

	if err != nil {
		panic(err)
	}

}

func insertProduc(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("INSERT INTO products(id, name, price) VALUES(?, ?, ?)")

	if err != nil {
		return err
	}

	defer stmt.Close()
	_, err = stmt.Exec(product.ID, product.Name, product.Price)

	if err != nil {
		return err
	}

	return nil
}

func updateProduct(db *sql.DB, product *Product) error {
	stmt, err := db.Prepare("UPDATE products SET name = ?, price = ? WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(product.Name, product.Price, product.ID)

	if err != nil {
		return err
	}

	return nil
}

func selectProduct(ctx context.Context, db *sql.DB, id string) (*Product, error) {
	stmt, err := db.Prepare("SELECT id, name, price FROM products WHERE id = ?")

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	var p Product

	//err = stmt.QueryRow(id).Scan(&p.ID, &p.Name, &p.Price)
	err = stmt.QueryRowContext(ctx, id).Scan(&p.ID, &p.Name, &p.Price)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func selectAllProducts(db *sql.DB) ([]Product, error) {

	rows, err := db.Query("SELECT id, name, price FROM products")

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var products []Product

	for rows.Next() {
		var p Product
		err = rows.Scan(&p.ID, &p.Name, &p.Price)
		if err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	return products, nil
}

func deleteProduct(db *sql.DB, id string) error {
	stmt, err := db.Prepare("DELETE FROM products WHERE id = ?")

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

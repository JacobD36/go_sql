package storage

import (
	"database/sql"
	"fmt"

	"github.com/jacobd39/edteam/go_sql/pkg/product"
)

const (
	mySQLMigrateProduct = `CREATE TABLE IF NOT EXISTS products(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		nombre VARCHAR(25) NOT NULL,
		observations VARCHAR(100),
		price INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	)`

	mySQLCreateProduct = `INSERT INTO products(nombre, observations, price, created_at) VALUES(?, ?, ?, ?)`

	mySQLGetAllProduct = `SELECT id, nombre, observations, price, created_at, updated_at FROM products`

	mySQLGetProductByID = mySQLGetAllProduct + " WHERE id = ?"

	mySQLUpdateProduct = `UPDATE products SET nombre = ?, observations = ?, price = ?, updated_at = ? WHERE id = ?`

	mySQLDeleteProduct = `DELETE FROM products WHERE id = ?`
)

//mySQLProduct is used for work with MySQL - product
type mySQLProduct struct {
	db *sql.DB
}

//newPsqlProduct return a new pointer of MySQLProduct
func newMySQLProduct(db *sql.DB) *mySQLProduct {
	return &mySQLProduct{db}
}

//Migrate implements the interface product.Storage
func (p *mySQLProduct) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateProduct)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	fmt.Println("Migración de producto ejecutada con éxito")
	return nil
}

//Create implements the interface product.Storage
func (p *mySQLProduct) Create(m *product.Model) error {
	stmt, err := p.db.Prepare(mySQLCreateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		m.Name,
		stringToNull(m.Observations),
		m.Price,
		m.CreatedAt,
	)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	m.ID = uint(id)

	fmt.Printf("Creación de producto ejecutada con éxito. ID: %d", m.ID)
	return nil
}

//GetAll implements the interface product.Storage
func (p *mySQLProduct) GetAll() (product.Models, error) {
	stmt, err := p.db.Prepare(mySQLGetAllProduct)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ms := make(product.Models, 0)
	for rows.Next() {
		m, err := scanRowProduct(rows)
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return ms, nil
}

//GetByID implements the interface product.Storage
func (p *mySQLProduct) GetByID(id uint) (*product.Model, error) {
	stmt, err := p.db.Prepare(mySQLGetProductByID)
	if err != nil {
		return &product.Model{}, err
	}

	defer stmt.Close()

	return scanRowProduct(stmt.QueryRow(id))
}

//Update implements the interface product.Storage
func (p *mySQLProduct) Update(m *product.Model) error {
	stmt, err := p.db.Prepare(mySQLUpdateProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(m.Name, stringToNull(m.Observations), m.Price, timeToNull(m.UpdatedAt), m.ID)

	if err != nil {
		return err
	}

	fmt.Println("Se actualizó el producto correctamente")
	return nil
}

//Delete implements the interface product.Storage
func (p *mySQLProduct) Delete(id uint) error {
	stmt, err := p.db.Prepare(mySQLDeleteProduct)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	fmt.Println("El registro se eliminó exitosamente")
	return nil
}

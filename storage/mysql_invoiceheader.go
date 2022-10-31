package storage

import (
	"database/sql"
	"fmt"

	"github.com/jacobd39/edteam/go_sql/pkg/invoiceheader"
)

const (
	mySQLMigrateInvoiceHeader = `CREATE TABLE IF NOT EXISTS invoice_headers(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		client VARCHAR(100) NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now()
	)`

	mySQLCreateInvoiceHeader = `INSERT INTO invoice_headers(client) VALUES(?)`
)

//mysqlInvoiceHeader is used for work with MySQL - invoiceheader
type mySQLInvoiceHeader struct {
	db *sql.DB
}

//newMySQLInvoiceHeader return a new pointer of MySQLInvoiceHeader
func newMySQLInvoiceHeader(db *sql.DB) *mySQLInvoiceHeader {
	return &mySQLInvoiceHeader{db}
}

//Migrate implements the interface invoiceheader.Storage
func (p *mySQLInvoiceHeader) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateInvoiceHeader)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	fmt.Println("Migración de InvoiceHeader ejecutada con éxito")
	return nil
}

//CreateTx implements the interface invoiceheader.Storage
func (p *mySQLInvoiceHeader) CreateTx(tx *sql.Tx, m *invoiceheader.Model) error {
	stmt, err := tx.Prepare(mySQLCreateInvoiceHeader)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(m.Client)

	if err != nil {
		return err
	}

	id, err := result.LastInsertId()

	if err != nil {
		return err
	}

	m.ID = uint(id)

	return nil
}

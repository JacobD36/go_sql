package storage

import (
	"database/sql"
	"fmt"

	"github.com/jacobd39/edteam/go_sql/pkg/invoiceitem"
)

const (
	mySQLMigrateInvoiceItem = `CREATE TABLE IF NOT EXISTS invoice_items(
		id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
		invoice_header_id INT NOT NULL,
		product_id INT NOT NULL,
		created_at TIMESTAMP NOT NULL DEFAULT now(),
		updated_at TIMESTAMP DEFAULT now(),
		CONSTRAINT invoice_items_invoice_header_id_fk FOREIGN KEY (invoice_header_id) REFERENCES invoice_headers (id) ON UPDATE RESTRICT ON DELETE RESTRICT,
		CONSTRAINT invoice_items_product_id_fk FOREIGN KEY (product_id) REFERENCES products (id) ON UPDATE RESTRICT ON DELETE RESTRICT
	)`

	mySQLCreateInvoiceItem = `INSERT INTO invoice_items(invoice_header_id, product_id) VALUES(?, ?)`
)

//mySQLInvoiceItem is used for work with MySQL - invoiceheader
type mySQLInvoiceItem struct {
	db *sql.DB
}

//newMySQLInvoiceItem return a new pointer of PsqlInvoiceHeader
func newMySQLInvoiceItem(db *sql.DB) *mySQLInvoiceItem {
	return &mySQLInvoiceItem{db}
}

//Migrate implements the interface invoiceitem.Storage
func (p *mySQLInvoiceItem) Migrate() error {
	stmt, err := p.db.Prepare(mySQLMigrateInvoiceItem)
	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec()

	if err != nil {
		return err
	}

	fmt.Println("Migración de InvoiceItem ejecutada con éxito")
	return nil
}

//CreateTx implements the interface invoiceitem.Storage
func (p *mySQLInvoiceItem) CreateTx(tx *sql.Tx, headerID uint, ms invoiceitem.Models) error {
	stmt, err := tx.Prepare(mySQLCreateInvoiceItem)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, item := range ms {
		result, err := stmt.Exec(headerID, item.ProductID)

		if err != nil {
			return err
		}

		id, err := result.LastInsertId()

		if err != nil {
			return err
		}

		item.ID = uint(id)
	}

	return nil
}

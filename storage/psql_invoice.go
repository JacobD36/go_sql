package storage

import (
	"database/sql"
	"fmt"

	"github.com/jacobd39/edteam/go_sql/pkg/invoice"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceheader"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceitem"
)

//PsqlIvoice is used to work with postgres - invoice
type PsqlInvoice struct {
	db            *sql.DB
	storageHeader invoiceheader.Storage
	storageItems  invoiceitem.Storage
}

//NewPsqlInvoice returns a new pointer of PsqlInvoice
func NewPsqlInvoice(db *sql.DB, h invoiceheader.Storage, i invoiceitem.Storage) *PsqlInvoice {
	return &PsqlInvoice{
		db:            db,
		storageHeader: h,
		storageItems:  i,
	}
}

//Create implements the interface invoice.Storage
func (p *PsqlInvoice) Create(m *invoice.Model) error {
	tx, err := p.db.Begin()
	if err != nil {
		return err
	}

	if err := p.storageHeader.CreateTx(tx, m.Header); err != nil {
		tx.Rollback()
		return fmt.Errorf("Header: %w", err)
	}

	if err := p.storageItems.CreateTx(tx, m.Header.ID, m.Items); err != nil {
		tx.Rollback()
		return fmt.Errorf("Items: %w", err)
	}

	return tx.Commit()
}

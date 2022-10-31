package storage

import (
	"database/sql"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jacobd39/edteam/go_sql/pkg/invoice"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceheader"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceitem"
	"github.com/jacobd39/edteam/go_sql/pkg/product"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

//Driver of storage
type Driver string

//Drivers
const (
	MySQL    Driver = "mysql"
	Postgres Driver = "postgres"
)

//New creates the conexion with the database
func New(d Driver) {
	switch d {
	case MySQL:
		newMySQLDB()
	case Postgres:
		newPostgresDB()
	}
}

//Inicia una nueva conexión a la Base de Datos Postgress
func newPostgresDB() {
	once.Do(func() {
		//Toda esta sección se ejecutará una sola vez - patrón Singleton
		var err error
		db, err = sql.Open("postgres", "postgres://jaime:kbjnfqfsfy79@localhost:5432/godb?sslmode=disable")
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("conectado a postgres")
	})
}

//Inicia una nueva conexión a la Base de Datos MySQL
func newMySQLDB() {
	once.Do(func() {
		//Toda esta sección se ejecutará una sola vez - patrón Singleton
		var err error
		db, err = sql.Open("mysql", "bay:bayental2019@tcp(localhost:3306)/godb?parseTime=true")
		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		if err = db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("conectado a MySQL")
	})
}

//Pool return an unique instance of db
func Pool() *sql.DB {
	return db
}

func stringToNull(s string) sql.NullString {
	null := sql.NullString{String: s}
	if null.String != "" {
		null.Valid = true
	}
	return null
}

func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{Time: t}
	if !null.Time.IsZero() {
		null.Valid = true
	}
	return null
}

//DAOProduct factory of product.Storage
func DAOProduct(driver Driver) (product.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlProduct(db), nil
	case MySQL:
		return newMySQLProduct(db), nil
	default:
		return nil, fmt.Errorf("Driver not implemented")
	}
}

//DAOInvoiceHeader factory of invoiceheader.Storage
func DAOInvoiceHeader(driver Driver) (invoiceheader.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlInvoiceHeader(db), nil
	case MySQL:
		return newMySQLInvoiceHeader(db), nil
	default:
		return nil, fmt.Errorf("Driver not implemented")
	}
}

//DAOInvoiceItem factory of invoiceitem.Storage
func DAOInvoiceItem(driver Driver) (invoiceitem.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlInvoiceItem(db), nil
	case MySQL:
		return newMySQLInvoiceItem(db), nil
	default:
		return nil, fmt.Errorf("Driver not implemented")
	}
}

//DAOInvoice factory of invoice.Storage
func DAOInvoice(driver Driver, h invoiceheader.Storage, i invoiceitem.Storage) (invoice.Storage, error) {
	switch driver {
	case Postgres:
		return newMySQLInvoice(db, h, i), nil
	case MySQL:
		return newMySQLInvoice(db, h, i), nil
	default:
		return nil, fmt.Errorf("Driver not implemented")
	}
}

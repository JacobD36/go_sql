package main

import (
	"log"

	"github.com/jacobd39/edteam/go_sql/pkg/invoice"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceheader"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceitem"
	"github.com/jacobd39/edteam/go_sql/pkg/product"
	"github.com/jacobd39/edteam/go_sql/storage"
)

func main() {
	driver := storage.MySQL
	storage.New(driver)

	//Migración de tablas principales - Creación de tablas
	myProductStorage, err := storage.DAOProduct(driver)

	if err != nil {
		log.Fatalf("DAOProduct: %v", err)
	}

	serviceProduct := product.NewService(myProductStorage)
	if err := serviceProduct.Migrate(); err != nil {
		log.Fatalf("product.Migrate: %v", err)
	}

	myInvoiceHeaderStorage, err := storage.DAOInvoiceHeader(driver)

	if err != nil {
		log.Fatalf("DAOInvoiceHeader: %v", err)
	}

	serviceInvoiceHeader := invoiceheader.NewService(myInvoiceHeaderStorage)
	if err := serviceInvoiceHeader.Migrate(); err != nil {
		log.Fatalf("invoiceHeader.Migrate: %v", err)
	}

	myInvoiceItemStorage, err := storage.DAOInvoiceItem(driver)

	if err != nil {
		log.Fatalf("DAOInvoiceItem: %v", err)
	}

	serviceInvoiceItem := invoiceitem.NewService(myInvoiceItemStorage)
	if err := serviceInvoiceItem.Migrate(); err != nil {
		log.Fatalf("invoiceItem.Migrate: %v", err)
	}

	myInvoiceStorage, err := storage.DAOInvoice(driver, myInvoiceHeaderStorage, myInvoiceItemStorage)

	if err != nil {
		log.Fatalf("DAOInvoice: %v", err)
	}

	serviceInvoice := invoice.NewService(myInvoiceStorage)

	//ms, err := serviceProduct.GetAll()
	//if err != nil {
	//	log.Fatalf("product.GetAll: %v", err)
	//}
	//
	//fmt.Println(ms)

	//Inserción de nuevos datos
	//m := &product.Model{
	//	Name:         "Curso de BD con Go",
	//	Price:        80,
	//	Observations: "Este curso está disponible",
	//}
	//if err := serviceProduct.Create(m); err != nil {
	//	log.Fatalf("product.Create: %v", err)
	//}
	//
	//fmt.Printf("%+v\n", m)

	//ms, err := serviceProduct.GetAll()
	//if err != nil {
	//	log.Fatalf("product.GetAll: %v", err)
	//}
	//
	//fmt.Println(ms)

	//m, err := serviceProduct.GetByID(5)
	//
	//switch {
	//case errors.Is(err, sql.ErrNoRows):
	//	fmt.Println("No se encontró un producto con el ID solicitado")
	//case err != nil:
	//	log.Fatalf("product.GetByID: %v", err)
	//default:
	//	fmt.Println(m)
	//}
	//
	//m.Observations = "Este curso está disponible"
	//m.Price = 90
	//
	//err = serviceProduct.Update(m)
	//
	//if err != nil {
	//	log.Fatalf("product.Update: %v", err)
	//}
	//
	//fmt.Println(m)

	//err := serviceProduct.Delete(8)
	//if err != nil {
	//	log.Fatalf("product.Delete: %v", err)
	//}

	m := &invoice.Model{
		Header: &invoiceheader.Model{
			Client: "Juan Pérez Silva",
		},
		Items: invoiceitem.Models{
			&invoiceitem.Model{ProductID: 1},
			&invoiceitem.Model{ProductID: 2},
			&invoiceitem.Model{ProductID: 5},
		},
	}

	if err := serviceInvoice.Create(m); err != nil {
		log.Fatalf("Invoice.Create: %v", err)
	}
}

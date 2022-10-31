package invoice

import (
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceheader"
	"github.com/jacobd39/edteam/go_sql/pkg/invoiceitem"
)

//Model of invoice
type Model struct {
	Header *invoiceheader.Model
	Items  invoiceitem.Models
}

//Storage interface that must implement a db storage
type Storage interface {
	Create(*Model) error
}

//Service of an invoice
type Service struct {
	storage Storage
}

//NewService returns a pointer of Service
func NewService(s Storage) *Service {
	return &Service{s}
}

//Create a new invoice
func (s *Service) Create(m *Model) error {
	return s.storage.Create(m)
}

package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ProductDao interface {
	FindAllProducts() ([]*models.Product, error)
}

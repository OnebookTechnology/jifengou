package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ProductDao interface {
	FindAllProducts() ([]*models.Product, error)
	FindProductByItemStatement(itemStatement string) (p *models.Product, err error)
	AddProduct(p *models.Product) error
	FindAllCategory() ([]*models.Category, error)
	FindAllProductByBusinessIdAndStatus(businessId int, status int, pageNum, pageCount int) ([]*models.Product, error)
}

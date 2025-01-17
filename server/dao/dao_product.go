package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ProductDao interface {
	FindAllProducts() ([]*models.Product, error)
	FindProductByItemStatement(itemStatement string) (p *models.Product, err error)
	AddProduct(p *models.Product) error
	UpdateProductStatusAndCode(id, status int, code string) error
	UpdateProductById(p *models.Product) error
	FindAllCategory() ([]*models.Category, error)
	FindProductById(productId int) (*models.Product, error)
	FindAllProductByBusinessIdAndStatus(businessId int, status int, pageNum, pageCount int) ([]*models.Product, error)
	FindAllProductCountByBusinessIdAndStatus(businessId int, status int) (int, error)

	FindAllProductsOrderByScore(pageNum, pageCount int, isDesc bool, status int) ([]*models.Product, error)
	FindAllProductsOrderByOnlineTime(pageNum, pageCount, status int) ([]*models.Product, error)
	FindAllProductsOrderByExchangeTime(pageNum, pageCount, status int) ([]*models.Product, error)
}

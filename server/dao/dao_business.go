package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type BusinessDao interface {
	AddBusiness(b *models.Business) error
	QueryBusinessCount() (int, error)
	FindBusinessById(id int) (*models.Business, error)
	FindBusinessByKeyword(keyword string, pageNum, pageCount int) ([]models.Business, error)
	FindAllBusiness(pageNum, pageCount int) ([]models.Business, error)
}

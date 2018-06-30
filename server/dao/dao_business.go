package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type BusinessDao interface {
	AddBusiness(b *models.Business) error
	QueryBusinessCount() (int, error)
	FindBusinessByNo(no string) (*models.Business, error)
	FindBusinessByKeyword(keyword string, pageNum, pageCount int) ([]models.Business, error)
	FindAllBusiness(pageNum, pageCount int) ([]models.Business, error)
}

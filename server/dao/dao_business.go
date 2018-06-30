package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type BusinessDao interface {
	AddBusiness(b *models.Business) error
}

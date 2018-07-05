package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type ImageDao interface {
	AddImage(image *models.Image) (int64, error)
}

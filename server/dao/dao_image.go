package dao

import "github.com/OnebookTechnology/Engineer/server/models"

type ImageDao interface {
	AddImage(image *models.Image) error
}

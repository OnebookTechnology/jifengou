package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type TokenDao interface {
	UpdateToken(token, expireTime string) error
	FindToken() (*models.WxToken, error)
}

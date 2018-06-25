package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type UserDao interface {
	// 查找用户
	FindUser(userName, password string) (*models.User, error)
}
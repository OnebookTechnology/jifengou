package _interface

import "github.com/OnebookTechnology/jifengou/server/dao"

type ServerDB interface {
	InitialDB(confPath string, tagName string) error
	dao.UserDao
	dao.ImageDao
	dao.ProductDao
	dao.CouponDao
	dao.BusinessDao
	dao.ExRecordDao
	dao.TokenDao
}

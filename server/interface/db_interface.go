package _interface

import "github.com/OnebookTechnology/jifengou/server/dao"

type ServerDB interface {
	InitialDB(confPath string, tagName string) error
	dao.UserDao
	dao.ImageDao
	dao.ProductDao
	dao.CouponDao
	//dao.BookDao
	//dao.BookListDao
	//dao.DiscoverDao
	//dao.PressDao
	//dao.RecordDao
	//dao.ExpenseCalendar
}

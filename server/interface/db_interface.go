package _interface

import "github.com/OnebookTechnology/jifengou/server/dao"

type ServerDB interface {
	InitialDB(confPath string, tagName string) error
	dao.UserDao
	//dao.BookDao
	//dao.BookListDao
	//dao.DiscoverDao
	//dao.PressDao
	//dao.RecordDao
	//dao.ExpenseCalendar
}

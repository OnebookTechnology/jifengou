package mysql

import (
	"database/sql"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/robfig/config"
)

type MysqlService struct {
	Cfg *mysql.Config
	Db  *sql.DB
}

func (m *MysqlService) InitialDB(confPath, tagName string) error {
	c, err := config.ReadDefault(confPath)
	if err != nil {
		panic(err)
		fmt.Println(err)
	}

	m.Cfg = mysql.NewConfig()
	m.Cfg.DBName, err = c.String(tagName, "dbName")
	if err != nil {
		return err
	}
	m.Cfg.Addr, err = c.String(tagName, "addr")
	if err != nil {
		return err
	}
	m.Cfg.Net, err = c.String(tagName, "net")
	if err != nil {
		return err
	}
	m.Cfg.User, err = c.String(tagName, "user")
	if err != nil {
		return err
	}
	m.Cfg.Passwd, err = c.String(tagName, "passwd")
	if err != nil {
		return err
	}
	m.Cfg.Params = make(map[string]string)
	m.Cfg.Params["charset"], err = c.String(tagName, "charset")
	if err != nil {
		return err
	}
	db, err := sql.Open("mysql", m.Cfg.FormatDSN())
	if err != nil {
		return err
	}
	m.Db = db
	maxOpenConns, err := c.Int(tagName, "maxOpenConns")
	if err != nil {
		return err
	}
	maxIdleConns, err := c.Int(tagName, "maxIdleConns")
	if err != nil {
		return err
	}
	m.Db.SetMaxOpenConns(maxOpenConns)
	m.Db.SetMaxIdleConns(maxIdleConns)
	err = m.Db.Ping()
	if err != nil {
		return err
	}
	return nil
}

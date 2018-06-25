package mysql

import (
	"github.com/OnebookTechnology/jifengou/server/models"
)

// 查找用户
func (m *MysqlService) FindUser(userName, password string) (*models.User, error) {
	row := m.Db.QueryRow("SELECT user_id, user_name FROM user WHERE user_name=? AND password=?",
		userName, password)
	u := new(models.User)
	err := row.Scan(&u.UserId, &u.UserName)
	if err != nil {
		return nil, err
	}
	return u, nil
}

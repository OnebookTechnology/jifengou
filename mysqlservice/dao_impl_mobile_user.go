package mysql

import (
	"github.com/OnebookTechnology/jifengou/server/models"
)

// 查找手机用户
func (m *MysqlService) FindMobileUser(phoneNumber uint64) (*models.MobileUser, error) {
	row := m.Db.QueryRow("SELECT user_id, phone_number FROM mobile_user WHERE phone_number=?",
		phoneNumber)
	u := new(models.MobileUser)
	err := row.Scan(&u.UserId, &u.PhoneNumber)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func (m *MysqlService) RegisterMobileUser(newUser *models.MobileUser) error {
	tx, err := m.Db.Begin()
	_, err = tx.Exec("INSERT INTO mobile_user(phone_number,register_time) VALUES(?,?)",
		newUser.PhoneNumber, newUser.RegisterTime)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}
	err = tx.Commit()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return err
	}
	return nil
}

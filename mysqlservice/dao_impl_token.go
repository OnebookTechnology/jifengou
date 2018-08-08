package mysql

import "github.com/OnebookTechnology/jifengou/server/models"

func (m *MysqlService) UpdateToken(token, expireTime string) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("UPDATE token SET token=?, expire_time=?", token, expireTime)
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

func (m *MysqlService) FindToken() (*models.WxToken, error) {
	row := m.Db.QueryRow("SELECT token, expire_time FROM token WHERE id=1")
	c := new(models.WxToken)
	err := row.Scan(&c.Token, &c.ExpireTime)
	if err != nil {
		return nil, err
	}
	return c, nil
}

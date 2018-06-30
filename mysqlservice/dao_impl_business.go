package mysql

import (
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
)

// 添加商户
func (m *MysqlService) AddBusiness(b *models.Business) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("INSERT INTO business(business_no,business_name,business_pwd,business_info,business_register_time,business_auth) "+
		"VALUES (?,?,?,?,?,?)", b.BusinessNo, b.BusinessName, b.BusinessPwd, b.BusinessInfo, b.BusinessRegisterTime, b.BusinessAuth)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return errors.New("AddBusiness err:" + err.Error())
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

// 更新商户权限
func (m *MysqlService) Update(code string, status int, updateTime string) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("UPDATE coupon SET coupon_status=?, update_time=? WHERE coupon_code=?", status, updateTime, code)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return errors.New("UPDATE coupon err:" + err.Error())
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// 查询所有商户数量
func (m *MysqlService) QueryBusinessCount() (int, error) {
	row := m.Db.QueryRow("SELECT count(*) FROM business")
	var c int
	err := row.Scan(&c)
	if err != nil {
		return 0, err
	}
	return c, nil
}

// 查询关键字商户
func (m *MysqlService) FindBusinessById(id int) (*models.Business, error) {
	row := m.Db.QueryRow("SELECT business_id, business_no,business_name,business_register_time,business_auth FROM business "+
		"WHERE business_id=? ", id)
	b := new(models.Business)
	err := row.Scan(&b.BusinessId, &b.BusinessNo, &b.BusinessName, &b.BusinessRegisterTime, &b.BusinessAuth)
	if err != nil {
		return nil, err
	}
	return b, nil
}

// 查找所有商户
func (m *MysqlService) FindBusinessByKeyword(keyword string, pageNum, pageCount int) ([]models.Business, error) {
	rows, err := m.Db.Query("SELECT business_id,business_no,business_name,business_register_time,business_auth FROM business "+
		"WHERE business_name LIKE '%?%' LIMIT ?,?", keyword, (pageNum-1)*pageCount, pageCount)
	if err != nil {
		return nil, nil
	}
	var bs []models.Business
	for rows.Next() {
		b := new(models.Business)
		err = rows.Scan(&b.BusinessId, &b.BusinessNo, &b.BusinessName, &b.BusinessRegisterTime, &b.BusinessAuth)
		if err != nil {
			return nil, err
		}
		bs = append(bs, *b)
	}
	return bs, nil
}

// 查找所有商户
func (m *MysqlService) FindAllBusiness(pageNum, pageCount int) ([]models.Business, error) {
	rows, err := m.Db.Query("SELECT business_no,business_name,business_register_time,business_auth FROM business "+
		"LIMIT ?,?", (pageNum-1)*pageCount, pageCount)
	if err != nil {
		return nil, nil
	}
	var bs []models.Business
	for rows.Next() {
		b := new(models.Business)
		err = rows.Scan(&b.BusinessNo, &b.BusinessName, &b.BusinessRegisterTime, &b.BusinessAuth)
		if err != nil {
			return nil, err
		}
		bs = append(bs, *b)
	}
	return bs, nil
}

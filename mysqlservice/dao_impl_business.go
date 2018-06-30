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

// 查询商户
func (m *MysqlService) FindBusiness(couponCode string) (*models.Coupon, error) {
	row := m.Db.QueryRow("SELECT product_id ,coupon_code, coupon_end_time, coupon_status FROM coupon WHERE coupon_code=?",
		couponCode)
	c := new(models.Coupon)
	err := row.Scan(&c.ProductID, &c.CouponCode, &c.CouponEndTime, &c.CouponStatus)
	if err != nil {
		return nil, err
	}
	return c, nil
}

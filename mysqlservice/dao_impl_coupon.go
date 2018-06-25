package mysql

import (
	"database/sql"
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
	"time"
)

// 更新券码状态
func (m *MysqlService) UpdateCouponStatusByCouponId(couponId uint, status int) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("UPDATE coupon SET coupon_status=? WHERE coupon_id=?", status, couponId)
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

// 添加券码
func (m *MysqlService) AddCoupon(cs []*models.Coupon) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")

	// begin transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	for _, c := range cs {
		_, err := m.FindCouponByCode(c.CouponCode)
		if err != nil {
			if err != sql.ErrNoRows {
				rollBackErr := tx.Rollback()
				if rollBackErr != nil {
					return rollBackErr
				}
				return errors.New("UPDATE Coupon err:" + err.Error())
			} else {
				_, err = tx.Exec("INSERT INTO coupon(coupon_code,coupon_start_time,coupon_end_time,coupon_status) VALUES(?,?,?,?)",
					c.CouponCode, currentTime, c.CouponEndTime, models.CouponNotUsed)
				if err != nil {
					rollBackErr := tx.Rollback()
					if rollBackErr != nil {
						return rollBackErr
					}
					return errors.New("UPDATE Coupon err:" + err.Error())
				}
			}
		}
		continue
	}

	// commit transaction
	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// 查询券码
func (m *MysqlService) FindCouponByCode(couponCode string) (*models.Coupon, error) {
	row := m.Db.QueryRow("SELECT coupon_code, coupon_end_time, coupon_status FROM coupon WHERE coupon_code=?",
		couponCode)
	c := new(models.Coupon)
	err := row.Scan(&c.CouponCode, &c.CouponEndTime, &c.CouponStatus)
	if err != nil {
		return nil, err
	}
	return c, nil
}
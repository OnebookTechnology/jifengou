package mysql

import (
	"database/sql"
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
	"time"
)

// 更新券码状态
func (m *MysqlService) UpdateCouponStatusByCouponCode(code string, status int, updateTime string) error {
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
	row := m.Db.QueryRow("SELECT product_id ,coupon_code, coupon_end_time, coupon_status FROM coupon WHERE coupon_code=?",
		couponCode)
	c := new(models.Coupon)
	err := row.Scan(&c.ProductID, &c.CouponCode, &c.CouponEndTime, &c.CouponStatus)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 查询所有券码
func (m *MysqlService) FindCouponsByCount(count int, buyTime string, startTime, endTime string) ([]*models.Coupon, error) {
	rows, err := m.Db.Query("SELECT coupon_id, coupon_status, coupon_code,update_time,coupon_start_time,coupon_end_time "+
		"FROM coupon WHERE coupon_status = ? AND update_time=? AND coupon_start_time=? AND coupon_end_time = ? LIMIT ?",
		models.CouponNotUsed, buyTime, startTime, endTime, count)
	if err != nil {
		return nil, nil
	}
	var coupons []*models.Coupon
	for rows.Next() {
		c := new(models.Coupon)
		err = rows.Scan(&c.CouponId, &c.CouponStatus, &c.CouponCode, &c.UpdateTime, &c.CouponStartTime, &c.CouponId)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil
}

// 查询券码库存
func (m *MysqlService) FindCouponCountByItemStatement(itemStatement string) (count int, err error) {
	row := m.Db.QueryRow("SELECT COUNT(coupon_id) AS count 	FROM coupon c LEFT JOIN product p ON c.product_id=p.product_id 	WHERE p.product_item_statement = ?",
		itemStatement)
	err = row.Scan(&count)
	return
}

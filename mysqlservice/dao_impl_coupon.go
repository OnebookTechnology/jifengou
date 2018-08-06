package mysql

import (
	"errors"
	"github.com/OnebookTechnology/jifengou/server/models"
)

// 更新券码和商家券码状态
func (m *MysqlService) UpdateCouponStatus(code string, status int, updateTime string) error {
	c, err := m.FindCouponByCode(code)
	if err != nil {
		return err
	}
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE coupon SET coupon_status=?, update_time=? WHERE coupon_code=?", status, updateTime, code)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return errors.New("UPDATE coupon err:" + err.Error())
	}

	_, err = tx.Exec("UPDATE bcoupon SET bc_status=?, bc_update_time=? WHERE pc_id=?", status, updateTime, c.CouponId)
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return rollBackErr
		}
		return errors.New("UPDATE bcoupon err:" + err.Error())
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

// 添加券码
func (m *MysqlService) AddCoupon(c *models.Coupon) (int, error) {
	// begin transaction
	tx, err := m.Db.Begin()
	if err != nil {
		return 0, err
	}
	r, err := tx.Exec("INSERT INTO coupon(coupon_code, product_id, coupon_start_time,coupon_end_time,coupon_status,update_time) VALUES(?,?,?,?,?,?)",
		c.CouponCode, c.ProductID, c.CouponStartTime, c.CouponEndTime, models.CouponNotReleased, nowFormat())
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, errors.New("Add Coupon err:" + err.Error())
	}

	id, err := r.LastInsertId()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, errors.New("LastInsertId() err:" + err.Error())
	}
	// commit transaction
	err = tx.Commit()
	if err != nil {
		rollBackErr := tx.Rollback()
		if rollBackErr != nil {
			return 0, rollBackErr
		}
		return 0, err
	}

	return int(id), nil
}

// 查询券码
func (m *MysqlService) FindCouponByCode(couponCode string) (*models.Coupon, error) {
	row := m.Db.QueryRow("SELECT coupon_id, product_id ,coupon_code, coupon_end_time, coupon_status FROM coupon WHERE coupon_code=?",
		couponCode)
	c := new(models.Coupon)
	err := row.Scan(&c.CouponId, &c.ProductID, &c.CouponCode, &c.CouponEndTime, &c.CouponStatus)
	if err != nil {
		return nil, err
	}
	return c, nil
}

// 根据商品id和状态查询券码
func (m *MysqlService) FindCouponsByProductId(productId, status, pageNum, pageCount int) ([]*models.Coupon, error) {
	rows, err := m.Db.Query("SELECT b.bc_code, c.coupon_id, c.coupon_status, c.coupon_code, c.update_time, DATE(c.coupon_start_time), DATE(c.coupon_end_time) "+
		"FROM coupon c LEFT JOIN bcoupon b ON c.coupon_id=b.pc_id WHERE c.product_id=? AND c.coupon_status = ? "+
		"LIMIT ?,?", productId, status, (pageNum-1)*pageCount, pageCount)
	if err != nil {
		return nil, err
	}
	var coupons []*models.Coupon
	cMap := make(map[int]*models.Coupon)
	for rows.Next() {
		c := new(models.Coupon)
		var code string
		err = rows.Scan(&code, &c.CouponId, &c.CouponStatus, &c.CouponCode, &c.UpdateTime, &c.CouponStartTime, &c.CouponEndTime)
		if err != nil {
			return nil, err
		}
		if cMap[c.CouponId] == nil {
			cMap[c.CouponId] = c
			cMap[c.CouponId].BCouponCodes = append(cMap[c.CouponId].BCouponCodes, code)
		} else {
			cMap[c.CouponId].BCouponCodes = append(cMap[c.CouponId].BCouponCodes, code)
		}

	}
	for _, v := range cMap {
		coupons = append(coupons, v)
	}
	return coupons, nil
}

// 根据商品编号查询所有券码
func (m *MysqlService) FindCouponsByItemStatement(itemStatement string, count int, buyTime string, startTime, endTime string) ([]*models.Coupon, error) {
	rows, err := m.Db.Query("SELECT c.coupon_id, c.coupon_status, c.coupon_code, c.update_time, DATE(c.coupon_start_time), DATE(c.coupon_end_time) "+
		"FROM coupon c LEFT JOIN product p ON c.product_id = p.product_id "+
		"WHERE c.coupon_status = ? AND p.product_item_statement=?"+
		" LIMIT ?",
		models.CouponNotReleased, itemStatement, count)
	if err != nil {
		return nil, err
	}
	var coupons []*models.Coupon
	for rows.Next() {
		c := new(models.Coupon)
		err = rows.Scan(&c.CouponId, &c.CouponStatus, &c.CouponCode, &c.UpdateTime, &c.CouponStartTime, &c.CouponEndTime)
		if err != nil {
			return nil, err
		}
		coupons = append(coupons, c)
	}
	return coupons, nil
}

// 查询券码库存
func (m *MysqlService) FindCouponCountByItemStatement(itemStatement string) (count int, err error) {
	row := m.Db.QueryRow("SELECT COUNT(coupon_id) AS count 	FROM coupon c LEFT JOIN product p ON c.product_id=p.product_id 	"+
		" WHERE p.product_item_statement = ? AND c.coupon_status=?",
		itemStatement, models.CouponNotReleased)
	err = row.Scan(&count)
	return
}

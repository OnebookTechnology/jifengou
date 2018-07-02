package mysql

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"time"
)

// 添加商户券码
func (m *MysqlService) AddBusinessCoupon(b *models.BCoupon) error {
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec("INSERT INTO bcoupon(bc_cart_id,bc_code,b_id,product_id,bc_start,bc_end,bc_status,bc_update_time) "+
		"VALUES (?,?,?,?,?,?,?,?)", b.BCCartId, b.BCCode, b.BId, b.ProductId, b.BCStart, b.BCEnd, b.BCStatus, time.Now().Format("2006-01-02 15:04:05"))
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

// 根据状态查询商家券码
func (m *MysqlService) FindBCouponByStatus(status, productId, pageNum, pageCount int) ([]*models.BCoupon, error) {
	rows, err := m.Db.Query("SELECT bc_id,bc_cart_id,bc_code,b_id,product_id,pc_id,bc_start,bc_end,bc_status,bc_update_time FROM bcoupon "+
		" WHERE bc_status=? AND product_id=? "+
		" LIMIT ?,?", status, productId, (pageNum-1)*pageCount, pageCount)
	if err != nil {
		return nil, nil
	}
	var bs []*models.BCoupon
	for rows.Next() {
		b := new(models.BCoupon)
		err := rows.Scan(&b.BCId, &b.BCCartId, &b.BCCode, &b.BId, &b.ProductId, &b.PCId, &b.BCStart, &b.BCEnd, &b.BCStatus, &b.BCUpdateTime)
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

// 根据平台券码查询商家券码
func (m *MysqlService) FindBCouponByCouponId(couponId int) ([]*models.BCoupon, error) {
	rows, err := m.Db.Query("SELECT bc_id,bc_cart_id,bc_code,b_id,product_id,pc_id,bc_start,bc_end,bc_status,bc_update_time FROM bcoupon "+
		" WHERE pc_id=? ", couponId)
	if err != nil {
		return nil, nil
	}
	var bs []*models.BCoupon
	for rows.Next() {
		b := new(models.BCoupon)
		err := rows.Scan(&b.BCId, &b.BCCartId, &b.BCCode, &b.BId, &b.ProductId, &b.PCId, &b.BCStart, &b.BCEnd, &b.BCStatus, &b.BCUpdateTime)
		if err != nil {
			return nil, err
		}
		bs = append(bs, b)
	}
	return bs, nil
}

// 更新券码状态 (同时更新coupon 和bcoupon)
func (m *MysqlService) UpdateBCouponStatusAndCouponIdById(couponId, bcId, status int) error {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	tx, err := m.Db.Begin()
	if err != nil {
		return err
	}
	// s1. update online book's last_op_time、last_op_phone_number、online_status
	_, err = tx.Exec("UPDATE bcoupon SET pc_id=?, bc_status=?, bc_update_time=? where bc_id=?", couponId, status, currentTime, bcId)
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

// 根据平台券码得到商家券码

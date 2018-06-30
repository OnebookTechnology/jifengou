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

// 绑定商户券码（添加coupon表记录 ）

// 更新券码状态 (同时更新coupon 和bcoupon)

package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type CouponDao interface {
	AddCoupon(c *models.Coupon) (int, error)
	FindCouponByCode(couponCode string) (*models.Coupon, error)
	UpdateCouponStatus(code string, status int, updateTime string) error
	FindCouponCountByItemStatement(itemStatement string) (count int, err error)
	FindCouponsByItemStatement(itemStatement string, count int, buyTime string, startTime, endTime string) ([]*models.Coupon, error)

	AddBusinessCoupon(b *models.BCoupon) (int64, error)
	FindBCouponByStatus(status, productId, pageNum, pageCount int) ([]*models.BCoupon, error)
	UpdateBCouponStatusAndCouponIdById(couponId, bcId, status int) error
	FindCouponsByProductId(productId, status, pageNum, pageCount int) ([]*models.Coupon, error)
	FindBCouponsByCoupon(code string) ([]*models.BCoupon, error)
}

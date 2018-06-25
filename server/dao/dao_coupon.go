package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type CouponDao interface {
	FindCouponByCode(couponCode string) (*models.Coupon, error)
	UpdateCouponStatusByCouponCode(code string, status int, updateTime string) error
	FindCouponCountByItemStatement(itemStatement string) (count int, err error)
}

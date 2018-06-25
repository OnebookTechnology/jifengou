package dao

import "github.com/OnebookTechnology/jifengou/server/models"

type CouponDao interface {
	FindCouponByCode(couponCode string) (*models.Coupon, error)
}

package models

type Coupon struct {
	CouponId        int    `json:"coupon_id"`
	CouponCode      string `json:"coupon_code"`
	CouponStartTime string `json:"coupon_start_time"`
	CouponEndTime   string `json:"coupon_end_time"`
	CouponStatus    int    `json:"coupon_status"`
	UpdateTime      string `json:"updatetime"`
}

//券码状态码
const (
	CouponNotReleased = iota - 1
	CouponNotUsed
	CouponUsed
	CouponOverdue
	CouponLocked
	CouponLogOut
)

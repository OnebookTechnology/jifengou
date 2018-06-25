package models

type Coupon struct {
	CouponId        int    `json:"coupon_id"`
	CouponCode      string `json:"coupon_code"`
	CouponStartTime string `json:"coupon_start_time"`
	CouponEndTime   string `json:"coupon_end_time"`
	CouponStatus    int    `json:"coupon_status"`
	CouponUseTime   string `json:"coupon_use_time"`
}

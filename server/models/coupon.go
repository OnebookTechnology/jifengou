package models

type Coupon struct {
	CouponId        int      `json:"coupon_id"`
	ProductID       int      `json:"product_id"`
	CouponCode      string   `json:"coupon_code"`
	CouponStartTime string   `json:"coupon_start_time"`
	CouponEndTime   string   `json:"coupon_end_time"`
	CouponStatus    int      `json:"coupon_status"`
	UpdateTime      string   `json:"updatetime"`
	BCouponCodes    []string `json:"b_coupon_codes,omitempty"`
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

var CouponStatusMap = [6]string{"已绑定", "未使用", "已使用", "已过期", "已锁定", "已注销"}

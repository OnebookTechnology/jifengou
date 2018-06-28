package models

// 商家
type Business struct {
	BusinessID           int         `json:"business_id"`            // 商家ID
	BusinessName         string      `json:"business_name"`          // 商家名称
	BusinessInfo         string      `json:"business_info"`          // 商家信息
	BusinessRegisterTime string      `json:"business_register_time"` // 商家注册时间
	BProducts            []*BProduct `json:"b_products"`             // 商家的商品
	BCoupons             []*BProduct `json:"b_coupons"`              // 商家的优惠券
}

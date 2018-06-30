package models

// 商家优惠券
type BCoupon struct {
	BCId         int    `json:"bc_id"`           // 商家优惠券ID，自增
	BCCartId     string `json:"bc_cart_id"`      // 商家优惠券使用密码（无则为空）
	BCCode       string `json:"bc_code"`         // 商家优惠券券码
	BId          int    `json:"b_id"`            // 商家ID
	ProductId    int    `json:"product_id"`      // 商家商品ID
	PCId         int    `json:"pc_id,omitempty"` // 平台优惠券ID
	BCStart      string `json:"bc_start"`        // 商家优惠券可用的起始时间
	BCEnd        string `json:"bc_end"`          // 商家优惠券可用的终止时间
	BCStatus     int    `json:"bc_status"`       // 商家优惠券的状态
	BCUpdateTime string `json:"bc_update_time"`  // 商家优惠券状态修改时间
}

const (
	BCADDED   = iota // 商家优惠券已添加 0
	BCBINDED         // 商家优惠券已绑定 1
	BCUSED           // 商家优惠券已使用 2
	BCEXPIRED        // 商家优惠券已超期 3
)

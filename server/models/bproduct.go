package models

// 商品A库存示例：
// 商品A包含券码：C1,C2,C3,C4,C5,C6
// 其中C1和C2已经使用，或超过使用期限
// 则商品A的库存量为 4

// 商家商品
type BProduct struct {
	BPId     int        `json:"bp_id"`     // 商家商品ID
	BId      int        `json:"b_id"`      // 商家ID
	BPName   string     `json:"bp_name"`   // 商家商品名称
	BPPic    string     `json:"bp_pic"`    // 商家商品图片
	BPPrice  int        `json:"bp_price"`  // 商家商品价格
	BPIntro  string     `json:"bp_intro"`  // 商家商品简介
	BPAlarm  int        `json:"bp_alarm"`  // 商家商品库存警告（注：库存指的是商家提供的该商品尚未使用的券码量）
	BCoupons []*BCoupon `json:"b_coupons"` // 商家商品包含的优惠券
}

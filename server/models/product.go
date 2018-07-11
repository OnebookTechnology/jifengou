package models

//商家商品信息
type Product struct {
	ProductId            int      `json:"product_id"`
	ProductItemStatement string   `json:"product_item_statement"` //商品在积分购的编号
	ProductName          string   `json:"product_name"`
	ProductInfo          string   `json:"product_info,omitempty"`
	ProductStatus        int      `json:"product_status"` //0 未上线 1 已上线 2 已移除
	BusinessId           int      `json:"business_id"`
	ProductCategory      int      `json:"product_category"` //类型
	ProductSubtitle      string   `json:"product_subtitle,omitempty"`
	ProductPrice         float64  `json:"product_price"`
	ProductStartTime     string   `json:"product_start_time"`
	ProductEndTime       string   `json:"product_end_time"`
	ProductAlertCount    int      `json:"product_alert_count"`
	ProductOnlineTime    string   `json:"product_online_time"`
	ProductBoundCount    int      `json:"product_bound_count"`    // 平台券所绑定商家券的数量
	ProductScore         int      `json:"product_score"`          //积分
	ProductCode          string   `json:"product_code"`           // 积分购平台代码
	ProductPics          []string `json:"product_pics"`           // 商品图片
	ExchangeTime         int      `json:"exchange_time"`          //兑换次数
	ExchangeInfo         string   `json:"exchange_info"`          //兑换说明
	ProductExchangePhone string   `json:"product_exchange_phone"` //兑换电话
}

const (
	ProductReviewing = iota
	ProductSaling
	ProductRemoved
)

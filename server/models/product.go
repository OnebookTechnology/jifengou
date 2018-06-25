package models

type Product struct {
	ProductId            int     `json:"product_id"`
	ProductName          string  `json:"product_name"`
	ProductInfo          string  `json:"product_info"`
	ProductStatus        int     `json:"product_status"`
	BusinessId           int     `json:"business_id"`
	ProductCategory      int     `json:"product_category"`
	ProductPic           string  `json:"product_pic"`
	ProductSubtitle      string  `json:"product_subtitle"`
	ProductPrice         float64 `json:"product_price"`
	ProductStartTime     string  `json:"product_start_time"`
	ProductEndTime       string  `json:"product_end_time"`
	ProductAlertCount    int     `json:"product_alert_count"`
	ProductItemStatement string  `json:"product_item_statement"`
}

const (
	ProductReviewing = iota
	ProductSaling
	ProductRemoved
)

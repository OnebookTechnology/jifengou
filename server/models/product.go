package models

type Product struct {
	ProductId   int    `json:"product_id"`
	ProductName string `json:"product_name"`
	ProductInfo string `json:"product_info"`
}

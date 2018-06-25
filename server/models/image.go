package models

type Image struct {
	ImageId   int    `json:"image_id"`
	ImageName string `json:"image_name"`
	ImageType int    `json:"image_type"`
	ProductId int    `json:"product_id"`
}

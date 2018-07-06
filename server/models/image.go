package models

type Image struct {
	ImageId    int    `json:"image_id"`
	ImageUrl   string `json:"image_url"`
	ImageType  int    `json:"image_type"`
	IsMain     int    `json:"is_main"`
	ProductId  int    `json:"product_id"`
	UploadTime string `json:"upload_time"`
}

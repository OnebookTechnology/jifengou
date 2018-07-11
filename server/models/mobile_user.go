package models

type MobileUser struct {
	UserId       int    `json:"user_id"`
	PhoneNumber  uint64 `json:"phone_number"`
	RegisterTime string `json:"register_time"`
}

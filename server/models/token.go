package models

type WxToken struct {
	Token      string `json:"token"`
	ExpireTime string `json:"expire_time"`
}

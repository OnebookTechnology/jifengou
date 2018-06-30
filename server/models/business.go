package models

type Business struct {
	BusinessId           int    `json:"business_id"`
	BusinessNo           int    `json:"business_no"`
	BusinessName         string `json:"business_name"`
	BusinessPwd          string `json:"business_pwd"`
	BusinessInfo         string `json:"business_info"`
	BusinessRegisterTime string `json:"business_register_time"`
	BusinessAuth         int    `json:"business_auth"`
}

package server

import (
	"github.com/gin-gonic/gin"
)

const BusinessId = "3866229787"
const BusinessKey = "96e295d126829290dc6e906133d6a1cd"

type JFGResponse struct {
	StatusCode int           `json:"status_code"`
	Message    string        `json:"message"`
	Data       *ResponseData `json:"data"`
}

type ResponseData struct {
	Result     int        `json:"result"`
	FailReason string     `json:"fail_reason"`
	ItemCount  int        `json:"item_count,omitempty"`
	ItemList   []ItemData `json:"item_list,omitempty"`
}

type ItemData struct {
	ItemStatement string  `json:"item_statement"`
	ItemName      string  `json:"item_name"`
	ItemPrice     float64 `json:"item_price"`
}

func PlayGround() {

}

//商品查询
func QueryProduct(ctx *gin.Context) {
	itemData := ItemData{
		ItemStatement: "11190",
		ItemName:      "温莎KTV测试测试",
		ItemPrice:     10.02,
	}
	var l []ItemData
	l = append(l, itemData)
	resData := &ResponseData{
		Result:     ResultOK,
		FailReason: "",
		ItemCount:  1,
		ItemList:   l,
	}
	res := &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data:       resData,
	}

	ctx.JSON(200, res)

}

//券码信息查询
func QueryCouponInfo(ctx *gin.Context) {

}

//券码状态查询
func QueryCouponStatus(ctx *gin.Context) {

}

//券码状态更新
func UpdateCouponStatus(ctx *gin.Context) {

}

//券码库存查询
func QueryCouponCount(ctx *gin.Context) {

}

//券码使用通知
func NotifyCouponUsed(ctx *gin.Context) {

}

package server

import (
	"github.com/gin-gonic/gin"
	"strconv"
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
	Status     int        `json:"status, omitempty"`
}

type ItemData struct {
	ItemStatement string  `json:"item_statement"`
	ItemName      string  `json:"item_name"`
	ItemPrice     float64 `json:"item_price"`
}

type QueryJsonCouponStatus struct {
	Code   string `json:"code,omitempty"`
	CardId string `json:"card_id,omitempty"`
	SpId   int    `json:"sp_id,omitempty"`
}

func PlayGround() {

}

func sendFailedJsonResponse(ctx *gin.Context, resultCode int) {
	resData := &ResponseData{
		Result:     resultCode,
		FailReason: "errcode:" + strconv.Itoa(resultCode),
	}
	res := &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data:       resData,
	}
	ctx.JSON(200, res)
}

//商品查询
func QueryProduct(ctx *gin.Context) {
	crossDomain(ctx)
	products, err := server.DB.FindAllProducts()
	if err != nil {
		sendFailedJsonResponse(ctx, RequestUrlErr)
		return
	}

	var l []ItemData

	for _, p := range products {
		itemData := ItemData{
			ItemStatement: p.ProductItemStatement,
			ItemName:      p.ProductName,
			ItemPrice:     p.ProductPrice,
		}
		l = append(l, itemData)
	}

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
	crossDomain(ctx)
	var queryJson QueryJsonCouponStatus
	ctx.ShouldBindJSON(&queryJson)
	if queryJson.Code == "" {
		ctx.JSON(200, &JFGResponse{
			StatusCode: RequestFail,
			Message:    "请求失败，缺少code参数",
			Data:       nil,
		})
		return
	}

	coupon, err := server.DB.FindCouponByCode(queryJson.Code)
	if err != nil {
		logger.Error(err.Error())
		sendFailedJsonResponse(ctx, RequestUrlErr)
		return
	}

	ctx.JSON(200, &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data: &ResponseData{
			Result:     ResultOK,
			FailReason: "",
			Status:     coupon.CouponStatus,
		},
	})
	return
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

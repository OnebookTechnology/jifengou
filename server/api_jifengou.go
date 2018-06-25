package server

import (
	"encoding/hex"
	"github.com/gin-gonic/gin"
	"sort"
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

// 华易平台请求结构
type RequestJson struct {
	Code       string `json:"code,omitempty"`
	CardId     string `json:"card_id,omitempty"`
	SpId       int    `json:"sp_id,omitempty"`
	Status     int    `json:"status,omitempty"`
	UpdateTime string `json:"update_time"`
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
	var requestJson RequestJson
	// 获取华易平台的请求参数
	err := ctx.ShouldBindJSON(&requestJson)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// code必须非空
	if requestJson.Code == "" {
		ctx.JSON(200, &JFGResponse{
			StatusCode: RequestFail,
			Message:    "请求失败，缺少code参数",
			Data:       nil,
		})
		return
	}
	logger.Info("Get coupon query request.", requestJson.Code)
	// 解密
	code, err := AESDecryptHexStringToOrigin(requestJson.Code, []byte(BusinessKey))
	if err != nil {
		handleError(ctx, err)
	}
	// 在数据库中查询指定coupon
	coupon, err := server.DB.FindCouponByCode(code)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// 构造返回结果
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
	crossDomain(ctx)
	//var requestJson *RequestJson
	//err := ctx.BindJSON(&requestJson)
	//if err != nil {
	//	handleError(ctx, err)
	//	return
	//}
}

//券码库存查询
func QueryCouponCount(ctx *gin.Context) {

}

//券码使用通知
func NotifyCouponUsed(ctx *gin.Context) {

}

//积分购平台签名算法
func CalcSign(key, data, timestamp string) string {
	md5Data := doMD5FromString(data)
	var sa = sort.StringSlice{key, md5Data, timestamp}
	sort.Strings(sa)
	sa[0], sa[2] = sa[2], sa[0]
	var str string
	for _, s := range sa {
		str += s
	}
	sha1 := doSHA1([]byte(str))
	return hex.EncodeToString(sha1)
}

func handleError(ctx *gin.Context, err error) {
	logger.Error(err.Error())
	sendFailedJsonResponse(ctx, RequestUrlErr)
	return
}

package server

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
)

type JFGEnv struct {
	BusinessId         string
	BusinessKey        string
	GetCouponStatusUrl string
}

var (
	testEnv = JFGEnv{
		BusinessId:         "279916728",
		BusinessKey:        "d29a2850596496ad0a0b9821747d80b4",
		GetCouponStatusUrl: "http://api.cwidp.com/1/get_coupon_status",
	}
	onlineEnv = JFGEnv{
		BusinessId:         "3866229787",
		BusinessKey:        "96e295d126829290dc6e906133d6a1cd",
		GetCouponStatusUrl: "http://api.1710086.cn/1/get_coupon_status",
	}
)

type JFGResponse struct {
	StatusCode int           `json:"status_code"`
	Message    string        `json:"message"`
	Data       *ResponseData `json:"data"`
}

type ResponseData struct {
	Result     int    `json:"result"`
	FailReason string `json:"fail_reason"`

	//商品查询
	ItemCount int        `json:"item_count,omitempty"`
	ItemList  []ItemData `json:"item_list,omitempty"`
	Status    int        `json:"status, omitempty"`

	//券码信息查询
	Statement   string       `json:"statement,omitempty"`
	CouponCount int          `json:"coupon_count,omitempty"`
	CouponList  []CouponData `json:"coupon_list,omitempty"`

	// 库存
	StockCount int `json:"stock_count, omitempty"`
}

type CouponData struct {
	CouponId    int    `json:"coupon_id,omitempty"`
	Code        string `json:"code,omitempty"`
	CreateTime  string `json:"create_time,omitempty"`
	Status      int    `json:"status, omitempty"`
	ExpireStart string `json:"expire_start,omitempty"`
	ExpireEnd   string `json:"expire_end,omitempty"`
}

type ItemData struct {
	ItemStatement string  `json:"item_statement"`
	ItemName      string  `json:"item_name"`
	ItemPrice     float64 `json:"item_price"`
}

// 请求结构
type RequestJson struct {
	Code       string `json:"code,omitempty"`
	CardId     string `json:"card_id,omitempty"`
	SpId       string `json:"sp_id,omitempty"`
	Status     string `json:"status,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`

	//券码信息查询
	ItemStatement string `json:"item_statement,omitempty"`
	Count         string `json:"count,omitempty"`
	BuyTime       string `json:"buy_time,omitempty"`
	Statement     string `json:"statement,omitempty"`
	ExpireStart   string `json:"expire_start,omitempty"`
	ExpireEnd     string `json:"expire_end,omitempty"`
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
	crossDomain(ctx)
	var requestJson RequestJson
	// 获取华易平台的请求参数
	err := ctx.ShouldBindJSON(&requestJson)
	if err != nil {
		handleError(ctx, err)
		return
	}
	// code必须非空
	if requestJson.ItemStatement == "" || requestJson.BuyTime == "" || requestJson.ExpireEnd == "" || requestJson.ExpireStart == "" {
		ctx.JSON(200, &JFGResponse{
			StatusCode: RequestFail,
			Message:    "请求失败，缺少参数",
			Data:       nil,
		})
		return
	}

	count, err := strconv.Atoi(requestJson.Count)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if count > 20 {
		sendFailedJsonResponse(ctx, CouponInfoErr)
		return
	}

	coupons, err := server.DB.FindCouponsByItemStatement(requestJson.ItemStatement, count, requestJson.BuyTime,
		requestJson.ExpireStart, requestJson.ExpireEnd)
	if err != nil {
		handleError(ctx, err)
		return
	}

	// 检查数量
	if len(coupons) < count {
		sendFailedJsonResponse(ctx, CountNotEnoughErr)
		return
	}

	var cList []CouponData
	for _, c := range coupons {
		code, _ := AESEncryptToHexString([]byte(c.CouponCode), []byte(server.Env.BusinessKey))
		coupon := CouponData{
			CouponId:    c.CouponId,
			Code:        code,
			CreateTime:  c.UpdateTime,
			ExpireStart: requestJson.ExpireStart,
			ExpireEnd:   requestJson.ExpireEnd,
			Status:      c.CouponStatus,
		}
		cList = append(cList, coupon)
	}

	ctx.JSON(200, &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data: &ResponseData{
			Result:      ResultOK,
			FailReason:  "",
			Statement:   requestJson.Statement,
			CouponCount: count,
			CouponList:  cList,
		},
	})
	return
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
	code, err := AESDecryptHexStringToOrigin(requestJson.Code, []byte(server.Env.BusinessKey))
	if err != nil {
		handleError(ctx, err)
		return
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
	// status可能为0，所以初值设置为-99
	// 以此来判断华易是否发送了status参数
	var requestJson RequestJson
	err := ctx.BindJSON(&requestJson)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if requestJson.Code == "" {
		handleError(ctx, errors.New("华易发起状态更新请求时缺少参数code"))
		return
	}
	if requestJson.UpdateTime == "" {
		handleError(ctx, errors.New("华易发起状态更新请求时缺少参数update_time"))
		return
	}
	if requestJson.Status == "" {
		handleError(ctx, errors.New("华易发起状态更新请求时缺少参数status"))
		return
	}
	// 解密
	code, err := AESDecryptHexStringToOrigin(requestJson.Code, []byte(server.Env.BusinessKey))
	if err != nil {
		handleError(ctx, err)
		return
	}
	logger.Debug("code:", code)
	// 从数据库查询coupon
	coupon, err := server.DB.FindCouponByCode(code)
	if err != nil {
		handleError(ctx, err)
		return
	}
	status, err := strconv.Atoi(requestJson.Status)
	if err != nil {
		handleError(ctx, err)
		return
	}
	err = server.DB.UpdateCouponStatusByCouponCode(coupon.CouponCode, status, requestJson.UpdateTime)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(200, &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data: &ResponseData{
			Result:     ResultOK,
			FailReason: ""}})
	return
}

//券码库存查询
func QueryCouponCount(ctx *gin.Context) {
	crossDomain(ctx)
	var requestJson RequestJson
	err := ctx.BindJSON(&requestJson)
	if err != nil {
		handleError(ctx, err)
		return
	}
	if requestJson.ItemStatement == "" {
		handleError(ctx, errors.New("华易请求券码库存缺少参数item_statement"))
		return
	}
	p, err := server.DB.FindProductByItemStatement(requestJson.ItemStatement)
	if err != nil || p == nil {
		handleError(ctx, err)
		return
	}
	count, err := server.DB.FindCouponCountByItemStatement(requestJson.ItemStatement)
	if err != nil {
		handleError(ctx, err)
		return
	}
	ctx.JSON(200, &JFGResponse{
		StatusCode: RequestOK,
		Message:    "请求成功",
		Data: &ResponseData{
			Result:     ResultOK,
			FailReason: "",
			StockCount: count,
		},
	})
	return
}

type ResponseFromJFG struct {
	StatusCode string               `json:"status_code"`
	Message    string               `json:"message"`
	Data       *ResponseFromJFGData `json:"data"`
}

type ResponseFromJFGData struct {
	Result     string `json:"result"`
	FailReason string `json:"fail_reason"`
}

// 请求结构
type RequestJsonToJFG struct {
	SpId       int    `json:"sp_id,omitempty"`
	Code       string `json:"code,omitempty"`
	CardId     string `json:"card_id,omitempty"`
	Status     string `json:"status,omitempty"`
	UpdateTime string `json:"update_time,omitempty"`

	//券码信息查询
	ItemStatement string `json:"item_statement,omitempty"`
	Count         string `json:"count,omitempty"`
	BuyTime       string `json:"buy_time,omitempty"`
	Statement     string `json:"statement,omitempty"`
	ExpireStart   string `json:"expire_start,omitempty"`
	ExpireEnd     string `json:"expire_end,omitempty"`
}

//向积分购查询券码状态
func QueryCouponStatusFromJFG(ctx *gin.Context) {
	code := ctx.Query("code")
	cryptCode, _ := AESEncryptToHexString([]byte(code), []byte(server.Env.BusinessKey))
	//id, _ := strconv.Atoi(BusinessId)
	reqJson := &RequestJson{
		SpId: server.Env.BusinessId,
		Code: cryptCode,
	}

	reqStr, err := jsoniter.MarshalToString(reqJson)
	if err != nil {
		ctx.String(http.StatusBadRequest, "MarshalToString err: %s", err.Error())
		return
	}
	now := nowTimestampString()
	sign := CalcSign(server.Env.BusinessKey, reqStr, now)
	logger.Debug("sign:", sign)
	var url = server.Env.GetCouponStatusUrl
	url += "?sign=" + sign + "&t=" + now
	fmt.Println(url)
	resp, err := http.Post(url, "application/json;charset=utf-8", bytes.NewBuffer([]byte(reqStr)))
	defer resp.Body.Close()
	if err != nil {
		ctx.String(http.StatusInternalServerError, "coupon use notify err: %s", err.Error())
		return
	}
	rb, _ := ioutil.ReadAll(resp.Body)
	logger.Debug("coupon use response:", string(rb))
	res := &ResponseFromJFG{Data: new(ResponseFromJFGData)}
	err = jsoniter.UnmarshalFromString(string(rb), res)
	if err != nil {
		ctx.String(http.StatusOK, "UnmarshalFromString err: %s", err.Error())
		return
	}
	if res.StatusCode != "200" {
		ctx.String(http.StatusOK, "status code is: %s, message: %s", res.StatusCode, res.Message)
		return
	}
	if res.Data.Result != "1000" {
		ctx.String(http.StatusOK, "result code is: %s, message: %s", res.Data.Result, res.Data.FailReason)
		return
	}
	//TODO: 更新本地数据库列表

	ctx.String(http.StatusOK, "ok")

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

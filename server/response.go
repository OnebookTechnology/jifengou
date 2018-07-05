package server

import (
	"fmt"
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"net/http"
)

const (
	Ok  = 0
	Err = -iota
	DuplicateBusinessErr
	BCouponCountErr
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Uri     string      `json:"uri"`
	Data    interface{} `json:"data,omitempty"`
}

//注册返回数据结构
type ResData struct {
	Businesses []*models.Business `json:"businesses,omitempty"`
	Business   *models.Business   `json:"business,omitempty"`
	Categories []*models.Category `json:"categories,omitempty"`
	Products   []*models.Product  `json:"products,omitempty"`
	Product    *models.Product    `json:"product,omitempty"`
	BCoupons   []*models.BCoupon  `json:"b_coupons,omitempty"`
	Coupons    []*models.Coupon   `json:"coupons,omitempty"`
	Coupon     *models.Coupon     `json:"coupon,omitempty"`
	Image      *models.Image      `json:"image,omitempty"`
}

func Options(ctx *gin.Context) {
	crossDomain(ctx)
}

func sendFailedResponse(ctx *gin.Context, code int, v ...interface{}) {
	msg := resFormat(v...)
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Uri:     ctx.Request.RequestURI,
		Message: msg,
	})
	logger.Error("[", ctx.Request.RequestURI, "]", "ErrCode:", code, "response:", msg)

}

func sendSuccessResponse(ctx *gin.Context, data *ResData) {
	ctx.JSON(http.StatusOK, Response{
		Code:    Ok,
		Uri:     ctx.Request.RequestURI,
		Message: "ok",
		Data:    data,
	})
	s, _ := jsoniter.MarshalToString(data)
	logger.Info("[", ctx.Request.RequestURI, "]", "response:", s)

}

func resFormat(v ...interface{}) string {
	formatStr := ""
	for i := 0; i < len(v); i++ {
		formatStr += "%v "
	}
	formatStr += "\n"
	return fmt.Sprintf(formatStr, v...)
}

func crossDomain(ctx *gin.Context) {
	ctx.Header("Access-Control-Allow-Origin", "*")
	ctx.Header("Access-Control-Allow-Headers", "Authorization, Origin, No-Cache, X_Requested_With, X-Requested-With, Content-Range, X_FILENAME, If-Modified-Since, Pragma, Last-Modified, Cache-Control, Expires, Content-Type, X-E4M-With")
}

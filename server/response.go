package server

import (
	"fmt"
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	Ok  = 0
	Err = -iota
	DuplicateBusinessErr
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Uri     string      `json:"uri"`
	Data    interface{} `json:"data,omitempty"`
}

//注册返回
type ResData struct {
	Businesses []*models.Business `json:"businesses,omitempty"`
	Business   *models.Business   `json:"business,omitempty"`
	Categories []*models.Category `json:"categories,omitempty"`
	Products   []*models.Product  `json:"products,omitempty"`
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
	logger.Info("[", ctx.Request.RequestURI, "]", "response:", fmt.Sprintf("%v", *data))

}

func resFormat(v ...interface{}) string {
	formatStr := ""
	for i := 0; i < len(v); i++ {
		formatStr += "%v "
	}
	formatStr += "\n"
	return fmt.Sprintf(formatStr, v...)
}

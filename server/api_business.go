package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddBusinessReq struct {
	BNo   string `json:"b_no"`
	BName string `json:"b_name"`
	BPwd  string `json:"b_pwd"`
}

//添加商户
func AddBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req AddBusinessReq
	if err := ctx.BindJSON(&req) ; err== nil {
		b := &models.Business{
			BusinessNo:           req.BNo,
			BusinessName:         req.BName,
			BusinessPwd:          req.BPwd,
			BusinessInfo:         "",
			BusinessAuth:         1,
			BusinessRegisterTime: nowFormat(),
		}
		err := server.DB.AddBusiness(b)
		if err != nil {
			ctx.String(http.StatusServiceUnavailable, "%s", err.Error())
			return
		}
		ctx.String(http.StatusOK, "ok")
		return
	} else {
		ctx.String(http.StatusServiceUnavailable, "bind request parameter err: %s", err.Error())
		return
	}
}

func Options(ctx *gin.Context) {
	crossDomain(ctx)
}

func QueryBusinessById(ctx *gin.Context) {
	crossDomain(ctx)

}

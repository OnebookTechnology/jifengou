package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type BusinessReq struct {
	BNo   string `json:"b_no,omitempty"`
	BName string `json:"b_name,omitempty"`
	BPwd  string `json:"b_pwd,omitempty"`

	PageNum   int `json:"page_num,omitempty"`
	PageCount int `json:"page_count,omitempty"`
}

//添加商户
func AddBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
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
		logger.Info()
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

func QueryBusinessByKeyWord(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
		bs, err := server.DB.FindBusinessByKeyword(req.BName, req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindBusinessByKeyword err:", err)
			return
		}
		res := ResData{
			Businesses: bs,
		}
		sendSuccessResponse(ctx, res)
		return

	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
}

func QueryAllBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
		bs, err := server.DB.FindAllBusiness(req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindAllBusiness err:", err)
			return
		}
		res := ResData{
			Businesses: bs,
		}
		sendSuccessResponse(ctx, res)
		return

	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
}

func QueryBusinessByNo(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
		no, err := strconv.Atoi(req.BNo)
		if err != nil {
			sendFailedResponse(ctx, Err, "string convert err. req.BNo:", req.BNo)
			return
		}
		bs, err := server.DB.FindBusinessById(no)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindBusinessById err:", err)
			return
		}
		res := ResData{
			Business: *bs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
}

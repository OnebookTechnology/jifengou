package server

import (
	"database/sql"
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type BusinessReq struct {
	BId   int    `json:"b_id,omitempty" form:"b_id"`
	BNo   string `json:"b_no,omitempty" form:"b_no"`
	BName string `json:"b_name,omitempty" form:"b_name"`
	BPwd  string `json:"b_pwd,omitempty" form:"b_pwd"`

	PageNum   int `json:"page_num,omitempty" form:"page_num"`
	PageCount int `json:"page_count,omitempty" form:"page_count"`
}

//添加商户
func AddBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
		_, err := server.DB.FindBusinessByNo(req.BNo)
		if err != nil {
			if err == sql.ErrNoRows {
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
					sendFailedResponse(ctx, Err, "AddBusiness err:", err)
					return
				}
				sendSuccessResponse(ctx, nil)
				return
			} else {
				sendFailedResponse(ctx, Err, "FindBusinessByKeyword err:", err)
				return
			}

		} else {
			sendFailedResponse(ctx, DuplicateBusinessErr, "duplicate business")
			return
		}

	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

func UpdateAvail(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.BindJSON(&req); err == nil {
		err = server.DB.UpdateAvail(req.BId)
		if err != nil {
			sendFailedResponse(ctx, Err, "UpdateAvail err:", err, "businessId:", req.BId)
			return
		}
		sendSuccessResponse(ctx, nil)
		return

	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}

}

func Options(ctx *gin.Context) {
	crossDomain(ctx)
}

func QueryBusinessByKeyWord(ctx *gin.Context) {
	crossDomain(ctx)
	var req BusinessReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
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
	if err := ctx.ShouldBindQuery(&req); err == nil {
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
	if err := ctx.ShouldBindQuery(&req); err == nil {
		if err != nil {
			sendFailedResponse(ctx, Err, "string convert err. req.BNo:", req.BNo)
			return
		}
		bs, err := server.DB.FindBusinessByNo(req.BNo)
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

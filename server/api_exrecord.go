package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
)

type ExRecordReq struct {
	PhoneNumber int      `json:"phone_number" form:"phone_number"`
	BCodes      string   `json:"b_codes"`
	BCodeArray  []string `json:"b_code_array"`
	PCode       string   `json:"p_code"`
	ExTime      string   `json:"ex_time"`
	PId         int      `json:"p_id"`

	PageNum   int `json:"page_num,omitempty" form:"page_num"`
	PageCount int `json:"page_count,omitempty" form:"page_count"`
}

func QueryExRecord(ctx *gin.Context) {
	crossDomain(ctx)
	var phoneNumber int
	var err error
	if phoneNumber, err = CheckUserSessionWithPhone(ctx); err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid session. err:", err)
		return
	}

	var req ExRecordReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		rs, err := server.DB.FindExchangeRecordByPhone(phoneNumber)
		if err != nil {
			if err == sql.ErrNoRows {
				sendSuccessResponse(ctx, nil)
				return
			}
			sendFailedResponse(ctx, Err, "FindExchangeRecordByPhone. err:", err)
			return
		}
		res := &ResData{
			ExRecords: rs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

func QueryAllExRecord(ctx *gin.Context) {
	crossDomain(ctx)
	var err error
	if err = CheckSession(ctx); err != nil {
		sendFailedResponse(ctx, SessionErr, "invalid session. err:", err)
		return
	}

	var req ExRecordReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		rs, err := server.DB.FindAllExchangeRecord(req.PageNum, req.PageCount)
		if err != nil {
			if err == sql.ErrNoRows {
				sendSuccessResponse(ctx, nil)
				return
			}
			sendFailedResponse(ctx, Err, "FindAllExchangeRecord. err:", err)
			return
		}

		c, err := server.DB.FindAllExchangeCount()
		if err != nil {
			if err == sql.ErrNoRows {
				sendSuccessResponse(ctx, nil)
				return
			}
			sendFailedResponse(ctx, Err, "FindAllExchangeCount. err:", err)
			return
		}
		res := &ResData{
			ExRecords:  rs,
			TotalCount: c,
		}

		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

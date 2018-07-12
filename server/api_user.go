package server

import "github.com/gin-gonic/gin"

type UserReq struct {
	PageNum   int `json:"page_num,omitempty" form:"page_num"`
	PageCount int `json:"page_count,omitempty" form:"page_count"`
}

func ListAllUser(ctx *gin.Context) {
	crossDomain(ctx)
	var req UserReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		us, err := server.DB.ListAllUser(req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "ListAllUser err:", err)
			return
		}
		res := &ResData{
			Users: us,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}

}

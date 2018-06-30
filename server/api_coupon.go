package server

import (
	"github.com/gin-gonic/gin"
)

type CouponReq struct {
	BusinessId   int      `json:"b_id"`
	ProductId    int      `json:"p_id"`
	BCouponCodes []string `json:"b_codes"`
}

//添加商品
func AddBusinessCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.BindJSON(&req); err == nil {
		for _, bcode := range req.BCouponCodes {
			logger.Info(bcode)
		}
		//bc := &models.BCoupon{}
		//err := server.DB.AddBusinessCoupon(bc)
		//if err != nil {
		//	sendFailedResponse(ctx, Err, "AddProduct err:", err)
		//	return
		//}
		sendSuccessResponse(ctx, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

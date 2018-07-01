package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"strings"
)

type CouponReq struct {
	BusinessId   int      `json:"b_id"`
	ProductId    int      `json:"p_id"`
	ProductStart string   `json:"p_start"`
	ProductEnd   string   `json:"p_end"`
	BCouponCodes []string `json:"b_codes"`
}

//添加商品
func AddBusinessCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.BindJSON(&req); err == nil {
		for _, bcode := range req.BCouponCodes {
			//去空格
			bcode = strings.TrimSpace(bcode)
			if bcode == "" {
				continue
			}
			logger.Info("add bcoupon:", bcode)
			coupon := strings.Split(bcode, ",")
			var cartId, code string
			if len(coupon) == 1 {
				code = coupon[0]
			} else {
				cartId = coupon[0]
				code = coupon[1]
			}
			bc := &models.BCoupon{
				BCCartId:     cartId,
				BCCode:       code,
				BId:          req.BusinessId,
				ProductId:    req.ProductId,
				BCStart:      req.ProductStart,
				BCEnd:        req.ProductEnd,
				BCStatus:     models.CouponNotReleased,
				BCUpdateTime: nowTimestampString(),
			}
			err := server.DB.AddBusinessCoupon(bc)
			if err != nil {
				sendFailedResponse(ctx, Err, "AddBusinessCoupon err:", err)
				return
			}
		}
		sendSuccessResponse(ctx, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

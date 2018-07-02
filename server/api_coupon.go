package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"strings"
)

type CouponReq struct {
	BusinessId   int      `json:"b_id"`
	ProductId    int      `json:"p_id" form:"p_id"`
	BCouponCodes []string `json:"b_codes"`
	BindIds      []int    `json:"codes"`
	CouponCode   string   `form:"code"`
	Status       int      `json:"status" form:"status"`
	Exchange     string   `json:"exchange,omitempty" form:"exchange"`
	UpdateTime   string   `json:"time"`

	PageNum   int `json:"page_num,omitempty" form:"page_num"`
	PageCount int `json:"page_count,omitempty" form:"page_count"`
}

//添加商家券码
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
			p, err := server.DB.FindProductById(req.ProductId)
			if err != nil {
				sendFailedResponse(ctx, Err, "FindProductById err:", err)
				return
			}
			bc := &models.BCoupon{
				BCCartId:     cartId,
				BCCode:       code,
				BId:          req.BusinessId,
				ProductId:    req.ProductId,
				BCStart:      p.ProductStartTime,
				BCEnd:        p.ProductEndTime,
				BCStatus:     models.CouponNotBind,
				BCUpdateTime: nowTimestampString(),
			}
			err = server.DB.AddBusinessCoupon(bc)
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

// 根据商品id和状态查询券码
func QueryCouponByProductAndStatus(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		cs, err := server.DB.FindCouponsByProductId(req.ProductId, req.Status, req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindCouponsByProductId err:", err)
			return
		}
		res := &ResData{
			Coupons: cs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
}

//根据商家商品查询券码
func QueryBCouponByStatus(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		bs, err := server.DB.FindBCouponByStatus(req.Status, req.ProductId, req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindBCouponByStatus err:", err)
			return
		}
		res := &ResData{
			BCoupons: bs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
}

// 绑定券码
func BindCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.BindJSON(&req); err == nil {
		p, err := server.DB.FindProductById(req.ProductId)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindProductById err:", err)
			return
		}
		c := &models.Coupon{
			ProductID:       req.ProductId,
			CouponCode:      "JFG" + nowTimestampString() + RandText(4),
			CouponStartTime: p.ProductStartTime,
			CouponEndTime:   p.ProductEndTime,
			CouponStatus:    models.CouponNotReleased,
			UpdateTime:      nowTimestampString(),
		}
		cId, err := server.DB.AddCoupon(c)
		if err != nil {
			sendFailedResponse(ctx, Err, "AddCoupon err:", err)
			return
		}
		for _, bcId := range req.BindIds {
			err = server.DB.UpdateBCouponStatusAndCouponIdById(cId, bcId, models.CouponNotReleased)
			if err != nil {
				sendFailedResponse(ctx, Err, "UpdateBCouponStatusAndCouponIdById err:", err, "data:", cId, bcId, req.Status)
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

// 根据平台券码查询商家券码
func QueryBCouponByCoupon(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		bs, err := server.DB.FindBCouponByStatus(req.Status, req.ProductId, req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindBCouponByStatus err:", err)
			return
		}
		res := &ResData{
			BCoupons: bs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "BindJSON err:", err)
		return
	}
	return
}

// 更新券码状态
func UpdateCodeStatus(ctx *gin.Context) {
	crossDomain(ctx)
	var req CouponReq
	if err := ctx.ShouldBindJSON(&req); err == nil {

		c, err := server.DB.FindCouponByCode(req.CouponCode)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindCouponByCode err:", err, "data:", req.CouponCode)
			return
		}

		if req.Status == models.CouponUsed && c.CouponStatus != models.CouponNotUsed {
			goto RETURN
		}

		err = server.DB.UpdateCouponStatus(req.CouponCode, req.Status, req.UpdateTime)
		if err != nil {
			sendFailedResponse(ctx, Err, "UpdateCouponStatus err:", err)
			return
		}
		// 已使用，则通知积分购
		go func() {
			if req.Status == models.CouponUsed {
				err := notifyJFGUseCoupon(req.CouponCode, req.UpdateTime, "")
				if err != nil {
					logger.Error("notifyJFGUseCoupon err:", err)
				}
			}
		}()

	RETURN:
		bcs, err := server.DB.FindBCouponsByCoupon(req.CouponCode)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindCouponByCode err:", err)
			return
		}
		res := &ResData{
			BCoupons: bcs,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

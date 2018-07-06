package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"strconv"
)

const (
	ProductPics = "productPics"
)

type ProductReq struct {
	ProductId         int      `json:"p_id" form:"p_id"`
	ProductName       string   `json:"p_name"`
	ProductInfo       string   `json:"p_info,omitempty"`
	BusinessId        int      `json:"b_id" form:"b_id"`
	ProductCategory   int      `json:"p_category"` //类型
	ProductStatus     int      `json:"p_status" form:"p_status"`
	ProductSubtitle   string   `json:"p_subtitle,omitempty"`
	ProductPrice      float64  `json:"p_price"`
	ProductStartTime  string   `json:"p_start_time"`
	ProductEndTime    string   `json:"p_end_time"`
	ProductAlertCount int      `json:"p_alert_count"`
	ProductBoundCount int      `json:"p_bound_count"`
	ProductScore      int      `json:"p_score"`
	ProductCode       string   `json:"p_code"`
	ProductPics       []string `json:"p_pics"`
	ExchangeInfo      string   `json:"p_ex_info"`

	PageNum   int `json:"page_num,omitempty" form:"page_num"`
	PageCount int `json:"page_count,omitempty" form:"page_count"`
}

//添加商品
func AddProduct(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.BindJSON(&req); err == nil {
		p := &models.Product{
			ProductItemStatement: "JFG_" + strconv.Itoa(req.BusinessId) + "_" + nowTimestampString(),
			ProductName:          req.ProductName,
			ProductInfo:          req.ProductInfo,
			ProductStatus:        models.ProductReviewing,
			BusinessId:           req.BusinessId,
			ProductCategory:      req.ProductCategory,
			ProductSubtitle:      req.ProductSubtitle,
			ProductPrice:         req.ProductPrice,
			ProductStartTime:     req.ProductStartTime,
			ProductEndTime:       req.ProductEndTime,
			ProductAlertCount:    req.ProductAlertCount,
			ProductOnlineTime:    nowFormat(),
			ProductBoundCount:    req.ProductBoundCount,
			ProductScore:         req.ProductScore,
			ProductPics:          req.ProductPics,
			ExchangeInfo:         req.ExchangeInfo,
		}
		logger.Info(p.ProductPics)
		err := server.DB.AddProduct(p)
		if err != nil {
			sendFailedResponse(ctx, Err, "AddProduct err:", err)
			return
		}
		sendSuccessResponse(ctx, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

//添加商品
func UpdateProduct(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.BindJSON(&req); err == nil {
		p := &models.Product{
			ProductId:         req.ProductId,
			ProductName:       req.ProductName,
			ProductInfo:       req.ProductInfo,
			ProductCategory:   req.ProductCategory,
			ProductSubtitle:   req.ProductSubtitle,
			ProductPrice:      req.ProductPrice,
			ProductStartTime:  req.ProductStartTime,
			ProductEndTime:    req.ProductEndTime,
			ProductAlertCount: req.ProductAlertCount,
			ProductScore:      req.ProductScore,
			ProductPics:       req.ProductPics,
			ExchangeInfo:      req.ExchangeInfo,
		}
		err := server.DB.UpdateProductById(p)
		if err != nil {
			sendFailedResponse(ctx, Err, "UpdateProductById err:", err)
			return
		}
		sendSuccessResponse(ctx, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

//根据商家查找商品
func FindAllProductByBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		ps, err := server.DB.FindAllProductByBusinessIdAndStatus(req.BusinessId, req.ProductStatus, req.PageNum, req.PageCount)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindAllProductByBusinessIdAndStatus err:", err)
			return
		}
		res := &ResData{
			Products: ps,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

//根据id查找商品
func FindAllProductById(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.ShouldBindQuery(&req); err == nil {
		ps, err := server.DB.FindProductById(req.ProductId)
		if err != nil {
			sendFailedResponse(ctx, Err, "FindAllProductByBusinessIdAndStatus err:", err)
			return
		}
		res := &ResData{
			Product: ps,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

func FindAllCategory(ctx *gin.Context) {
	crossDomain(ctx)
	cs, err := server.DB.FindAllCategory()
	if err != nil {
		sendFailedResponse(ctx, Err, "FindAllCategory err:", err)
		return
	}
	res := &ResData{
		Categories: cs,
	}
	sendSuccessResponse(ctx, res)
	return
}

//添加商品
func UpdateProductStatus(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.BindJSON(&req); err == nil {
		err := server.DB.UpdateProductStatusAndCode(req.ProductId, req.ProductStatus, req.ProductCode)
		if err != nil {
			sendFailedResponse(ctx, Err, "UpdateProductStatus err:", err)
			return
		}
		sendSuccessResponse(ctx, nil)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

//根据条件查找商品
func FindAllProductByCondition(ctx *gin.Context) {
	crossDomain(ctx)
	cond := ctx.Param("condition")

	var req ProductReq
	var ps []*models.Product
	if err := ctx.ShouldBindQuery(&req); err == nil {
		switch cond {
		case "score_aesc":
			ps, err = server.DB.FindAllProductsOrderByScore(req.PageNum, req.PageCount, false)
			if err != nil {
				sendFailedResponse(ctx, Err, "FindAllProductsOrderByScore AESC err:", err)
				return
			}
		case "score_desc":
			ps, err = server.DB.FindAllProductsOrderByScore(req.PageNum, req.PageCount, true)
			if err != nil {
				sendFailedResponse(ctx, Err, "FindAllProductsOrderByScore DESC err:", err)
				return
			}
		case "exchange":
			ps, err = server.DB.FindAllProductsOrderByExchangeTime(req.PageNum, req.PageCount)
			if err != nil {
				sendFailedResponse(ctx, Err, "FindAllProductsOrderByExchangeTime err:", err)
				return
			}
		case "latest":
			ps, err = server.DB.FindAllProductsOrderByOnlineTime(req.PageNum, req.PageCount)
			if err != nil {
				sendFailedResponse(ctx, Err, "FindAllProductsOrderByOnlineTime err:", err)
				return
			}
		default:
			sendFailedResponse(ctx, Err, "invalid condition.", "data:", cond)
			return
		}
		res := &ResData{
			Products: ps,
		}
		sendSuccessResponse(ctx, res)
		return
	} else {
		sendFailedResponse(ctx, Err, "bind request parameter err:", err)
		return
	}
}

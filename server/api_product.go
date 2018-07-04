package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"
)

const (
	ProductPics = "productPics"
)

type ProductReq struct {
	ProductId         int     `json:"p_id" form:"p_id"`
	ProductName       string  `json:"p_name"`
	ProductInfo       string  `json:"p_info,omitempty"`
	BusinessId        int     `json:"b_id" form:"b_id"`
	ProductCategory   int     `json:"p_category"` //类型
	ProductStatus     int     `json:"p_status" form:"p_status"`
	ProductSubtitle   string  `json:"p_subtitle,omitempty"`
	ProductPrice      float64 `json:"p_price"`
	ProductStartTime  string  `json:"p_start_time"`
	ProductEndTime    string  `json:"p_end_time"`
	ProductAlertCount int     `json:"p_alert_count"`
	ProductBoundCount int     `json:"p_bound_count"`
	ProductScore      int     `json:"product_score"`

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
		}
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
		err := server.DB.UpdateProductStatus(req.ProductId, req.ProductStatus)
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

func savePics(form *multipart.Form, picType string, p *models.Product, ctx *gin.Context) {
	var picCount = 0
	for _, pic := range form.File[picType] {
		picCount++
		reader, _ := pic.Open()
		var fileName = ""
		logger.Info(picType)
		fileName = strconv.Itoa(p.ProductId) + "_pics" + strconv.Itoa(picCount) + ".jpg"

		data, _ := ioutil.ReadAll(reader)
		// 上传图片到本地
		ioutil.WriteFile("./images/"+fileName, data, 0777)

		// 保存图片原信息到数据库
		image := &models.Image{
			ImageName: fileName,
			ImageType: 0,
			ProductId: p.ProductId,
		}
		err := server.DB.AddImage(image)
		if err != nil {
			ctx.String(http.StatusInternalServerError, "err: %s", err.Error())
			return
		}
	}
}

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
	ProductName       string  `json:"p_name"`
	ProductInfo       string  `json:"p_info,omitempty"`
	BusinessId        int     `json:"b_id"`
	ProductCategory   int     `json:"p_category"` //类型
	ProductSubtitle   string  `json:"p_subtitle,omitempty"`
	ProductPrice      float64 `json:"p_price"`
	ProductStartTime  string  `json:"p_start_time"`
	ProductEndTime    string  `json:"p_end_time"`
	ProductAlertCount int     `json:"p_alert_count"`
}

//添加商品
func AddProduct(ctx *gin.Context) {
	crossDomain(ctx)
	var req ProductReq
	if err := ctx.BindJSON(&req); err == nil {
		p := &models.Product{
			ProductItemStatement:
		}
		err := server.DB.AddProduct(p)
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

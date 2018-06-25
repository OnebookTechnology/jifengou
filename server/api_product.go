package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"strconv"
)

const (
	ProductPics = "productPics"
)

func savePics(form *multipart.Form, picType string, p *models.Product, ctx *gin.Context) {
	var picCount = 0
	for _, pic := range form.File[picType] {
		picCount++
		reader, _ := pic.Open()
		var fileName string = ""
		logger.Info(picType)
		fileName = strconv.FormatUint(b.ISBN, 10) + "_pics" + strconv.Itoa(picCount) + ".jpg"

		// 上传图片到OSS上

		// 保存图片原信息到数据库
		err = server.DB.AddImage(image)
		if err != nil {
			sendJsonResponse(ctx, Err, "AddImage failed, error: %s", err.Error())
			return
		}
	}
}

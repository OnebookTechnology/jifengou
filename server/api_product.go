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

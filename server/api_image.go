package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"strconv"
)

func AddProductPic(ctx *gin.Context) {
	crossDomain(ctx)
	productIdStr := ctx.Query("p_id")
	productId, err := strconv.Atoi(productIdStr)
	if err != nil {
		sendFailedResponse(ctx, Err, "p_id is invalid. data:", productIdStr)
		return
	}

	form, _ := ctx.MultipartForm()
	for _, pic := range form.File["file"] {
		file, err := pic.Open()
		defer file.Close()
		if err != nil {
			sendFailedResponse(ctx, Err, "save pics:", err)
			return
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			sendFailedResponse(ctx, Err, "save pics:", err)
			return
		}
		err = ioutil.WriteFile(server.ueditorConf.ImagePath+pic.Filename, data, 0777)
		if err != nil {
			sendFailedResponse(ctx, Err, "save pics:", err)
			return
		}

		image := &models.Image{
			ImageUrl:   "http://47.93.17.108/images/" + pic.Filename,
			ImageType:  0,
			ProductId:  productId,
			UploadTime: nowFormat(),
		}

		id, err := server.DB.AddImage(image)
		if err != nil {
			sendFailedResponse(ctx, Err, "AddImage err:", err)
			return
		}
		image.ImageId = int(id)
		res := &ResData{
			Image: image,
		}
		sendSuccessResponse(ctx, res)

	}
}

package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"io/ioutil"
)

func AddProductPic(ctx *gin.Context) {
	crossDomain(ctx)

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
			ImageUrl:   "http://47.93.17.108/images/" + nowTimestampString() + "_" + pic.Filename,
			ImageType:  0,
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

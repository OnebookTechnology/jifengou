package server

import (
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

const (
	ImageUrl       = "http://47.93.17.108/ueditor?action=uploadimage"
	ImagePath      = "/images/"
	ImageFieldName = "upfile"
	ImageMaxSize   = 2048
)

var ImageAllowFiles = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}

type UEditorConfig struct {
	ImageUrl        string   `json:"imageUrl"`
	ImagePath       string   `json:"imagePath"`
	ImageFieldName  string   `json:"imageFieldName"`
	ImageMaxSize    int      `json:"imageMaxSize"`
	ImageAllowFiles []string `json:"imageAllowFiles"`
}

func UEditorHandler(ctx *gin.Context) {
	action := ctx.Query("action")
	switch action {
	case "config":
		GetUEditorConfig(ctx)
		return
	case "uploadimage":
		SavePics(ctx)
		return
	}
}

func GetUEditorConfig(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, UEditorConfig{
		ImageUrl:        ImageUrl,
		ImagePath:       ImagePath,
		ImageFieldName:  ImageFieldName,
		ImageMaxSize:    ImageMaxSize,
		ImageAllowFiles: ImageAllowFiles,
	})
}

func SavePics(ctx *gin.Context) {
	form, _ := ctx.MultipartForm()
	for _, pic := range form.File["upfile"] {
		file, err := pic.Open()
		defer file.Close()
		if err != nil {
			logger.Error("save pics:", err)
			ctx.String(http.StatusOK, "%s", err.Error())
			return
		}
		data, err := ioutil.ReadAll(file)
		if err != nil {
			logger.Error("save pics:", err)
			ctx.String(http.StatusOK, "%s", err.Error())
			return
		}
		err = ioutil.WriteFile(ImagePath+pic.Filename, data, 0777)
		if err != nil {
			logger.Error("save pics:", err)
			ctx.String(http.StatusOK, "%s", err.Error())
			return
		}
		ctx.String(http.StatusOK, "path:%s", ImagePath+pic.Filename)
	}
}

package server

import (
	"github.com/gin-gonic/gin"
	"github.com/json-iterator/go"
	"io/ioutil"
	"net/http"
)

var ImageAllowFiles = []string{".png", ".jpg", ".jpeg", ".gif", ".bmp"}

type UEditorConfig struct {
	ImageUrl        string   `json:"imageUrl"`
	ImagePath       string   `json:"imagePath"`
	ImageFieldName  string   `json:"imageFieldName"`
	ImageMaxSize    int      `json:"imageMaxSize"`
	ImageAllowFiles []string `json:"imageAllowFiles"`
	ImageActionName string   `json:"imageActionName"`
	ImageUrlPrefix  string   `json:"imageUrlPrefix"`
	StoragePath     string   `json:"storagePath"`
}

func UEditorHandler(ctx *gin.Context) {
	crossDomain(ctx)
	if ctx.Request.Method == "OPTIONS" {
		return
	}
	action := ctx.Query("action")
	callback := ctx.Query("callback")
	switch action {
	case "config":
		GetUEditorConfig(ctx, callback)
		return
	case "uploadimage":
		SavePics(ctx)
		return
	default:
		logger.Debug("actions:", action)
		return
	}
}

func GetUEditorConfig(ctx *gin.Context, callback string) {
	c := UEditorConfig{
		ImageUrl:        server.ueditorConf.ImageUrl,
		ImagePath:       server.ueditorConf.ImagePath,
		ImageFieldName:  server.ueditorConf.ImageFieldName,
		ImageMaxSize:    server.ueditorConf.ImageMaxSize,
		ImageAllowFiles: server.ueditorConf.ImageAllowFiles,
		ImageActionName: server.ueditorConf.ImageActionName,
	}
	logger.Info("ueditor config:", c)
	s, _ := jsoniter.MarshalToString(c)
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Writer.Write([]byte(callback))
	ctx.Writer.Write([]byte("("))
	ctx.Writer.Write([]byte(s))
	ctx.Writer.Write([]byte(")"))
}

type PicResponse struct {
	State    string `json:"state"`
	Url      string `json:"url"`
	Title    string `json:"title"`
	Original string `json:"original"`
	Type     string `json:"type"`
	Size     string `json:"size"`
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
		picName := nowTimestampString() + "_" + doMD5FromString(pic.Filename)
		err = ioutil.WriteFile(server.ueditorConf.ImagePath+picName, data, 0777)
		if err != nil {
			logger.Error("save pics:", err)
			ctx.String(http.StatusOK, "%s", err.Error())
			return
		}
		res := &PicResponse{
			State: "SUCCESS",
			// 解决ue加载的问题
			Url: "http://" + server.Conf.domain + "/images/" + picName + "?t=JFGTIMESTAMP",
		}
		s, _ := jsoniter.MarshalToString(res)
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Writer.Write([]byte("<div id=\"jsonData\">"))
		ctx.Writer.Write([]byte(s))
		ctx.Writer.Write([]byte("</div>"))
	}
}

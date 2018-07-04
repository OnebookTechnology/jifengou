package server

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func UploadPic(ctx *gin.Context) {
	crossDomain(ctx)
}

func savePic(picType int, filename string) (err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	name := file.Name()
	fmt.Println(name)
	//保存

	// 上传完毕之后删除缓存的图片

	// 保存图片原信息到数据库

	return
}

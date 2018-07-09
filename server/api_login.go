package server

import (
	"database/sql"
	"fmt"
	"github.com/cxt90730/xxtea-go/xxtea"
	"github.com/gin-gonic/gin"
	"net/http"
)

type LoginInfo struct {
	UserName string `json:"user_name"`
}

func Login(ctx *gin.Context) {
	crossDomain(ctx)
	user := ctx.PostForm("username")
	pwd := ctx.PostForm("password")

	u, err := server.DB.FindUser(user, pwd)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.String(http.StatusUnauthorized, "user or password is incorrect.")
			return
		}
		ctx.String(http.StatusInternalServerError, "err: %s", err.Error())
		return
	}

	session := xxtea.EncryptStdToURLString("a", "a")
	fmt.Println(session)
	//ctx.Header("SESSION")
	ctx.String(http.StatusOK, "login success: %s", u.UserName)
}

func checkSession(ctx *gin.Context) {
	//session := ctx.GetHeader("SESSION")

}

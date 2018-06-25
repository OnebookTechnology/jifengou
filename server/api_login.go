package server

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"net/http"
)

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
	ctx.String(http.StatusOK, "login success: %s", u.UserName)
}

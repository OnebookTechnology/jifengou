package server

import (
	"database/sql"
	"github.com/cxt90730/xxtea-go/xxtea"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

const XXTEA_KEY = "JiFenGou"
const SessionPrefix = "/session/"
const MaxSessionTimeout = 60 //1 minute for test

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

	session := xxtea.EncryptStdToURLString(SessionPrefix+user+"/"+nowTimestampString(), XXTEA_KEY)

	ts := strconv.FormatInt(time.Now().Add(MaxSessionTimeout*time.Second).Unix(), 10)
	err = server.Consist.Put(session, ts, MaxSessionTimeout*time.Second)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "err: %s", err)
	}
	ctx.Header("SESSION", session)
	ctx.String(http.StatusOK, "%s", u.UserName)
}

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		crossDomain(c)
		// Check if header contain SESSION
		session := c.GetHeader("SESSION")
		if len(session) == 0 {
			sendFailedResponse(c, SessionErr, "invalid session.")
			c.Abort()
			return
		}

		//TODO: Consist session
		// Check if etcd contain this sessionKey
		outTime, err := server.Consist.Get(session)
		if len(outTime) == 0 {
			sendFailedResponse(c, SessionErr, "invalid session key:", session)
			c.Abort()
			return
		}
		if err != nil {
			sendFailedResponse(c, SessionErr, "Consist.Get err:", err)
			c.Abort()
			return
		}
		// check if login object's timestamp is over time
		now := time.Now().Unix()
		sessionTime, _ := strconv.ParseInt(outTime, 10, 64)
		diff := now - sessionTime
		if diff > MaxSessionTimeout {
			// over time means current user without any operation over 30 minutes
			sendFailedResponse(c, SessionErr, "session time out")
			c.Abort()
			return
		}
		c.Next()
	}
}

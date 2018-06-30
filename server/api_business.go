package server

import (
	"github.com/OnebookTechnology/jifengou/server/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

type AddBusinessReq struct {
	BNo   int    `form:"b_no"`
	BName string `form:"b_name"`
	BPwd  string `form:"b_pwd"`
}

func AddBusiness(ctx *gin.Context) {
	crossDomain(ctx)
	var req *AddBusinessReq
	if ctx.ShouldBind(req) != nil {
		b := &models.Business{
			BusinessNo:   req.BNo,
			BusinessName: req.BName,
			BusinessPwd:  req.BPwd,
			BusinessInfo: "",
			BusinessAuth: 1,
		}
		err := server.DB.AddBusiness(b)
		if err != nil {
			ctx.String(http.StatusServiceUnavailable, "%s", err.Error())
			return
		}
		ctx.String(http.StatusOK, "ok")
		return
	} else {
		ctx.String(http.StatusServiceUnavailable, "%s", "bind request parameter err.")
		return
	}

}

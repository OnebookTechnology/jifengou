package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRouter(router *gin.Engine) {
	router.GET("/whoami", func(context *gin.Context) {
		context.String(http.StatusOK, "I am %s", server.ServerName)
	})

	myRouter := router.Group("/jifengou")
	myRouter.POST("/login", Login)

	myRouter.Any("/query_product", QueryProduct)
	myRouter.POST("/query_coupon", QueryCouponInfo)
	myRouter.POST("/update_coupon", UpdateCouponStatus)
	myRouter.POST("/query_coupon_status", QueryCouponStatus)
	myRouter.POST("/query_count", QueryCouponCount)

	myRouter.GET("/query_jfg_status", QueryCouponStatusFromJFG)
	myRouter.POST("/notify_jfg_used", NotifyCouponUsedToJFG)

	myRouter.POST("/business/add", AddBusiness)
	myRouter.POST("/business/query", QueryBusinessById)
}

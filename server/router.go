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

	businessRouter := myRouter.Group("/business")
	{
		businessRouter.POST("/add", AddBusiness)
		businessRouter.OPTIONS("/add", Options)
		businessRouter.GET("/query_keyword", QueryBusinessByKeyWord)
		businessRouter.OPTIONS("/query_keyword", Options)
		businessRouter.GET("/query_no", QueryBusinessByNo)
		businessRouter.OPTIONS("/query_no", Options)
		businessRouter.GET("/query_all", QueryAllBusiness)
		businessRouter.OPTIONS("/query_all", Options)
		businessRouter.POST("/update", UpdateAvail)
		businessRouter.OPTIONS("/update", Options)
	}

	productRouter := myRouter.Group("/product")
	{
		productRouter.GET("/category", FindAllCategory)
		productRouter.OPTIONS("/category", Options)
		productRouter.POST("/add", AddProduct)
		productRouter.OPTIONS("/add", Options)
		productRouter.GET("/query_by_bid", FindAllProductByBusiness)
		productRouter.OPTIONS("/query_by_bid", Options)
	}

	couponRouter := myRouter.Group("/coupon")
	{
		couponRouter.POST("/business/add", AddBusinessCoupon)
		couponRouter.OPTIONS("/business/add", Options)
	}

}

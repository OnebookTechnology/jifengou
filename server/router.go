package server

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadRouter(router *gin.Engine) {
	router.GET("/whoami", func(context *gin.Context) {
		context.String(http.StatusOK, "I am %s", server.ServerName)
	})

	router.Any("/ueditor", UEditorHandler)

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
	businessRouter.Use(TokenAuthMiddleware())
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
		productRouter.POST("/update_status", UpdateProductStatus)
		productRouter.OPTIONS("/update_status", Options)
		productRouter.POST("/add_pic", AddProductPic)
		productRouter.OPTIONS("/add_pic", Options)
		productRouter.GET("/query", FindProductById)
		productRouter.OPTIONS("/query", Options)
		productRouter.GET("/query_all/:condition", FindAllProductByCondition)
		productRouter.OPTIONS("/query_all/:condition", Options)
		productRouter.POST("/update", UpdateProduct)
		productRouter.OPTIONS("/update", Options)
	}

	couponRouter := myRouter.Group("/coupon")
	couponRouter.Use(TokenAuthMiddleware())
	{
		couponRouter.POST("/business/add", AddBusinessCoupon)
		couponRouter.OPTIONS("/business/add", Options)
		couponRouter.GET("/business/query", QueryBCouponByStatus)
		couponRouter.OPTIONS("/business/query", Options)
		couponRouter.POST("/bind", BindCoupon)
		couponRouter.OPTIONS("/bind", Options)
		couponRouter.GET("/query", QueryCouponByProductAndStatus)
		couponRouter.OPTIONS("/query", Options)
	}

	myRouter.POST("/coupon/update", UpdateCodeStatus)
	myRouter.OPTIONS("/coupon/update", Options)
	myRouter.GET("/coupon/query_by_item", FindProductByStatement)
	myRouter.OPTIONS("/coupon/query_by_item", Options)

	captchaGroup := myRouter.Group("/captcha")
	{
		captchaGroup.GET("/getkey", GetKey)
		captchaGroup.GET("/image/:key", ShowImage)
		captchaGroup.POST("/verify", Verify)
		captchaGroup.OPTIONS("/getkey", Options)
		captchaGroup.OPTIONS("/image/:key", Options)
		captchaGroup.OPTIONS("/verify", Options)
	}

	vcodeGroup := myRouter.Group("/vcode")
	{
		vcodeGroup.POST("/send/:key", SendVerifyCode)
		vcodeGroup.POST("/verify", VerifyVCode)
		vcodeGroup.OPTIONS("/send/:key", Options)
		vcodeGroup.OPTIONS("/verify", Options)
	}

	userRouter := myRouter.Group("/user")
	userRouter.Use(TokenAuthMiddleware())
	{
		userRouter.GET("/list", ListAllUser)
		userRouter.OPTIONS("/list", Options)
	}

}

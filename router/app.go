package router

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"weibo/docs"
	"weibo/middleware"
	"weibo/service"
)

func Router() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// 要在路由组之前全局使用「跨域中间件」, 否则OPTIONS会返回404
	r.Use(middleware.Cors())

	// swagger
	docs.SwaggerInfo.BasePath = ""
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.POST("/login", service.Login)
	r.POST("/code", service.Code)

	r.GET("/register", service.Register)

	//apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	//{
	//	...
	//}

	err := r.Run(":" + viper.GetString("server.HttpPort")) // 默认是localhost:8080
	if err != nil {
		return
	}
}

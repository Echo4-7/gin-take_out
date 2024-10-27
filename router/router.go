package router

import (
	api "Take_Out/api/v1"
	_ "Take_Out/docs"
	"Take_Out/middleware"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	gs "github.com/swaggo/gin-swagger"
	"net/http"
)

func NewRouter() *gin.Engine {

	r := gin.Default()
	r.Use(middleware.Cors())
	r.StaticFS("/static", http.Dir("./static")) // 加载静态文件
	r.GET("/swagger/*any", gs.WrapHandler(swaggerFiles.Handler))

	v1 := r.Group("api/v1")
	{
		//测试
		v1.GET("ping", func(c *gin.Context) {
			c.String(http.StatusOK, "pong")
		})

		// 用户操作
		v1.POST("user/register", api.UserRegister)
		v1.POST("user/login", api.UserLogin)
		// 轮播图
		//v1.GET("carousels", api.ListCarousel)

		auth := v1.Group("/") // 需要登陆保护  api/v1
		auth.Use(middleware.JWT())
		{
			// 更新操作
			auth.PUT("user", api.UserUpdate)
			// 上传头像
			auth.POST("avatar", api.UploadAvatar)
			// 发送邮件
			auth.POST("user/sending-email", api.SendEmail)
			// 验证邮箱
			auth.POST("/user/valid-email", api.ValidEmail)
		}
	}

	//r.NoRoute(func(c *gin.Context) {
	//	c.JSON(http.StatusOK, gin.H{
	//		"message": "404",
	//	})
	//})
	return r
}

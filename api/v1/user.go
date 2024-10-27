package v1

import (
	"Take_Out/pkg/util"
	"Take_Out/service"
	"github.com/gin-gonic/gin"
	"net/http"
)

// UserRegister 用户注册接口

func UserRegister(c *gin.Context) {
	var userRegister service.UserService
	if err := c.ShouldBind(&userRegister); err == nil {
		res := userRegister.Register(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserLogin 用户登陆接口
func UserLogin(c *gin.Context) {
	var userLogin service.UserService
	if err := c.ShouldBind(&userLogin); err == nil {
		res := userLogin.Login(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UserUpdate 用户更新接口
func UserUpdate(c *gin.Context) {
	var userUpdate service.UserService
	claims, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&userUpdate); err == nil {
		res := userUpdate.Update(c.Request.Context(), claims.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// UploadAvatar 上传头像
func UploadAvatar(c *gin.Context) {
	file, _, _ := c.Request.FormFile("file")
	var uploadAvatar service.UserService
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&uploadAvatar); err == nil {
		res := uploadAvatar.Post(c.Request.Context(), claim.ID, file)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// SendEmail 发送邮件
func SendEmail(c *gin.Context) {
	var sendEmail service.UserService
	claim, _ := util.ParseToken(c.GetHeader("Authorization"))
	if err := c.ShouldBind(&sendEmail); err == nil {
		res := sendEmail.Send(c.Request.Context(), claim.ID)
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

// ValidEmail 验证邮箱
func ValidEmail(c *gin.Context) {
	var validEmail service.UserService
	if err := c.ShouldBind(&validEmail); err == nil {
		res := validEmail.Valid(c.Request.Context(), c.GetHeader("Authorization"))
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

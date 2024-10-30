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

// SendCheckCode 发送验证码
func SendCheckCode(c *gin.Context) {
	var userSendCheckCode service.UserService
	email := c.Query("email")
	if email == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "邮箱不能为空！"})
		return
	}
	res := userSendCheckCode.SendCheckCode(c.Request.Context(), email)
	c.JSON(http.StatusOK, res)

}

// FindPwd 找回密码
func FindPwd(c *gin.Context) {
	var userFindPwd service.FindPwdService
	if err := c.ShouldBind(&userFindPwd); err == nil {
		res := userFindPwd.FindPwd(c.Request.Context())
		c.JSON(http.StatusOK, res)
	} else {
		c.JSON(http.StatusBadRequest, err)
	}
}

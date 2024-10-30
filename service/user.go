package service

import (
	"Take_Out/cache"
	"Take_Out/dao"
	"Take_Out/model"
	"Take_Out/pkg/e"
	"Take_Out/pkg/util"
	"Take_Out/serializer"
	"context"
	"fmt"
	"math/rand"
	"mime/multipart"
	"strconv"
	"time"
)

type UserService struct {
	NickName      string  `json:"nick_name" form:"nick_name"`
	Email         string  `json:"email" form:"email"`
	Password      string  `json:"password" form:"password"`
	OperationType uint    `json:"operation_type" form:"operation_type"` // 1.绑定邮箱 2.解绑邮箱 3. 改密码
	Status        string  `json:"status" form:"status"`
	Money         float64 `json:"money" form:"money"`
}

type FindPwdService struct {
	Email     string `json:"email" form:"email"`
	NewPwd    string `json:"new_pwd" form:"new_pwd"`
	CheckCode string `json:"check_code" form:"check_code"`
}

//type UserService struct {
//	UserId        int64   `json:"userId" form:"user_id"`
//	NickName      string  `json:"nickName" form:"nick_name"`
//	Password      string  `json:"password" form:"password" binding:"required"`
//	Email         string  `json:"email" form:"email"`
//	OperationType uint    `json:"operation_type" form:"operation_type"` //1 绑定邮箱 2 解绑邮箱 3 改密码
//	Status        string  `json:"status" form:"status"`
//	Money         float64 `json:"money" form:"money"`
//}

func (service *UserService) Register(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS

	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	if exist {
		code = e.ErrorExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 激活用户
	user = &model.User{
		Email:    service.Email,
		NickName: service.NickName,
		Avatar:   "avatar.jpg",
		Status:   model.Active,
	}

	// 密码加密
	if err = user.SetPassword(service.Password); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 创建用户
	err = userDao.CreateUser(user)
	if err != nil {
		code = e.ERROR
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *UserService) Login(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS

	userDao := dao.NewUserDao(ctx)

	// 判断用户是否存在
	user, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil || !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在，请先注册！",
		}
	}
	// 校验密码
	if user.CheckPassword(service.Password) == false {
		code = e.ErrorNotCompare
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "密码错误，请重新输入！",
		}
	}
	// http 无状态（认证，带上token)
	token, err := util.GenerateToken(user.ID, user.Status)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
	}
}

func (service *UserService) Update(ctx context.Context, uid uint) serializer.Response {
	var user *model.User
	// 找到用户
	userDao := dao.NewUserDao(ctx)
	code := e.SUCCESS

	user, err := userDao.GetUserByID(uid)

	// 修改昵称
	if service.NickName != "" {
		user.NickName = service.NickName
	}
	err = userDao.UpdateUserByID(user, uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// Post 头像更新
func (service *UserService) Post(ctx context.Context, uid uint, file multipart.File) serializer.Response {
	code := e.SUCCESS
	var user *model.User
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByID(uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 保存图片到本地
	userID := strconv.Itoa(int(user.ID))
	path, err := UploadAvatarToLocalStatic(file, uid, userID)
	if err != nil {
		code = e.ErrorUploadFile
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user.Avatar = path
	err = userDao.UpdateUserByID(user, uid)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
		Data:   serializer.BuildUser(user),
	}
}

// SendCheckCode 发送验证码
func (service *UserService) SendCheckCode(ctx context.Context, email string) serializer.Response {
	code := e.SUCCESS
	// 去数据库查找该用户是否已经注册
	userDao := dao.NewUserDao(ctx)
	_, exist, err := userDao.ExistOrNotExist(email)
	if err != nil || !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "用户不存在",
		}
	}

	checkCode := fmt.Sprintf("%06d", rand.Intn(1000000)) // 生成 6 位数验证码
	// 发送邮件
	err = util.SendEmail(email, checkCode, "Take_Out")
	if err != nil {
		code = e.ErrorSendEmail
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 做缓存
	if err = cache.RedisClient.Set("CHECK_CODE_MAIL:"+email, checkCode, 5*time.Minute).Err(); err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "存储验证码失败",
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

func (service *FindPwdService) FindPwd(ctx context.Context) serializer.Response {
	var user *model.User
	code := e.SUCCESS
	if service.Email == "" || service.NewPwd == "" || service.CheckCode == "" {
		code = e.InvalidParams
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 查询用户是否存在
	userDao := dao.NewUserDao(ctx)
	user, exist, err := userDao.ExistOrNotExist(service.Email)
	if err != nil || !exist {
		code = e.ErrorNotExistUser
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	// 检查验证码
	check, err := cache.RedisClient.Get("CHECK_CODE_MAIL:" + service.Email).Result()
	if check != service.CheckCode {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
			Data:   "验证码错误",
		}
	}
	// 更新密码
	if err = user.SetPassword(service.NewPwd); err != nil {
		code = e.ErrorFailEncryption
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	err = userDao.UpdateUserByID(user, user.ID)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	return serializer.Response{
		Status: code,
		Msg:    e.GetMsg(code),
	}
}

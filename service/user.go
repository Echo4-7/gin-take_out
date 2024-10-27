package service

import (
	"Take_Out/config"
	"Take_Out/dao"
	"Take_Out/model"
	"Take_Out/pkg/e"
	"Take_Out/pkg/util"
	"Take_Out/serializer"
	"context"
	"gopkg.in/mail.v2"
	"mime/multipart"
	"strconv"
	"strings"
	"time"
)

type UserService struct {
	NickName string `json:"nick_name" form:"nick_name"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	//Key      string `json:"key" form:"key"`
	OperationType uint    `json:"operation_type" form:"operation_type"` // 1.绑定邮箱 2.解绑邮箱 3. 改密码
	Status        string  `json:"status" form:"status"`
	Money         float64 `json:"money" form:"money"`
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

	//if service.Key == "" || len(service.Key) != 16 {
	//	code = e.ERROR
	//	return serializer.Response{
	//		Status: code,
	//		Msg:    e.GetMsg(code),
	//		Error:  "密钥长度不足",
	//	}
	//}

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

	//token, err := jwt.GenerateEmailToken(req.Email)
	//// 绑定邮箱
	//var address string
	//token, err := jwt.GenerateEmailToken(0, 1, service.NickName, service.Email, service.Password)
	//if err != nil {
	//	code = e.CodeInvalidToken
	//	return serializer.Response{
	//		Status: code,
	//		Msg:    e.GetMsg(code),
	//	}
	//}
	//// 发送邮件
	//address = config.Conf.ValidEmail + token
	//emailStr := "您正在绑定邮箱Email"
	//emailText := strings.Replace(emailStr, "Email", address, -1)
	//m := mail.NewMessage()
	//m.SetHeader("From", config.Conf.SmtpEmail)
	//m.SetHeader("To", service.Email)
	//m.SetHeader("Subject", "Take_out")
	//m.SetBody("text/html", emailText)
	////创建一个新的 SMTP 发送器实例
	//d := mail.NewDialer(config.Conf.Email.SmtpHost, 465, config.Conf.Email.SmtpEmail, config.Conf.Email.SmtpPass)
	//d.StartTLSPolicy = mail.MandatoryStartTLS
	//if err = d.DialAndSend(m); err != nil {
	//	code = e.CodeErrSendEmail
	//	return serializer.Response{
	//		Status: code,
	//		Msg:    e.GetMsg(code),
	//	}
	//}
	//
	//return serializer.Response{
	//	Status: code,
	//	Msg:    "请验证邮箱",
	//	Data:   serializer.BuildUser(user),
	//}

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

// Send 发送邮箱
func (service *UserService) Send(ctx context.Context, uid uint) serializer.Response {
	var notice *model.Notice
	code := e.SUCCESS
	token, err := util.GenerateEmailToken(uid, service.OperationType, service.Email, service.NickName, service.Password)
	if err != nil {
		code = e.ErrorAuthToken
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	noticeDao := dao.NewNoticeDao(ctx)
	notice, err = noticeDao.GetNoticeById(service.OperationType)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	address := config.Config.Email.ValidEmail + token // 发送方
	mailStr := notice.Text
	mailText := strings.Replace(mailStr, "Email", "<a href=\""+address+"\">"+address+"</a>", -1)
	m := mail.NewMessage()
	m.SetHeader("From", config.Config.Email.SmtpEmail)
	m.SetHeader("To", service.Email)
	m.SetHeader("Subject", "Take_Out")
	m.SetBody("text/html", mailText)
	d := mail.NewDialer(config.Config.Email.SmtpHost, 465, config.Config.Email.SmtpEmail, config.Config.Email.SmtpPass)
	d.StartTLSPolicy = mail.MandatoryStartTLS
	if err = d.DialAndSend(m); err != nil {
		code = e.ErrorSendEmail
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

// Valid 校验邮箱
func (service *UserService) Valid(ctx context.Context, token string) serializer.Response {
	var userId uint
	var email string
	var password string
	var operationType uint
	var nickName string
	code := e.SUCCESS
	if token == "" {
		code = e.InvalidParams
	} else {
		claims, err := util.ParseEmailToken(token)
		if err != nil {
			code = e.ErrorAuthToken
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = e.ErrorAuthCheckTokenTimeout
		} else {
			userId = claims.UserID
			email = claims.Email
			password = claims.Password
			operationType = claims.OperationType
			nickName = claims.Nickname
		}
	}
	if code != e.SUCCESS {
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}

	// 获取用户信息
	userDao := dao.NewUserDao(ctx)
	user, err := userDao.GetUserByID(userId)
	if err != nil {
		code = e.ERROR
		return serializer.Response{
			Status: code,
			Msg:    e.GetMsg(code),
		}
	}
	user = &model.User{
		NickName: nickName,
		Status:   model.Active,
		//Avatar:
	}

	if operationType == 1 {
		user.Email = email
	} else if operationType == 2 {
		user.Email = ""
	} else if operationType == 3 {
		err = user.SetPassword(password)
		if err != nil {
			code = e.ERROR
			return serializer.Response{
				Status: code,
				Msg:    e.GetMsg(code),
			}
		}
	}
	err = userDao.UpdateUserByID(user, userId)
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

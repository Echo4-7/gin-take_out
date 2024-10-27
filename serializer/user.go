package serializer

import (
	"Take_Out/config"
	"Take_Out/model"
)

// User 返回给前端的结构体
type User struct {
	ID       uint    `json:"id"`
	Email    string  `json:"email"`
	NickName string  `json:"nick_name"`
	Status   string  `json:"status"`
	Avatar   string  `json:"avatar"`
	CreateAt int64   `json:"create_at"`
	Money    float64 `json:"money"`
}

func BuildUser(user *model.User) *User {
	return &User{
		ID:       user.ID,
		Email:    user.Email,
		NickName: user.NickName,
		Status:   user.Status,
		Avatar:   config.Config.Path.PhotoHost + config.Config.System.HttpPort + config.Config.Path.AvatarPath + user.Avatar,
		CreateAt: user.CreatedAt.Unix(),
		Money:    user.Money,
	}
}

package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	UserName       string
	Email          string
	NickName       string
	PasswordDigest string
	Avatar         string
	Status         string
	Money          float64
}

const (
	PasswordCost        = 12 // 密码加密难度
	Pending      string = "Pending"
	Active       string = "Active"
)

// 密码加密

func (user *User) SetPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), PasswordCost)
	if err != nil {
		return err
	}
	user.PasswordDigest = string(bytes)
	return nil
}

// 校验密码

func (user *User) CheckPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordDigest), []byte(password))
	return err == nil
}

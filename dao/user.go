package dao

import (
	"Take_Out/model"
	"context"
	"gorm.io/gorm"
)

type UserDao struct {
	*gorm.DB
}

// 根据上下文 ctx 创建一个新的 UserDao 实例

func NewUserDao(ctx context.Context) *UserDao {
	return &UserDao{NewDBClient(ctx)}
}

func NewUserDaoByDB(db *gorm.DB) *UserDao {
	return &UserDao{db}
}

// ExistOrNotExist 检查用户是否存在
func (dao *UserDao) ExistOrNotExist(email string) (user *model.User, exist bool, err error) {
	var count int64
	err = dao.DB.Model(&model.User{}).Where("email = ?", email).Find(&user).Count(&count).Error
	if count == 0 {
		return nil, false, err
	}
	return user, true, nil
}

// CreateUser 创建新用户
func (dao *UserDao) CreateUser(user *model.User) error {
	return dao.DB.Model(&model.User{}).Create(&user).Error
}

// GetUserByID 根据ID获取用户
func (dao *UserDao) GetUserByID(uid uint) (user *model.User, err error) {
	err = dao.DB.Model(&model.User{}).Where("id = ?", uid).First(&user).Error
	return
}

// UpdateUserByID 根据ID更新用户信息
func (dao *UserDao) UpdateUserByID(user *model.User, uid uint) error {
	return dao.DB.Model(&model.User{}).Where("id = ?", uid).Updates(&user).Error
}

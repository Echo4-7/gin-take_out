package service

import (
	"Take_Out/config"
	"io"
	"mime/multipart"
	"os"
	"strconv"
)

// UploadAvatarToLocalStatic 更新头像到本地
func UploadAvatarToLocalStatic(file multipart.File, uid uint, userID string) (filePath string, err error) {
	bId := strconv.Itoa(int(uid)) // 路径拼接
	basePath := "." + config.Config.Path.AvatarPath + "user" + bId + "/"
	if !DirExistOrNot(basePath) {
		CreateDir(basePath)
	}
	avatarPath := basePath + userID + ".jpg"
	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}
	err = os.WriteFile(avatarPath, content, 06666)
	if err != nil {
		return
	}
	return "user" + bId + "/" + userID + ".jpg", err
}

// DirExistOrNot 判断路径是否存在
func DirExistOrNot(fileAddr string) bool {
	s, err := os.Stat(fileAddr)
	if err != nil {
		return false
	}
	return s.IsDir()
}

// CreateDir 创建文件夹
func CreateDir(dirName string) bool {
	err := os.MkdirAll(dirName, 0755)
	if err != nil {
		return false
	}
	return true
}

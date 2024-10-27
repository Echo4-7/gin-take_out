package dao

import (
	"Take_Out/model"
	"fmt"
)

func migration() {
	err := _db.Set("gorm:table_options", "charset=utf8mb4").
		AutoMigrate(&model.User{}, &model.Notice{})

	if err != nil {
		fmt.Println("err:", err)
	}
	return
}

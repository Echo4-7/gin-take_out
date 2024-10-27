package dao

import (
	"context"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"
	"time"
)

var _db *gorm.DB

func DBEngine(connRead, connWrite string) {
	var ormLogger logger.Interface
	if gin.Mode() == "debug" {
		ormLogger = logger.Default.LogMode(logger.Info)
	} else {
		ormLogger = logger.Default
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       connRead,
		DefaultStringSize:         256,  //string默认字段长度
		DisableDatetimePrecision:  true, //禁用datetime精度
		DontSupportRenameIndex:    true, //重命名索引，就要把索引先删除再重建
		DontSupportRenameColumn:   true, //用change重命名列
		SkipInitializeWithVersion: false,
	}), &gorm.Config{
		Logger: ormLogger,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		return
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)  //设置连接池
	sqlDB.SetMaxOpenConns(100) //打开的连接数
	sqlDB.SetConnMaxLifetime(time.Second * 30)
	_db = db
	//主从配置
	_ = _db.Use(dbresolver.Register(dbresolver.Config{
		Sources:  []gorm.Dialector{mysql.Open(connWrite)},
		Replicas: []gorm.Dialector{mysql.Open(connRead), mysql.Open(connRead)},
		Policy:   dbresolver.RandomPolicy{},
	}))
	migration()
}

func NewDBClient(ctx context.Context) *gorm.DB {
	db := _db
	return db.WithContext(ctx)
}

package config

import (
	"Take_Out/dao"
	"fmt"
	"github.com/spf13/viper"
	"os"
	"strings"
)

var Config *Conf

type Conf struct {
	System *System `yaml:"system"`
	Mysql  *MySql  `yaml:"mysql"`
	Redis  *Redis  `yaml:"redis"`
	Email  *Email  `yaml:"email"`
	Path   *Path   `yaml:"path"`
}

type System struct {
	AppEnv      string `yaml:"env"`
	Domain      string `yaml:"domain"`
	Version     string `yaml:"version"`
	HttpPort    string `yaml:"HttpPort"`
	Host        string `yaml:"Host"`
	UploadModel string `yaml:"UploadModel"`
}

type MySql struct {
	Dialect  string `yaml:"dialect"`
	DbHost   string `yaml:"dbHost"`
	DbPort   string `yaml:"dbPort"`
	DbName   string `yaml:"dbName"`
	UserName string `yaml:"userName"`
	Password string `yaml:"password"`
	Charset  string `yaml:"charset"`
}

type Redis struct {
	RedisHost string `yaml:"redisHost"`
	RedisPort string `yaml:"redisPort"`
	//RedisUsername string `yaml:"redisUsername"`
	RedisPassword string `yaml:"redisPassword"`
	RedisDbName   int    `yaml:"redisDbName"`
	RedisNetwork  string `yaml:"redisNetwork"`
}

type Email struct {
	ValidEmail string `yaml:"ValidEmail"`
	SmtpHost   string `yaml:"SmtpHost"`
	SmtpEmail  string `yaml:"SmtpEmail"`
	SmtpPass   string `yaml:"SmtpPass"`
}

// EncryptSecret 加密的东西
type EncryptSecret struct {
	JwtSecret   string `yaml:"jwtSecret"`
	EmailSecret string `yaml:"emailSecret"`
	PhoneSecret string `yaml:"phoneSecret"`
	//MoneySecret string `yaml:"moneySecret"`
}

type Path struct {
	PhotoHost  string `yaml:"photoHost"`
	AvatarPath string `yaml:"avatarPath"`
}

func Init() (err error) {
	workDir, _ := os.Getwd()
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(workDir + "/config")
	viper.AddConfigPath(workDir)
	//viper.SetConfigFile("./config/config.yaml") // 指定配置文件路径
	// 读取配置信息
	err = viper.ReadInConfig()
	if err != nil {
		fmt.Printf("viper.ReadInConfig failed, err: %v\n", err)
		return
	}
	// 把读取到的配置信息反序列化到 Conf 变量中
	if err := viper.Unmarshal(&Config); err != nil {
		fmt.Printf("viper.Unmarshal failed, err: %v\n", err)
	}
	pathRead := strings.Join([]string{Config.Mysql.UserName, ":", Config.Mysql.Password, "@tcp(", Config.Mysql.DbHost, ":", Config.Mysql.DbPort, ")/", Config.Mysql.DbName, "?charset=utf8mb4&parseTime=true&loc=Local"}, "")
	pathWrite := strings.Join([]string{Config.Mysql.UserName, ":", Config.Mysql.Password, "@tcp(", Config.Mysql.DbHost, ":", Config.Mysql.DbPort, ")/", Config.Mysql.DbName, "?charset=utf8mb4&parseTime=true&loc=Local"}, "")
	dao.DBEngine(pathRead, pathWrite)
	return
}

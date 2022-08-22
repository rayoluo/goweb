package utils

import (
	"fmt"
	"gopkg.in/ini.v1"
)

var (
	AppMode  string
	HttpPort string
	JwtKey   string

	Db         string
	DbHost     string
	DbPort     string
	DbUser     string
	DbPassWord string
	DbName     string

	AccessKey  string
	SecretKey  string
	Bucket     string
	QiniuSever string

	MailAddress string
	MailPass    string

	RedisAddr string
	RedisPort string
	RedisDB   int

	IdleConnSize   int
	SessionNetwork string
	StoreAddr      string
	StorePort      string
	StoreSecret    string

	Avatar string
)

func init() {
	file, err := ini.Load("config/config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
	}
	LoadServer(file)
	LoadData(file)
	LoadQiniu(file)
	LoadMail(file)
	LoadRedis(file)
	LoadUserInfo(file)
	LoadSessionStore(file)
}

func LoadServer(file *ini.File) {
	AppMode = file.Section("server").Key("AppMode").MustString("debug")
	HttpPort = file.Section("server").Key("HttpPort").MustString(":3000")
	JwtKey = file.Section("server").Key("JwtKey").MustString("89js82js72")
}

func LoadData(file *ini.File) {
	Db = file.Section("database").Key("Db").MustString("debug")
	DbHost = file.Section("database").Key("DbHost").MustString("localhost")
	DbPort = file.Section("database").Key("DbPort").MustString("3306")
	DbUser = file.Section("database").Key("DbUser").MustString("ginblog")
	DbPassWord = file.Section("database").Key("DbPassWord").MustString("admin123")
	DbName = file.Section("database").Key("DbName").MustString("ginblog")
}

func LoadQiniu(file *ini.File) {
	AccessKey = file.Section("qiniu").Key("AccessKey").String()
	SecretKey = file.Section("qiniu").Key("SecretKey").String()
	Bucket = file.Section("qiniu").Key("Bucket").String()
	QiniuSever = file.Section("qiniu").Key("QiniuSever").String()
}

func LoadMail(file *ini.File) {
	MailAddress = file.Section("mail").Key("MailAddress").String()
	MailPass = file.Section("mail").Key("MailPass").String()
}

func LoadRedis(file *ini.File) {
	RedisAddr = file.Section("redis").Key("RedisAddr").String()
	RedisPort = file.Section("redis").Key("RedisPort").String()
	RedisDB, _ = file.Section("redis").Key("RedisDB").Int()
}

func LoadUserInfo(file *ini.File) {
	Avatar = file.Section("userinfo").Key("Avatar").String()
}

func LoadSessionStore(file *ini.File) {
	IdleConnSize, _ = file.Section("sessionStore").Key("IdleConnSize").Int()
	SessionNetwork = file.Section("sessionStore").Key("SessionNetwork").String()
	StoreAddr = file.Section("sessionStore").Key("StoreAddr").String()
	StorePort = file.Section("sessionStore").Key("StorePort").String()
	StoreSecret = file.Section("sessionStore").Key("StoreSecret").String()
}

package initialize

import "github.com/go-ini/ini"

var HttpPort string
var LogLevel int
var QueueNum int
var Host string
var Port string
var PassWord string

func NewConfig(file *ini.File) {
	HttpPort = file.Section("server").Key("HttpPort").MustString("8080")
	LogLevel = file.Section("server").Key("LogLevel").MustInt(1)
	QueueNum = file.Section("server").Key("QueueNum").MustInt(10)
	Host = file.Section("redis").Key("Host").MustString("127.0.0.1")
	Port = file.Section("redis").Key("Port").MustString("Port")
	PassWord = file.Section("redis").Key("PassWord").MustString("")
}

//[server]
//HttpPort = :8080
//LogLevel = 1
//QueueNum = 10
//
//
//
//[redis]
//Host = 127.0.0.1
//Port = 3306
//PassWord =

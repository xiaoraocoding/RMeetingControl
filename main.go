package main

import (
	"RMeetingControl/api"
	"RMeetingControl/initialize"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-ini/ini"
)

func main() {
	file, err := ini.Load("./config.ini")
	if err != nil {
		fmt.Println("配置文件读取错误，请检查文件路径:", err)
		return
	}
	initialize.NewConfig(file)
	initialize.NewLog()
	initialize.Log.SetLevel(initialize.LogLevel)
	initialize.NewMeetingGroup()
	initialize.NewChanControl()
	initialize.NewQueue(initialize.QueueNum)
	initialize.Queue.StartWorkPool()
	// 下一次迭代使用yaml文件更改配置，目前只需要再此处即可进行修改
	initialize.ConnectRedis(initialize.Host+":"+initialize.Port, initialize.PassWord, 0)
	r := gin.Default()
	r.POST("/CreateMeeting", api.CreateMeeting)
	r.PUT("/LeaveMeeting", api.LeaveMeeting)
	r.GET("/AddMeeting", api.AddMeeting)
	r.GET("/Export", api.Export)
	// 下面这些接口最好是直接使用客户端的长链接，将修改发送给服务端，这里为了理解业务，方便后端测试，
	r.Run(":" + initialize.HttpPort) // listen and serve on 0.0.0.0:8080
}

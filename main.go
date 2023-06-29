package main

import (
	"RMeetingControl/api"
	"RMeetingControl/initialize"
	"github.com/gin-gonic/gin"
)

func main() {
	initialize.NewLog()
	initialize.Log.SetLevel(initialize.Info)
	initialize.NewMeetingGroup()
	initialize.NewChanControl()
	initialize.NewQueue(10)
	initialize.Queue.StartWorkPool()
	// 下一次迭代使用yaml文件更改配置，目前只需要再此处即可进行修改
	initialize.ConnectRedis("127.0.0.1:6379", "", 0)
	r := gin.Default()
	r.POST("/CreateMeeting", api.CreateMeeting)
	r.PUT("/LeaveMeeting", api.LeaveMeeting)
	r.GET("/AddMeeting", api.AddMeeting)
	r.GET("/Export", api.Export)
	// 下面这些接口最好是直接使用客户端的长链接，将修改发送给服务端，这里为了理解业务，方便后端测试，
	r.Run() // listen and serve on 0.0.0.0:8080
}

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
	// 下一次迭代使用yaml文件更改配置，目前只需要再此处即可进行修改
	initialize.ConnectRedis("127.0.0.1:6379", "", 0)
	r := gin.Default()
	r.POST("/CreateMeeting", api.CreateMeeting)
	r.PUT("/LeaveMeeting", api.LeaveMeeting)
	r.GET("/AddMeeting", api.AddMeeting)
	r.GET("/Export", api.Export)
	r.PUT("/ChangeVideoStatus", api.ChangeVideoStatus)
	r.PUT("/ChangeMuteStatus", api.ChangeMuteStatus)
	r.PUT("/ChangeNameStatus", api.ChangeNameStatus)
	r.Run() // listen and serve on 0.0.0.0:8080
}

package initialize

import (
	"RMeetingControl/model"
	"github.com/gorilla/websocket"
	"time"
)

var (
	pingPeriod = 10 * time.Second
)

//开始进行心跳，这里会对每一个conn开启一个goroutine
func StartHeartbeat(meetingUid string, userUid string) {
	heartbeatTicker := time.NewTicker(5 * time.Second) // 每隔 5 秒发送一次心跳消息
	defer heartbeatTicker.Stop()

	for {
		select {
		case <-heartbeatTicker.C:
			// 发送心跳消息
			if err := model.Group[meetingUid].AllConn[userUid].Conn.WriteMessage(websocket.TextMessage, []byte("heartbeat")); err != nil {
				Log.Error("Error Failed to send heartbeat:", err)
				model.Chan.AllChan[userUid] <- 1
				return
			}
		}
	}
}

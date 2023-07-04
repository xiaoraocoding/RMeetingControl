package api

// 会议相关router
import (
	"RMeetingControl/initialize"
	"RMeetingControl/model"
	"RMeetingControl/util"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/websocket"
	"net/http"
)

func LeaveMeeting(ctx *gin.Context) {
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	model.Chan.AllChan[user.UserUid] <- 1 // 给通道传递消息，结束掉conn链接
	model.Chan.HeartChan[user.UserUid] <- 1
	ctx.JSON(http.StatusOK, gin.H{
		"msg":      "leaveMeeting success!",
		"userUid":  user.UserUid,
		"username": user.Username,
	})
}

func CreateMeeting(ctx *gin.Context) {
	client := model.Client{}
	if err := ctx.ShouldBindJSON(&client); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	meetingUid := util.CreateUUid()
	userUid := util.CreateUUid()
	// 使用redis事务
	tx := initialize.Redis.Client.TxPipeline()
	// redis存放会议id，会议名称
	tx.Set(initialize.Redis.Context, meetingUid, client.MeetingName, redis.KeepTTL)
	// redis存放会议管理者id
	tx.Set(initialize.Redis.Context, meetingUid+"status", userUid, redis.KeepTTL)
	// redis存放当前会议所有用户信息
	tx.HSet(initialize.Redis.Context, meetingUid+"all", userUid, client.Username)
	_, err := tx.Exec(initialize.Redis.Context)
	if err != nil {
		initialize.Log.Error("Error: redis exec failed")
		return
	}
	initialize.Log.Info("Info: redis exec success", client.Username)
	mc := initialize.NewMeetingClient(client.IsVideo, client.MeetingName, client.IsMute)
	model.Group[meetingUid] = mc
	ctx.JSON(http.StatusOK, gin.H{
		"meetingUid":  meetingUid,
		"meetingName": client.MeetingName,
		"status":      "success",
		"userUid":     userUid,
	})
}

func AddMeeting(ctx *gin.Context) {
	username := ctx.Query("username")
	meetingUid := ctx.Query("meetingUid")
	userUuid := ctx.Query("userUid")
	if userUuid == "" {
		userUuid = util.CreateUUid()
	}
	// 将 HTTP 请求升级为 WebSocket 连接
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	defer conn.Close()
	if err != nil {
		initialize.Log.Error("Error Upgrade error:", err)
		return
	}
	clientc := model.ClientConnection{
		conn,
		username,
		userUuid,
	}
	if ok := initialize.Redis.Hset(meetingUid+"all", userUuid, username); !ok {
		initialize.Log.Error("Error: 加入用户失败")
		return
	}
	initialize.Log.Info("Info: 加入用户成功")
	msg := model.Message{
		Notice: "有成员加入会议室",
		Name:   username,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		initialize.Log.Error("Error json error:", err)
		return
	}
	model.Group[meetingUid].AllConn[userUuid] = clientc
	model.Chan.AllChan[userUuid] = make(chan int, 10)
	model.Chan.HeartChan[userUuid] = make(chan int, 10)
	content := initialize.Content{
		Message:    jsonData,
		MeetingUid: meetingUid,
	}
	initialize.Queue.SendMsgToQueue(content)
	control := model.Control{
		IsVideo:      model.Group[meetingUid].IsVideo,
		IsMute:       model.Group[meetingUid].IsMute,
		IsChangeName: model.Group[meetingUid].IsChangeName,
		UserUid:      userUuid,
	}
	controlData, err := json.Marshal(control)
	if err != nil {
		initialize.Log.Error("Error json error:", err)
		return
	}
	model.Group[meetingUid].Mutex.Lock()
	// 广播给当前加入的新人，会议室的状态
	err = conn.WriteMessage(websocket.TextMessage, controlData)
	if err != nil {
		initialize.Log.Error("Error send new people failed:", err)
		return
	}
	model.Group[meetingUid].Mutex.Unlock()
	messageCh := make(chan []byte)
	// 每进行了一个conn链接，就为此conn开启一个心跳
	go initialize.StartHeartbeat(meetingUid, userUuid)
	go func() {
		for {
			messageType, message, err := conn.ReadMessage()
			if err != nil {
				initialize.Log.Error("Error conn ReadMessage failed:", err)
				return
			}
			if messageType == websocket.TextMessage || messageType == websocket.BinaryMessage {
				messageCh <- message
			}
		}
	}()
	for {
		select {
		case <-model.Chan.AllChan[userUuid]: //说明那边断开了业务
			if ok := initialize.Redis.HDel(meetingUid+"all", userUuid); !ok {
				initialize.Log.Error("Error: 删除用户失败")
				return
			}
			msg := model.Message{
				Notice: "有成员离开会议室",
				Name:   username,
			}
			jsonData, err := json.Marshal(msg)
			if err != nil {
				initialize.Log.Error("Error json error:", err)
				return
			}
			content := initialize.Content{
				Message:    jsonData,
				MeetingUid: meetingUid,
			}
			initialize.Queue.SendMsgToQueue(content)
			model.Group[meetingUid].Mutex.Lock()
			for _, allConn := range model.Group[meetingUid].AllConn {
				if allConn.Uid == userUuid {
					delete(model.Group[meetingUid].AllConn, userUuid)
					fmt.Println("这里删除成功")
					return
				}
			}
			model.Group[meetingUid].Mutex.Unlock()
			return
		case message := <-messageCh: // 注意：此处逻辑正在修改，这里准备维护一个channel链接池，将任务直接丢进去，链接池直接进行消费
			if err != nil {
				initialize.Log.Error("Error Failed to read message:", err)
				break
			}
			// 管理员进行了权限的修改
			if util.IsSpecialMessage(message) {
				err = util.Modify(message)
				if err != nil {
					initialize.Log.Error("Error Modify message error:", err)
					return
				}
			} else {
				msg := model.Message{
					Notice: string(message),
					Name:   username,
				}
				jsonData, err := json.Marshal(msg)
				if err != nil {
					initialize.Log.Error("Error json error:", err)
					return
				}
				content := initialize.Content{
					Message:    jsonData,
					MeetingUid: meetingUid,
				}
				initialize.Queue.SendMsgToQueue(content)
			}
		}
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许跨域请求
		return true
	},
}

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
	"log"
	"net/http"
	"time"
)

func LeaveMeeting(ctx *gin.Context) {
	user := model.User{}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	msg := model.Message{
		Notice: "有成员离开会议室",
		Name:   user.Username,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		initialize.Log.Error("Error json error:", err)
		return
	}
	model.Chan.AllChan[user.UserUid] <- 1 // 给通道传递消息，结束掉conn链接
	// 开始上锁
	model.Group[user.MeetingUid].Mutex.RLock()
	for _, allConn := range model.Group[user.MeetingUid].AllConn {
		err := allConn.Conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			initialize.Log.Error("Error writeMessage error:", err)
			return
		}
	}
	model.Group[user.MeetingUid].Mutex.RUnlock()
	time.Sleep(5 * time.Second)
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
	msg := model.Message{
		Notice: "有成员加入会议室",
		Name:   username,
	}
	jsonData, err := json.Marshal(msg)
	if err != nil {
		initialize.Log.Error("Error json error:", err)
		return
	}
	model.Chan.AllChan[userUuid] = make(chan int)
	model.Group[meetingUid].Mutex.Lock()
	model.Group[meetingUid].AllConn[userUuid] = clientc
	// 有新的链接进来，广播给该会议室所有的人新来人的姓名
	for _, allConn := range model.Group[meetingUid].AllConn {
		err = allConn.Conn.WriteMessage(websocket.TextMessage, jsonData)
		if err != nil {
			initialize.Log.Error("Error Send Message to new all people error:", err)
			return
		}
	}
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
	// 广播给当前加入的新人，会议室的状态
	err = conn.WriteMessage(websocket.TextMessage, controlData)
	if err != nil {
		initialize.Log.Error("Error send Message  to conn error:", err)
		return
	}
	model.Group[meetingUid].Mutex.Unlock()
	messageCh := make(chan []byte)
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

	select {
	case <-model.Chan.AllChan[userUuid]: //说明那边断开了业务
		fmt.Println("断开了业务")
		fmt.Println("chan ", userUuid)
		if ok := initialize.Redis.HDel(meetingUid+"all", userUuid); !ok {
			initialize.Log.Error("Error: 删除用户失败")
			return
		}
		fmt.Println("准备上锁")
		model.Group[meetingUid].Mutex.Lock()
		for _, allConn := range model.Group[meetingUid].AllConn {
			fmt.Println("在这里：", userUuid)
			if allConn.Uid == userUuid {
				delete(model.Group[meetingUid].AllConn, userUuid)
				break
			}
		}
		model.Group[meetingUid].Mutex.Unlock()
		break
	case message := <-messageCh: // 注意：此处逻辑正在修改，这里准备维护一个channel链接池，将任务直接丢进去，链接池直接进行消费
		if err != nil {
			log.Println("Failed to read message:", err)
			break
		}
		model.Group[meetingUid].Mutex.RLock()
		for _, allConn := range model.Group[meetingUid].AllConn {
			err := allConn.Conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				initialize.Log.Error("Error writeMessage error:", err)
				return
			}
		}
		model.Group[meetingUid].Mutex.RUnlock()
	default:
		fmt.Println("开始循环")
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		// 允许跨域请求
		return true
	},
}

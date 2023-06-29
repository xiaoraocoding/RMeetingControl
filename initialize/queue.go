package initialize

import (
	"RMeetingControl/model"
	"fmt"
	"github.com/gorilla/websocket"
	"math/rand"
	"time"
)

type QueueControl struct {
	WorkQueue []chan Content //队列使用channel,传输的数据为[]byte
	QueueSize int            //队列的个数
}

type Content struct {
	Message    []byte
	MeetingUid string
}

var Queue QueueControl

func NewQueue(num int) {
	Queue = QueueControl{
		WorkQueue: make([]chan Content, num),
		QueueSize: num,
	}
}

// 进行广播
func (q QueueControl) SendMsg(msg []byte, MeetingUid string) error {
	fmt.Println("准备打印")
	model.Group[MeetingUid].Mutex.Lock()
	for _, allConn := range model.Group[MeetingUid].AllConn {
		err := allConn.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			Log.Error("Error Conn WriteMessage error:", err)
			return err
		}
	}
	model.Group[MeetingUid].Mutex.Unlock()
	return nil
}

// 开启worker
func (q QueueControl) StartWorker(queue chan Content) error {
	for {
		select {
		case msg := <-queue:
			err := q.SendMsg(msg.Message, msg.MeetingUid)
			if err != nil {
				Log.Error("Error SendMsg error:", err)
				return err
			}
		}
	}
}

func (q QueueControl) StartWorkPool() {
	for i := 0; i < q.QueueSize; i++ {
		q.WorkQueue[i] = make(chan Content, 1024) //这里暂时写死
		fmt.Println("i", "初始化链接")
		go q.StartWorker(q.WorkQueue[i])
	}
}

func (q QueueControl) SendMsgToQueue(content Content) {
	rand.Seed(time.Now().UnixNano())
	// 生成 0 到 9 之间的随机数
	randomNumber := rand.Intn(q.QueueSize)
	fmt.Println("随机数学:", randomNumber)
	q.WorkQueue[randomNumber] <- content
	fmt.Println("发送数据完成")
}

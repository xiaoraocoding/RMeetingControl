package initialize

import (
	"RMeetingControl/model"
	"github.com/gorilla/websocket"
)

type QueueControl struct {
	WorkQueue []chan []byte //队列使用channel,传输的数据为[]byte
	QueueSize int           //队列的个数
}

var Queue QueueControl

func NewQueue() {
	Queue = QueueControl{
		WorkQueue: make([]chan []byte, 0),
	}
}

// 进行广播
func (q QueueControl) SendMsg(msg []byte, MeetingUid string) error {
	model.Group[MeetingUid].Mutex.RLock()
	for _, allConn := range model.Group[MeetingUid].AllConn {
		err := allConn.Conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			Log.Error("Error Conn WriteMessage error:", err)
			return err
		}
	}
	model.Group[MeetingUid].Mutex.RUnlock()
	return nil
}

// 开启worker
func (q QueueControl) StartWorker(queue chan []byte, MeetingUid string) error {
	for {
		select {
		case msg := <-queue:
			err := q.SendMsg(msg, MeetingUid)
			if err != nil {
				Log.Error("Error SendMsg error:", err)
				return err
			}
		}
	}
}

func (q QueueControl) StartWorkPool(MeetingUid string) {
	for i := 0; i < q.QueueSize; i++ {
		q.WorkQueue[i] = make(chan []byte, 1024) //这里暂时写死
		go q.StartWorker(q.WorkQueue[i], MeetingUid)
	}
}

package initialize

// 初始化conn链接组等
import (
	"RMeetingControl/model"
	"sync"
)

func NewMeetingGroup() {
	model.Group = model.MeetingGroup{}
}

func NewMeetingClient(videoEnabled bool, meetingName string, isMute bool) *model.MeetingClient {
	return &model.MeetingClient{
		AllConn:      map[string]model.ClientConnection{},
		IsVideo:      videoEnabled,
		MeetingName:  meetingName,
		IsMute:       isMute,
		IsChangeName: true,
		Mutex:        new(sync.RWMutex),
	}
}

func NewChanControl() {
	model.Chan = model.ChanControl{
		AllChan:   map[string]chan int{},
		HeartChan: map[string]chan int{},
	}
}

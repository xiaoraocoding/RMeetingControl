package model

import (
	"github.com/gorilla/websocket"
	"sync"
)

// 每一个客户端链接的具体结构体
type ClientConnection struct {
	Conn *websocket.Conn
	Name string // 客户端名称
	Uid  string
}

type MeetingClient struct {
	AllConn      map[string]ClientConnection
	IsVideo      bool   // 是否能开启视频
	MeetingName  string // 会议的名称
	IsMute       bool   // 是否开启静音
	IsChangeName bool   // 是否允许修改名称
	Mutex        *sync.RWMutex
}

type MeetingGroup map[string]*MeetingClient

var Group MeetingGroup

type ChanControl struct {
	AllChan   map[string]chan int
	HeartChan map[string]chan int
}

var Chan ChanControl

type Client struct {
	MeetingName string `json:"meeting_name,omitempty"`
	Username    string `json:"username,omitempty"`
	IsMute      bool   `json:"is_mute,omitempty"`
	IsVideo     bool   `json:"is_video,omitempty"`
}

type User struct {
	Username   string `json:"username,omitempty"`
	MeetingUid string `json:"meetingUid,omitempty"`
	UserUid    string `json:"userUid,omitempty"`
}

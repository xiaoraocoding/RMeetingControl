package model

type ChangeVideo struct {
	IsVideo      bool   `json:"is_video"` // 是否能开启视频
	MeetingUid   string `json:"meeting_uid"`
	UserUid      string `json:"user_uid"`
	IsChangeName bool   `json:"is_change_name"`
	IsMute       bool   `json:"is_mute"`
}
type ChangeName struct {
	IsChangeName bool   `json:"is_change_name"` // 是否能开启视频
	MeetingUid   string `json:"meeting_uid"`
	UserUid      string `json:"user_uid"`
}

type ChangeRes struct {
	IsVideo      int `json:"is_video"`
	IsMute       int `json:"is_mute"`
	IsChangeName int `json:"is_change_name"`
}
type ChangeMute struct {
	IsMute     bool   `json:"is_video"` // 是否能开启视频
	MeetingUid string `json:"meeting_uid"`
	UserUid    string `json:"user_uid"`
}

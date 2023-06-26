package model

type ChangeVideo struct {
	IsVideo    bool   `json:"is_video,omitempty"` // 是否能开启视频
	MeetingUid string `json:"meeting_uid,omitempty"`
	UserUid    string `json:"user_uid,omitempty"`
}
type ChangeName struct {
	IsChangeName bool   `json:"is_change_name,omitempty"` // 是否能开启视频
	MeetingUid   string `json:"meeting_uid,omitempty"`
	UserUid      string `json:"user_uid,omitempty"`
}

type ChangeRes struct {
	IsVideo      bool `json:"is_video,omitempty"`
	IsMute       bool `json:"is_mute,omitempty"`
	IsChangeName bool `json:"is_change_name,omitempty"`
}
type ChangeMute struct {
	IsMute     bool   `json:"is_video,omitempty"` // 是否能开启视频
	MeetingUid string `json:"meeting_uid,omitempty"`
	UserUid    string `json:"user_uid,omitempty"`
}

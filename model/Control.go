package model

type Message struct {
	Notice string `json:"notice"`
	Name   string `json:"name"`
}

type Control struct {
	IsMute       bool   `json:"is_mute"`        // 是否开启静音
	IsChangeName bool   `json:"is_change_name"` // 是否允许修改名称
	IsVideo      bool   `json:"is_video"`       // 是否能开启视频
	UserUid      string `json:"user_uid"`       // 自动生成的uid
}

type ChangeIsMute struct {
	Notice string `json:"notice"`
	IsMute bool   `json:"isMute"`
}

type ChangeIsVideo struct {
	Notice  string `json:"notice"`
	IsVideo bool   `json:"isVideo"`
}

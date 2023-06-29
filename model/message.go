package model

type Msg struct {
	UserUid    string `json:"userUid"`
	MeetingUid string `json:"meetingUid"`
}

type Mute struct {
	IsMute     bool   `json:"isMute"`
	UserUid    string `json:"userUid"`
	MeetingUid string `json:"meetingUid"`
}

type Name struct {
	IsName     bool   `json:"isName"`
	UserUid    string `json:"userUid"`
	MeetingUid string `json:"meetingUid"`
}

type Video struct {
	IsVideo    bool   `json:"isVideo"`
	UserUid    string `json:"userUid"`
	MeetingUid string `json:"meetingUid"`
}

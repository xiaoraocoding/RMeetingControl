package util

import (
	"RMeetingControl/initialize"
	"RMeetingControl/model"
	"encoding/json"
	uuid "github.com/satori/go.uuid"
)

func CreateUUid() string {
	u1 := uuid.NewV4()
	return u1.String()
}

// 此处的特定消息为 管理员修改当前会议室的权限
func IsSpecialMessage(message []byte) bool {
	// 判断消息是否为特殊消息的逻辑判断，可以根据实际需求进行定义
	// 这里假设特殊消息的格式为 "特殊消息:" 开头
	return string(message[:4]) == "特殊消息"
}

const (
	Mute  = "1"
	Video = "2"
)

//func IsTypeMessage(message []byte) string {
//	if string(message[4]) == Mute {
//		return
//
//	}
//}

func Modify(message []byte) error {
	if string(message[4]) == Mute {
		var msg model.Mute
		err := json.Unmarshal(message[5:], &msg)
		if err != nil {
			initialize.Log.Error("Error Failed to parse message:", err)
		}
		managerUid := initialize.Redis.Get(msg.MeetingUid + "status")
		if msg.UserUid == managerUid {
			model.Group[msg.MeetingUid].IsMute = msg.IsMute
			mute := model.ChangeIsMute{
				Notice: "会议室状态发送变更",
				IsMute: msg.IsMute,
			}
			jsonData, err := json.Marshal(mute)
			if err != nil {
				initialize.Log.Error("Error json error:", err)
				return err
			}
			conten := initialize.Content{
				Message:    jsonData,
				MeetingUid: msg.MeetingUid,
			}
			initialize.Queue.SendMsgToQueue(conten)
		}
	} else if string(message[4]) == Video {
		var msg model.Video
		err := json.Unmarshal(message, &msg)
		if err != nil {
			initialize.Log.Error("Error Failed to parse message:", err)
		}
		managerUid := initialize.Redis.Get(msg.MeetingUid + "status")
		if msg.UserUid == managerUid {
			model.Group[msg.MeetingUid].IsVideo = msg.IsVideo
			mute := model.ChangeIsVideo{
				Notice:  "会议室状态发送变更",
				IsVideo: msg.IsVideo,
			}
			jsonData, err := json.Marshal(mute)
			if err != nil {
				initialize.Log.Error("Error json error:", err)
				return err
			}
			conten := initialize.Content{
				Message:    jsonData,
				MeetingUid: msg.MeetingUid,
			}
			initialize.Queue.SendMsgToQueue(conten)

		}
	}
	return nil
}


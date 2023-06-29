package api

// 会中管控的部分router
import (
	"RMeetingControl/initialize"
	"RMeetingControl/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ChangeVideoStatus(ctx *gin.Context) {
	change := model.ChangeVideo{}
	if err := ctx.ShouldBindJSON(&change); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	managerUid := initialize.Redis.Get(change.MeetingUid + "status")
	if change.UserUid != managerUid {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前只能管理员能修改会议状态",
		})
	}
	changeRes := model.ChangeRes{IsVideo: change.IsVideo}
	changeData, err := json.Marshal(changeRes)
	if err != nil {
		initialize.Log.Error("Error Marshal changeRes error:", err)
		return
	}
	content := initialize.Content{
		changeData,
		change.MeetingUid,
	}
	initialize.Queue.SendMsgToQueue(content)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功！",
	})
}

func ChangeNameStatus(ctx *gin.Context) {
	change := model.ChangeName{}
	if err := ctx.ShouldBindJSON(&change); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	managerUid := initialize.Redis.Get(change.MeetingUid + "status")
	if change.UserUid != managerUid {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前只能管理员能修改会议状态",
		})
	}
	changeRes := model.ChangeRes{IsChangeName: change.IsChangeName}
	changeData, err := json.Marshal(changeRes)
	if err != nil {
		initialize.Log.Error("Error Marshal changeRes error:", err)
		return
	}
	content := initialize.Content{
		changeData,
		change.MeetingUid,
	}
	initialize.Queue.SendMsgToQueue(content)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功！",
	})
}

func ChangeMuteStatus(ctx *gin.Context) {
	change := model.ChangeMute{}
	if err := ctx.ShouldBindJSON(&change); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	managerUid := initialize.Redis.Get(change.MeetingUid + "status")
	if change.UserUid != managerUid {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前只能管理员能修改会议状态",
		})
	}
	changeRes := model.ChangeRes{IsMute: change.IsMute}
	changeData, err := json.Marshal(changeRes)
	if err != nil {
		initialize.Log.Error("Error Marshal changeRes error:", err)
		return
	}
	content := initialize.Content{
		changeData,
		change.MeetingUid,
	}
	initialize.Queue.SendMsgToQueue(content)
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功！",
	})
}

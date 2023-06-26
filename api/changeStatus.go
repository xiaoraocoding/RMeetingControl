package api

// 会中管控的部分router
import (
	"RMeetingControl/initialize"
	"RMeetingControl/model"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
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
	model.Group[change.MeetingUid].Mutex.Lock()
	model.Group[change.MeetingUid].IsVideo = change.IsVideo
	for _, allConn := range model.Group[change.MeetingUid].AllConn {
		err = allConn.Conn.WriteMessage(websocket.TextMessage, changeData)
		if err != nil {
			initialize.Log.Error("Error change meeting status error:", err)
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "修改会议状态失败！",
			})
			return
		}
	}
	model.Group[change.MeetingUid].Mutex.Unlock()
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
	model.Group[change.MeetingUid].Mutex.Lock()
	model.Group[change.MeetingUid].IsChangeName = change.IsChangeName
	for _, allConn := range model.Group[change.MeetingUid].AllConn {
		err = allConn.Conn.WriteMessage(websocket.TextMessage, changeData)
		if err != nil {
			initialize.Log.Error("Error change meeting status error:", err)
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "修改会议状态失败！",
			})
			return
		}
	}
	model.Group[change.MeetingUid].Mutex.Unlock()
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
	model.Group[change.MeetingUid].Mutex.Lock()
	model.Group[change.MeetingUid].IsMute = change.IsMute
	for _, allConn := range model.Group[change.MeetingUid].AllConn {
		err = allConn.Conn.WriteMessage(websocket.TextMessage, changeData)
		if err != nil {
			initialize.Log.Error("Error change meeting status error:", err)
			ctx.JSON(http.StatusOK, gin.H{
				"msg": "修改会议状态失败！",
			})
			return
		}
	}
	model.Group[change.MeetingUid].Mutex.Unlock()
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "修改成功！",
	})
}

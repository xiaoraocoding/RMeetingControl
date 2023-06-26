package api

// 导出参会人员router
import (
	"RMeetingControl/initialize"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Export(ctx *gin.Context) {
	userUid := ctx.Query("userUid")
	meetingUid := ctx.Query("meetingUid")
	managerUid := initialize.Redis.Get(meetingUid + "status")
	if userUid != managerUid {
		ctx.JSON(http.StatusOK, gin.H{
			"msg": "当前只能管理员能导出参会者",
		})
	}
	res := initialize.Redis.HgetAll(meetingUid + "all")
	var attendee []string
	for _, v := range res {
		attendee = append(attendee, v)
	}
	ctx.JSON(http.StatusOK, gin.H{
		"attendee": attendee,
	})
}

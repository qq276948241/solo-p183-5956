package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

const (
	CodeSuccess           = 0
	CodeParamError        = 10001
	CodeUnauthorized      = 10002
	CodeNotFound          = 10003
	CodeDuplicate         = 10004
	CodeInternalError     = 10005
	CodeScheduleConflict  = 20001
	CodeAppointmentExists = 20002
	CodeInvalidTimeSlot   = 20003
	CodeDoctorNoSchedule  = 20004
)

var codeMsgMap = map[int]string{
	CodeSuccess:           "成功",
	CodeParamError:        "参数错误",
	CodeUnauthorized:      "未授权或token无效",
	CodeNotFound:          "资源不存在",
	CodeDuplicate:         "数据重复",
	CodeInternalError:     "服务器内部错误",
	CodeScheduleConflict:  "排班时间冲突",
	CodeAppointmentExists: "该时段已有预约",
	CodeInvalidTimeSlot:   "无效的时间段",
	CodeDoctorNoSchedule:  "该医生在此时段无排班",
}

func GetMsg(code int) string {
	if msg, ok := codeMsgMap[code]; ok {
		return msg
	}
	return "未知错误"
}

func OK(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code: CodeSuccess,
		Msg:  GetMsg(CodeSuccess),
		Data: data,
	})
}

func Fail(c *gin.Context, httpStatus int, code int) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  GetMsg(code),
	})
}

func FailWithMsg(c *gin.Context, httpStatus int, code int, msg string) {
	c.JSON(httpStatus, Response{
		Code: code,
		Msg:  msg,
	})
}

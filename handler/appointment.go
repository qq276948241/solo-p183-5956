package handler

import (
	"clinic/response"
	"clinic/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	Service *service.AppointmentService
}

func (h *AppointmentHandler) handleErr(c *gin.Context, err error) {
	if biz, ok := err.(*service.BizError); ok {
		if biz.Msg != "" {
			response.FailWithMsg(c, httpStatusForCode(biz.Code), biz.Code, biz.Msg)
		} else {
			response.Fail(c, httpStatusForCode(biz.Code), biz.Code)
		}
		return
	}
	response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
}

func httpStatusForCode(code int) int {
	switch code {
	case response.CodeParamError, response.CodeInvalidTimeSlot, response.CodeDoctorNoSchedule:
		return http.StatusBadRequest
	case response.CodeUnauthorized:
		return http.StatusUnauthorized
	case response.CodeNotFound:
		return http.StatusNotFound
	case response.CodeDuplicate, response.CodeScheduleConflict, response.CodeAppointmentExists:
		return http.StatusConflict
	default:
		return http.StatusInternalServerError
	}
}

type CompleteRequest struct {
	Diagnosis    string `json:"diagnosis"`
	Prescription string `json:"prescription"`
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var a struct {
		PatientID uint   `json:"patient_id"`
		DoctorID  uint   `json:"doctor_id"`
		AppDate   string `json:"app_date"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
		Remark    string `json:"remark"`
	}
	if err := c.ShouldBindJSON(&a); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	req := &service.CreateAppointmentRequest{
		PatientID: a.PatientID,
		DoctorID:  a.DoctorID,
		AppDate:   a.AppDate,
		StartTime: a.StartTime,
		EndTime:   a.EndTime,
		Remark:    a.Remark,
	}
	appt, err := h.Service.Create(req)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, appt)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	appt, err := h.Service.Cancel(uint(id))
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, appt)
}

func (h *AppointmentHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var req CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	appt, record, err := h.Service.Complete(&service.CompleteAppointmentRequest{
		AppointmentID: uint(id),
		Diagnosis:     req.Diagnosis,
		Prescription:  req.Prescription,
	})
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, gin.H{
		"appointment":  appt,
		"visit_record": record,
	})
}

func (h *AppointmentHandler) ListByDate(c *gin.Context) {
	date := c.Query("date")
	list, err := h.Service.ListByDate(date)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, list)
}

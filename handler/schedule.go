package handler

import (
	"clinic/model"
	"clinic/response"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ScheduleHandler struct {
	DB *gorm.DB
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	var s model.DoctorSchedule
	if err := c.ShouldBindJSON(&s); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if s.DoctorID == 0 || s.Weekday < 1 || s.Weekday > 7 || s.StartTime == "" || s.EndTime == "" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "doctor_id、weekday(1-7)、start_time、end_time 均为必填")
		return
	}
	if !isValidTimeFmt(s.StartTime) || !isValidTimeFmt(s.EndTime) {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidTimeSlot)
		return
	}
	if s.StartTime >= s.EndTime {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
		return
	}
	var doctor model.Doctor
	if err := h.DB.First(&doctor, s.DoctorID).Error; err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeNotFound, "医生不存在")
		return
	}
	var conflicts []model.DoctorSchedule
	h.DB.Unscoped().
		Where("doctor_id = ? AND weekday = ? AND start_time < ? AND end_time > ?", s.DoctorID, s.Weekday, s.EndTime, s.StartTime).
		Find(&conflicts)
	if len(conflicts) > 0 {
		response.Fail(c, http.StatusConflict, response.CodeScheduleConflict)
		return
	}
	if err := h.DB.Create(&s).Error; err != nil {
		if isDupErr(err) {
			response.Fail(c, http.StatusConflict, response.CodeDuplicate)
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, s)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	result := h.DB.Delete(&model.DoctorSchedule{}, id)
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	response.OK(c, nil)
}

func (h *ScheduleHandler) ListByDoctor(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctor_id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var list []model.DoctorSchedule
	if err := h.DB.Where("doctor_id = ?", doctorID).Order("weekday, start_time").Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, list)
}

func (h *ScheduleHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var s model.DoctorSchedule
	if err := h.DB.First(&s, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	var input model.DoctorSchedule
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if input.Weekday >= 1 && input.Weekday <= 7 {
		s.Weekday = input.Weekday
	}
	if input.StartTime != "" {
		if !isValidTimeFmt(input.StartTime) {
			response.Fail(c, http.StatusBadRequest, response.CodeInvalidTimeSlot)
			return
		}
		s.StartTime = input.StartTime
	}
	if input.EndTime != "" {
		if !isValidTimeFmt(input.EndTime) {
			response.Fail(c, http.StatusBadRequest, response.CodeInvalidTimeSlot)
			return
		}
		s.EndTime = input.EndTime
	}
	if s.StartTime >= s.EndTime {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
		return
	}
	var conflicts []model.DoctorSchedule
	h.DB.Unscoped().
		Where("id != ? AND doctor_id = ? AND weekday = ? AND start_time < ? AND end_time > ?", s.ID, s.DoctorID, s.Weekday, s.EndTime, s.StartTime).
		Find(&conflicts)
	if len(conflicts) > 0 {
		response.Fail(c, http.StatusConflict, response.CodeScheduleConflict)
		return
	}
	if err := h.DB.Save(&s).Error; err != nil {
		if isDupErr(err) {
			response.Fail(c, http.StatusConflict, response.CodeDuplicate)
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, s)
}

func isValidTimeFmt(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}

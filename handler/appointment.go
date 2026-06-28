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

type AppointmentHandler struct {
	DB *gorm.DB
}

func (h *AppointmentHandler) Create(c *gin.Context) {
	var a model.Appointment
	if err := c.ShouldBindJSON(&a); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if a.PatientID == 0 || a.DoctorID == 0 || a.AppDate == "" || a.StartTime == "" || a.EndTime == "" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "patient_id、doctor_id、app_date、start_time、end_time 均为必填")
		return
	}
	if _, err := time.Parse("2006-01-02", a.AppDate); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "app_date 格式应为 YYYY-MM-DD")
		return
	}
	if !isValidTimeFmt(a.StartTime) || !isValidTimeFmt(a.EndTime) {
		response.Fail(c, http.StatusBadRequest, response.CodeInvalidTimeSlot)
		return
	}
	if a.StartTime >= a.EndTime {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
		return
	}
	var patient model.Patient
	if err := h.DB.First(&patient, a.PatientID).Error; err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeNotFound, "患者不存在")
		return
	}
	var doctor model.Doctor
	if err := h.DB.First(&doctor, a.DoctorID).Error; err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeNotFound, "医生不存在")
		return
	}
	apptDate, _ := time.Parse("2006-01-02", a.AppDate)
	weekday := int(apptDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	var schedules []model.DoctorSchedule
	h.DB.Where("doctor_id = ? AND weekday = ?", a.DoctorID, weekday).Find(&schedules)
	matched := false
	for _, s := range schedules {
		if s.StartTime <= a.StartTime && s.EndTime >= a.EndTime {
			matched = true
			break
		}
	}
	if !matched {
		response.Fail(c, http.StatusBadRequest, response.CodeDoctorNoSchedule)
		return
	}
	var existing model.Appointment
	err := h.DB.Unscoped().
		Where("doctor_id = ? AND app_date = ? AND start_time < ? AND end_time > ? AND status = ?",
			a.DoctorID, a.AppDate, a.EndTime, a.StartTime, "booked").
		First(&existing).Error
	if err == nil {
		response.Fail(c, http.StatusConflict, response.CodeAppointmentExists)
		return
	}
	a.Status = "booked"
	if err := h.DB.Create(&a).Error; err != nil {
		if isDupErr(err) {
			response.Fail(c, http.StatusConflict, response.CodeAppointmentExists)
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	h.DB.Preload("Patient").Preload("Doctor").First(&a, a.ID)
	response.OK(c, a)
}

func (h *AppointmentHandler) Cancel(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var a model.Appointment
	if err := h.DB.First(&a, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	if a.Status != "booked" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "只能取消状态为 booked 的预约")
		return
	}
	a.Status = "cancelled"
	if err := h.DB.Save(&a).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	h.DB.Preload("Patient").Preload("Doctor").First(&a, a.ID)
	response.OK(c, a)
}

type CompleteRequest struct {
	Diagnosis    string `json:"diagnosis"`
	Prescription string `json:"prescription"`
}

func (h *AppointmentHandler) Complete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var a model.Appointment
	if err := h.DB.First(&a, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	if a.Status != "booked" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "只能标记状态为 booked 的预约为已完成")
		return
	}
	var req CompleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if req.Diagnosis == "" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "诊断结论不能为空")
		return
	}
	tx := h.DB.Begin()
	a.Status = "completed"
	if err := tx.Save(&a).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	record := model.VisitRecord{
		AppointmentID: a.ID,
		Diagnosis:     req.Diagnosis,
		Prescription:  req.Prescription,
	}
	if err := tx.Create(&record).Error; err != nil {
		tx.Rollback()
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	tx.Commit()
	h.DB.Preload("Patient").Preload("Doctor").First(&a, a.ID)
	h.DB.Preload("Appointment").First(&record, record.ID)
	response.OK(c, gin.H{
		"appointment":  a,
		"visit_record": record,
	})
}

func (h *AppointmentHandler) ListByDate(c *gin.Context) {
	date := c.Query("date")
	if _, err := time.Parse("2006-01-02", date); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "date 参数格式应为 YYYY-MM-DD")
		return
	}
	var list []model.Appointment
	if err := h.DB.Preload("Patient").Preload("Doctor").
		Where("app_date = ?", date).
		Order("start_time").Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, list)
}

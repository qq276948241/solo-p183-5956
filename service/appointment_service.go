package service

import (
	"clinic/model"
	"clinic/response"
	"errors"
	"time"

	"gorm.io/gorm"
)

type AppointmentService struct {
	DB *gorm.DB
}

func NewAppointmentService(db *gorm.DB) *AppointmentService {
	return &AppointmentService{DB: db}
}

type CreateAppointmentRequest struct {
	PatientID uint
	DoctorID  uint
	AppDate   string
	StartTime string
	EndTime   string
	Remark    string
}

type CompleteAppointmentRequest struct {
	AppointmentID uint
	Diagnosis     string
	Prescription  string
}

func (s *AppointmentService) Create(req *CreateAppointmentRequest) (*model.Appointment, error) {
	if req.PatientID == 0 || req.DoctorID == 0 || req.AppDate == "" || req.StartTime == "" || req.EndTime == "" {
		return nil, NewBizError(response.CodeParamError, "patient_id、doctor_id、app_date、start_time、end_time 均为必填")
	}
	if _, err := time.Parse("2006-01-02", req.AppDate); err != nil {
		return nil, NewBizError(response.CodeParamError, "app_date 格式应为 YYYY-MM-DD")
	}
	if !isValidTimeFmt(req.StartTime) || !isValidTimeFmt(req.EndTime) {
		return nil, NewBizError(response.CodeInvalidTimeSlot, "")
	}
	if req.StartTime >= req.EndTime {
		return nil, NewBizError(response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
	}

	if err := s.DB.First(&model.Patient{}, req.PatientID).Error; err != nil {
		return nil, NewBizError(response.CodeNotFound, "患者不存在")
	}
	if err := s.DB.First(&model.Doctor{}, req.DoctorID).Error; err != nil {
		return nil, NewBizError(response.CodeNotFound, "医生不存在")
	}

	if err := s.checkDoctorSchedule(req.DoctorID, req.AppDate, req.StartTime, req.EndTime); err != nil {
		return nil, err
	}

	tx := s.DB.Begin()
	if err := s.checkAppointmentConflictForUpdate(tx, req.DoctorID, req.AppDate, req.StartTime, req.EndTime); err != nil {
		tx.Rollback()
		return nil, err
	}

	a := &model.Appointment{
		PatientID: req.PatientID,
		DoctorID:  req.DoctorID,
		AppDate:   req.AppDate,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    "booked",
		Remark:    req.Remark,
	}
	if err := tx.Create(a).Error; err != nil {
		tx.Rollback()
		if isDupErr(err) {
			return nil, NewBizError(response.CodeAppointmentExists, "")
		}
		return nil, NewBizError(response.CodeInternalError, "")
	}
	tx.Commit()

	s.DB.Preload("Patient").Preload("Doctor").First(a, a.ID)
	return a, nil
}

func (s *AppointmentService) Cancel(apptID uint) (*model.Appointment, error) {
	var a model.Appointment
	if err := s.DB.First(&a, apptID).Error; err != nil {
		return nil, NewBizError(response.CodeNotFound, "")
	}
	if a.Status != "booked" {
		return nil, NewBizError(response.CodeParamError, "只能取消状态为 booked 的预约")
	}
	a.Status = "cancelled"
	if err := s.DB.Save(&a).Error; err != nil {
		return nil, NewBizError(response.CodeInternalError, "")
	}
	s.DB.Preload("Patient").Preload("Doctor").First(&a, a.ID)
	return &a, nil
}

func (s *AppointmentService) Complete(req *CompleteAppointmentRequest) (*model.Appointment, *model.VisitRecord, error) {
	var a model.Appointment
	if err := s.DB.First(&a, req.AppointmentID).Error; err != nil {
		return nil, nil, NewBizError(response.CodeNotFound, "")
	}
	if a.Status != "booked" {
		return nil, nil, NewBizError(response.CodeParamError, "只能标记状态为 booked 的预约为已完成")
	}
	if req.Diagnosis == "" {
		return nil, nil, NewBizError(response.CodeParamError, "诊断结论不能为空")
	}

	tx := s.DB.Begin()
	a.Status = "completed"
	if err := tx.Save(&a).Error; err != nil {
		tx.Rollback()
		return nil, nil, NewBizError(response.CodeInternalError, "")
	}
	record := &model.VisitRecord{
		AppointmentID: a.ID,
		Diagnosis:     req.Diagnosis,
		Prescription:  req.Prescription,
	}
	if err := tx.Create(record).Error; err != nil {
		tx.Rollback()
		return nil, nil, NewBizError(response.CodeInternalError, "")
	}
	tx.Commit()

	s.DB.Preload("Patient").Preload("Doctor").First(&a, a.ID)
	s.DB.Preload("Appointment").First(record, record.ID)
	return &a, record, nil
}

func (s *AppointmentService) ListByDate(date string) ([]model.Appointment, error) {
	if _, err := time.Parse("2006-01-02", date); err != nil {
		return nil, NewBizError(response.CodeParamError, "date 参数格式应为 YYYY-MM-DD")
	}
	var list []model.Appointment
	if err := s.DB.Preload("Patient").Preload("Doctor").
		Where("app_date = ?", date).
		Order("start_time").Find(&list).Error; err != nil {
		return nil, NewBizError(response.CodeInternalError, "")
	}
	return list, nil
}

func (s *AppointmentService) checkDoctorSchedule(doctorID uint, appDate, startTime, endTime string) error {
	apptDate, _ := time.Parse("2006-01-02", appDate)
	weekday := int(apptDate.Weekday())
	if weekday == 0 {
		weekday = 7
	}
	var schedules []model.DoctorSchedule
	s.DB.Where("doctor_id = ? AND weekday = ?", doctorID, weekday).Find(&schedules)
	matched := false
	for _, sc := range schedules {
		if sc.StartTime <= startTime && sc.EndTime >= endTime {
			matched = true
			break
		}
	}
	if !matched {
		return NewBizError(response.CodeDoctorNoSchedule, "")
	}
	return nil
}

func (s *AppointmentService) checkAppointmentConflictForUpdate(tx *gorm.DB, doctorID uint, appDate, startTime, endTime string) error {
	var existing model.Appointment
	err := tx.Set("gorm:query_option", "FOR UPDATE").
		Where("doctor_id = ? AND app_date = ? AND start_time < ? AND end_time > ? AND status = ?",
			doctorID, appDate, endTime, startTime, "booked").
		First(&existing).Error
	if err == nil {
		return NewBizError(response.CodeAppointmentExists, "")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return NewBizError(response.CodeInternalError, "")
	}
	return nil
}

func isValidTimeFmt(t string) bool {
	_, err := time.Parse("15:04", t)
	return err == nil
}

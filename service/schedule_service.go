package service

import (
	"clinic/model"
	"clinic/response"

	"gorm.io/gorm"
)

type ScheduleService struct {
	DB *gorm.DB
}

func NewScheduleService(db *gorm.DB) *ScheduleService {
	return &ScheduleService{DB: db}
}

func (s *ScheduleService) Create(doctorID uint, weekday int, startTime, endTime string) (*model.DoctorSchedule, error) {
	if doctorID == 0 || weekday < 1 || weekday > 7 || startTime == "" || endTime == "" {
		return nil, NewBizError(response.CodeParamError, "doctor_id、weekday(1-7)、start_time、end_time 均为必填")
	}
	if !isValidTimeFmt(startTime) || !isValidTimeFmt(endTime) {
		return nil, NewBizError(response.CodeInvalidTimeSlot, "")
	}
	if startTime >= endTime {
		return nil, NewBizError(response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
	}
	if err := s.DB.First(&model.Doctor{}, doctorID).Error; err != nil {
		return nil, NewBizError(response.CodeNotFound, "医生不存在")
	}
	var conflicts []model.DoctorSchedule
	s.DB.Unscoped().
		Where("doctor_id = ? AND weekday = ? AND start_time < ? AND end_time > ?",
			doctorID, weekday, endTime, startTime).
		Find(&conflicts)
	if len(conflicts) > 0 {
		return nil, NewBizError(response.CodeScheduleConflict, "")
	}
	sc := &model.DoctorSchedule{
		DoctorID:  doctorID,
		Weekday:   weekday,
		StartTime: startTime,
		EndTime:   endTime,
	}
	if err := s.DB.Create(sc).Error; err != nil {
		if isDupErr(err) {
			return nil, NewBizError(response.CodeDuplicate, "")
		}
		return nil, NewBizError(response.CodeInternalError, "")
	}
	return sc, nil
}

func (s *ScheduleService) Update(id uint, weekday int, startTime, endTime string) (*model.DoctorSchedule, error) {
	var sc model.DoctorSchedule
	if err := s.DB.First(&sc, id).Error; err != nil {
		return nil, NewBizError(response.CodeNotFound, "")
	}
	if weekday >= 1 && weekday <= 7 {
		sc.Weekday = weekday
	}
	if startTime != "" {
		if !isValidTimeFmt(startTime) {
			return nil, NewBizError(response.CodeInvalidTimeSlot, "")
		}
		sc.StartTime = startTime
	}
	if endTime != "" {
		if !isValidTimeFmt(endTime) {
			return nil, NewBizError(response.CodeInvalidTimeSlot, "")
		}
		sc.EndTime = endTime
	}
	if sc.StartTime >= sc.EndTime {
		return nil, NewBizError(response.CodeInvalidTimeSlot, "开始时间必须早于结束时间")
	}
	var conflicts []model.DoctorSchedule
	s.DB.Unscoped().
		Where("id != ? AND doctor_id = ? AND weekday = ? AND start_time < ? AND end_time > ?",
			sc.ID, sc.DoctorID, sc.Weekday, sc.EndTime, sc.StartTime).
		Find(&conflicts)
	if len(conflicts) > 0 {
		return nil, NewBizError(response.CodeScheduleConflict, "")
	}
	if err := s.DB.Save(&sc).Error; err != nil {
		if isDupErr(err) {
			return nil, NewBizError(response.CodeDuplicate, "")
		}
		return nil, NewBizError(response.CodeInternalError, "")
	}
	return &sc, nil
}

func (s *ScheduleService) Delete(id uint) error {
	result := s.DB.Delete(&model.DoctorSchedule{}, id)
	if result.RowsAffected == 0 {
		return NewBizError(response.CodeNotFound, "")
	}
	return nil
}

func (s *ScheduleService) ListByDoctor(doctorID uint) ([]model.DoctorSchedule, error) {
	var list []model.DoctorSchedule
	if err := s.DB.Where("doctor_id = ?", doctorID).
		Order("weekday, start_time").Find(&list).Error; err != nil {
		return nil, NewBizError(response.CodeInternalError, "")
	}
	return list, nil
}

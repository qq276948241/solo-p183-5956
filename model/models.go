package model

import (
	"time"

	"gorm.io/gorm"
)

type Patient struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Phone     string         `json:"phone" gorm:"type:varchar(20);not null;uniqueIndex"`
	Gender    string         `json:"gender" gorm:"type:varchar(10)"`
	BirthDate *time.Time     `json:"birth_date,omitempty" gorm:"type:date"`
	Address   string         `json:"address" gorm:"type:varchar(200)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Doctor struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"type:varchar(50);not null"`
	Phone     string         `json:"phone" gorm:"type:varchar(20)"`
	Dept      string         `json:"dept" gorm:"type:varchar(50)"`
	Title     string         `json:"title" gorm:"type:varchar(50)"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type DoctorSchedule struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	DoctorID  uint           `json:"doctor_id" gorm:"not null;uniqueIndex:idx_doc_week_start_end"`
	Weekday   int            `json:"weekday" gorm:"not null;uniqueIndex:idx_doc_week_start_end;comment:1=周一...7=周日"`
	StartTime string         `json:"start_time" gorm:"type:varchar(10);not null;uniqueIndex:idx_doc_week_start_end;comment:格式 HH:mm"`
	EndTime   string         `json:"end_time" gorm:"type:varchar(10);not null;uniqueIndex:idx_doc_week_start_end;comment:格式 HH:mm"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type Appointment struct {
	ID         uint           `json:"id" gorm:"primaryKey"`
	PatientID  uint           `json:"patient_id" gorm:"not null;uniqueIndex:idx_appt_unique"`
	DoctorID   uint           `json:"doctor_id" gorm:"not null;uniqueIndex:idx_appt_unique"`
	AppDate    string         `json:"app_date" gorm:"type:date;not null;uniqueIndex:idx_appt_unique"`
	StartTime  string         `json:"start_time" gorm:"type:varchar(10);not null;uniqueIndex:idx_appt_unique;comment:格式 HH:mm"`
	EndTime    string         `json:"end_time" gorm:"type:varchar(10);not null;uniqueIndex:idx_appt_unique;comment:格式 HH:mm"`
	Status     string         `json:"status" gorm:"type:varchar(20);not null;default:booked;comment:booked/cancelled/completed"`
	Remark     string         `json:"remark" gorm:"type:varchar(500)"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `json:"-" gorm:"index"`
	Patient    Patient        `json:"patient" gorm:"foreignKey:PatientID"`
	Doctor     Doctor         `json:"doctor" gorm:"foreignKey:DoctorID"`
}

func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&Patient{}, &Doctor{}, &DoctorSchedule{}, &Appointment{})
}

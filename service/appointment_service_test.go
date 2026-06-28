package service

import (
	"clinic/model"
	"fmt"
	"sync"
	"testing"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
	dsn := "clinic:clinic123@tcp(127.0.0.1:3306)/clinic_test?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Skip("跳过并发测试：无法连接测试数据库")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(20)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Minute)

	db.Exec("DROP TABLE IF EXISTS visit_records, appointments, doctor_schedules, doctors, patients")
	if err := db.AutoMigrate(&model.Patient{}, &model.Doctor{}, &model.DoctorSchedule{}, &model.Appointment{}, &model.VisitRecord{}); err != nil {
		t.Fatalf("迁移失败: %v", err)
	}
	if !db.Migrator().HasColumn(&model.Appointment{}, "booked_slot_key") {
		if err := db.Exec(`
			ALTER TABLE appointments
			ADD COLUMN booked_slot_key VARCHAR(100) GENERATED ALWAYS AS (
				CASE WHEN status = 'booked'
					THEN CONCAT(doctor_id, '-', app_date, '-', start_time, '-', end_time)
					ELSE NULL
				END
			) STORED
		`).Error; err != nil {
			t.Fatalf("添加生成列失败: %v", err)
		}
	}
	if !db.Migrator().HasIndex(&model.Appointment{}, "idx_booked_slot") {
		if err := db.Exec(`CREATE UNIQUE INDEX idx_booked_slot ON appointments(booked_slot_key)`).Error; err != nil {
			t.Fatalf("添加唯一索引失败: %v", err)
		}
	}

	return db
}

func TestConcurrentAppointmentCreate(t *testing.T) {
	db := setupTestDB(t)

	doctor := model.Doctor{Name: "测试医生", Dept: "全科", Title: "主治医师"}
	db.Create(&doctor)

	schedule := model.DoctorSchedule{
		DoctorID:  doctor.ID,
		Weekday:   1,
		StartTime: "08:00",
		EndTime:   "12:00",
	}
	db.Create(&schedule)

	apptDate := nextMonday()
	numPatients := 10
	patients := make([]model.Patient, numPatients)
	for i := 0; i < numPatients; i++ {
		patients[i] = model.Patient{
			Name:  fmt.Sprintf("患者%d", i),
			Phone: fmt.Sprintf("138%08d", i),
		}
		db.Create(&patients[i])
	}

	svc := NewAppointmentService(db)

	var wg sync.WaitGroup
	results := make(chan error, numPatients)

	for i := 0; i < numPatients; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			_, err := svc.Create(&CreateAppointmentRequest{
				PatientID: patients[idx].ID,
				DoctorID:  doctor.ID,
				AppDate:   apptDate,
				StartTime: "08:00",
				EndTime:   "08:30",
			})
			results <- err
		}(i)
	}
	wg.Wait()
	close(results)

	successCount := 0
	conflictCount := 0
	for err := range results {
		if err == nil {
			successCount++
		} else if biz, ok := err.(*BizError); ok && biz.Code == 20002 {
			conflictCount++
		} else {
			t.Errorf("意外的错误: %v", err)
		}
	}

	if successCount != 1 {
		t.Errorf("期望1个预约成功，实际%d个成功", successCount)
	}
	if conflictCount != numPatients-1 {
		t.Errorf("期望%d个冲突，实际%d个冲突", numPatients-1, conflictCount)
	}

	var count int64
	db.Model(&model.Appointment{}).Where("doctor_id = ? AND app_date = ? AND start_time = ? AND status = ?",
		doctor.ID, apptDate, "08:00", "booked").Count(&count)
	if count != 1 {
		t.Errorf("数据库中应有1条booked记录，实际%d条", count)
	}

	t.Logf("并发测试通过：%d个请求，1个成功，%d个冲突拒绝", numPatients, conflictCount)
}

func TestConcurrentOverlappingSlotCreate(t *testing.T) {
	db := setupTestDB(t)

	doctor := model.Doctor{Name: "测试医生2", Dept: "全科", Title: "主治医师"}
	db.Create(&doctor)

	schedule := model.DoctorSchedule{
		DoctorID:  doctor.ID,
		Weekday:   1,
		StartTime: "08:00",
		EndTime:   "12:00",
	}
	db.Create(&schedule)

	apptDate := nextMonday()
	p1 := model.Patient{Name: "患者A", Phone: "13900001111"}
	p2 := model.Patient{Name: "患者B", Phone: "13900002222"}
	db.Create(&p1)
	db.Create(&p2)

	svc := NewAppointmentService(db)

	_, err1 := svc.Create(&CreateAppointmentRequest{
		PatientID: p1.ID,
		DoctorID:  doctor.ID,
		AppDate:   apptDate,
		StartTime: "08:00",
		EndTime:   "08:30",
	})
	if err1 != nil {
		t.Fatalf("首次预约应该成功: %v", err1)
	}

	_, err2 := svc.Create(&CreateAppointmentRequest{
		PatientID: p2.ID,
		DoctorID:  doctor.ID,
		AppDate:   apptDate,
		StartTime: "08:15",
		EndTime:   "08:45",
	})
	if err2 == nil {
		t.Fatal("重叠时段的预约应该被拒绝")
	}

	t.Logf("重叠时段测试通过：首个预约成功，重叠预约被拒绝（%v）", err2)
}

func nextMonday() string {
	now := time.Now()
	for i := 1; i <= 7; i++ {
		d := now.AddDate(0, 0, i)
		if d.Weekday() == time.Monday {
			return d.Format("2006-01-02")
		}
	}
	return now.Format("2006-01-02")
}

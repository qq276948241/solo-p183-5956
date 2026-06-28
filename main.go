package main

import (
	"clinic/config"
	"clinic/handler"
	"clinic/middleware"
	"clinic/model"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(mysql.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("连接数据库失败: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	if err := model.AutoMigrate(db); err != nil {
		log.Fatalf("数据库迁移失败: %v", err)
	}
	fmt.Println("数据库迁移完成")

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	auth := r.Group("")
	auth.Use(middleware.TokenAuth(cfg))
	{
		patientH := handler.PatientHandler{DB: db}
		auth.POST("/patients", patientH.Create)
		auth.PUT("/patients/:id", patientH.Update)
		auth.DELETE("/patients/:id", patientH.Delete)
		auth.GET("/patients/:id", patientH.GetByID)
		auth.GET("/patients", patientH.List)

		doctorH := handler.DoctorHandler{DB: db}
		auth.POST("/doctors", doctorH.Create)
		auth.GET("/doctors", doctorH.List)
		auth.GET("/doctors/:id", doctorH.GetByID)

		scheduleH := handler.ScheduleHandler{DB: db}
		auth.POST("/doctors/:doctor_id/schedules", scheduleH.Create)
		auth.PUT("/schedules/:id", scheduleH.Update)
		auth.DELETE("/schedules/:id", scheduleH.Delete)
		auth.GET("/doctors/:doctor_id/schedules", scheduleH.ListByDoctor)

		appointmentH := handler.AppointmentHandler{DB: db}
		auth.POST("/appointments", appointmentH.Create)
		auth.PUT("/appointments/:id/cancel", appointmentH.Cancel)
		auth.PUT("/appointments/:id/complete", appointmentH.Complete)
		auth.GET("/appointments", appointmentH.ListByDate)
	}

	addr := ":" + cfg.ServerPort
	fmt.Printf("服务启动于 http://localhost%s\n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatalf("服务启动失败: %v", err)
	}
}

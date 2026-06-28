package handler

import (
	"clinic/model"
	"clinic/response"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type DoctorHandler struct {
	DB *gorm.DB
}

func (h *DoctorHandler) Create(c *gin.Context) {
	var d model.Doctor
	if err := c.ShouldBindJSON(&d); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if d.Name == "" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "医生姓名不能为空")
		return
	}
	if err := h.DB.Create(&d).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, d)
}

func (h *DoctorHandler) List(c *gin.Context) {
	var list []model.Doctor
	if err := h.DB.Order("id desc").Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, list)
}

func (h *DoctorHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	var d model.Doctor
	if err := h.DB.First(&d, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	response.OK(c, d)
}

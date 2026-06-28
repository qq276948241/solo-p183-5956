package handler

import (
	"clinic/model"
	"clinic/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PatientHandler struct {
	DB *gorm.DB
}

func (h *PatientHandler) Create(c *gin.Context) {
	var p model.Patient
	if err := c.ShouldBindJSON(&p); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if p.Name == "" || p.Phone == "" {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "姓名和手机号不能为空")
		return
	}
	if err := h.DB.Create(&p).Error; err != nil {
		if isDupErr(err) {
			response.Fail(c, http.StatusConflict, response.CodeDuplicate)
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, p)
}

func (h *PatientHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var p model.Patient
	if err := h.DB.First(&p, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	var input model.Patient
	if err := c.ShouldBindJSON(&input); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	if input.Name != "" {
		p.Name = input.Name
	}
	if input.Phone != "" {
		p.Phone = input.Phone
	}
	if input.Gender != "" {
		p.Gender = input.Gender
	}
	if input.BirthDate != nil {
		p.BirthDate = input.BirthDate
	}
	if input.Address != "" {
		p.Address = input.Address
	}
	if err := h.DB.Save(&p).Error; err != nil {
		if isDupErr(err) {
			response.Fail(c, http.StatusConflict, response.CodeDuplicate)
			return
		}
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, p)
}

func (h *PatientHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	result := h.DB.Delete(&model.Patient{}, id)
	if result.RowsAffected == 0 {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	response.OK(c, nil)
}

func (h *PatientHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var p model.Patient
	if err := h.DB.First(&p, id).Error; err != nil {
		response.Fail(c, http.StatusNotFound, response.CodeNotFound)
		return
	}
	response.OK(c, p)
}

func (h *PatientHandler) List(c *gin.Context) {
	var list []model.Patient
	query := h.DB.Order("id desc")
	if phone := c.Query("phone"); phone != "" {
		query = query.Where("phone LIKE ?", "%"+phone+"%")
	}
	if name := c.Query("name"); name != "" {
		query = query.Where("name LIKE ?", "%"+name+"%")
	}
	if err := query.Find(&list).Error; err != nil {
		response.Fail(c, http.StatusInternalServerError, response.CodeInternalError)
		return
	}
	response.OK(c, list)
}

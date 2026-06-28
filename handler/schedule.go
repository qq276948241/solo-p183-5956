package handler

import (
	"clinic/response"
	"clinic/service"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ScheduleHandler struct {
	Service *service.ScheduleService
}

func (h *ScheduleHandler) Create(c *gin.Context) {
	doctorID, err := strconv.ParseUint(c.Param("doctor_id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	var s struct {
		Weekday   int    `json:"weekday"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&s); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	sc, err := h.Service.Create(uint(doctorID), s.Weekday, s.StartTime, s.EndTime)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, sc)
}

func (h *ScheduleHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		response.Fail(c, http.StatusBadRequest, response.CodeParamError)
		return
	}
	if err := h.Service.Delete(uint(id)); err != nil {
		h.handleErr(c, err)
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
	list, err := h.Service.ListByDoctor(uint(doctorID))
	if err != nil {
		h.handleErr(c, err)
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
	var s struct {
		Weekday   int    `json:"weekday"`
		StartTime string `json:"start_time"`
		EndTime   string `json:"end_time"`
	}
	if err := c.ShouldBindJSON(&s); err != nil {
		response.FailWithMsg(c, http.StatusBadRequest, response.CodeParamError, "请求参数无效: "+err.Error())
		return
	}
	sc, err := h.Service.Update(uint(id), s.Weekday, s.StartTime, s.EndTime)
	if err != nil {
		h.handleErr(c, err)
		return
	}
	response.OK(c, sc)
}

func (h *ScheduleHandler) handleErr(c *gin.Context, err error) {
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

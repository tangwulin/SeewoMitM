package main

import (
	"github.com/gin-gonic/gin"
)
import "SeewoMitM/model"

// GetScreensaverContentHandler 屏保列表接口
// @Summary 屏保列表接口
// @Description 获取所有屏保列表
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} model.GetScreensaverContentResponse
// @Router /api/v1/getScreensaverContent [get]
func GetScreensaverContentHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
		"data":    GetScreensaverContentList(),
		"code":    200,
	})
}

// GetScreensaverContentByIDHandler 获取屏保详情接口
// @Summary 获取屏保详情接口
// @Description 根据屏保id获取屏保详情
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Param id body model.GetScreensaverContentByIDRequest true "屏保id"
// @Success 200 {object} model.GetScreensaverContentByIDResponse
// @Router /api/v1/getScreensaverContentByID [post]
func GetScreensaverContentByIDHandler(c *gin.Context) {
	var req model.GetScreensaverContentByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
			"data":    nil,
			"code":    400,
		})
		return
	}

	content, err := GetScreensaverContentByID(req.ID)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
			"data":    nil,
			"code":    500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    content,
		"code":    200,
	})
}

// AddScreensaverContentHandler 添加屏保接口
// @Summary 添加屏保接口
// @Description 添加屏保
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Param content body model.AddScreensaverContentRequest true "屏保内容"
// @Success 200 {object} model.AddScreensaverContentResponse
// @Router /api/v1/addScreensaverContent [post]
func AddScreensaverContentHandler(c *gin.Context) {
	var req model.AddScreensaverContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
			"data":    nil,
			"code":    400,
		})
		return
	}

	id, err := AddScreensaverContent(*req.Content)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
			"data":    nil,
			"code":    500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    id,
		"code":    200,
	})
}

// UpdateScreensaverContentHandler 更新屏保接口
// @Summary 更新屏保接口
// @Description 更新屏保
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Param id body model.UpdateScreensaverContentRequest true "屏保id"
// @Param content body model.UpdateScreensaverContentRequest true "屏保内容"
// @Success 200 {object} model.UpdateScreensaverContentResponse
// @Router /api/v1/updateScreensaverContent [post]
func UpdateScreensaverContentHandler(c *gin.Context) {
	var req model.UpdateScreensaverContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
			"data":    nil,
			"code":    400,
		})
		return
	}

	err := UpdateScreensaverContent(req.ID, *req.Content)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
			"code":    500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"code":    200,
	})
}

// DeleteScreensaverContentHandler 删除屏保接口
// @Summary 删除屏保接口
// @Description 删除屏保
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Param id body model.DeleteScreensaverContentRequest true "屏保id"
// @Success 200 {object} model.DeleteScreensaverContentResponse
// @Router /api/v1/deleteScreensaverContent [post]
func DeleteScreensaverContentHandler(c *gin.Context) {
	var req model.DeleteScreensaverContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
			"code":    400,
		})
		return
	}

	err := DeleteScreensaverContent(req.ID)

	if err != nil {
		c.JSON(200, gin.H{
			"message": err,
			"code":    500,
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"code":    200,
	})
}

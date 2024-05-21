package main

import "github.com/gin-gonic/gin"
import "SeewoMitM/model"

// GetScreensaverContentHandler 屏保列表接口
// @Summary 屏保列表接口
// @Description 获取所有屏保列表
// @Tags 屏保相关接口
// @Accept application/json
// @Produce application/json
// @Success 200 {object} []ScreensaverContent
// @Router /api/v1/getScreensaverContent [get]
func GetScreensaverContentHandler(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "ok",
		"data":    globalConfig.ScreensaverConfig.Contents,
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
// @Success 200 {object} ScreensaverContent
// @Router /api/v1/getScreensaverContentByID [post]
func GetScreensaverContentByIDHandler(c *gin.Context) {
	var req model.GetScreensaverContentByIDRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{
			"message": "参数错误",
			"code":    400,
		})
		return
	}
	content := globalConfig.ScreensaverConfig.GetContentByID(req.ID)
	if content == nil {
		c.JSON(200, gin.H{
			"message": "未找到该屏保",
			"code":    404,
		})
	}

	c.JSON(200, gin.H{
		"message": "ok",
		"data":    content,
		"code":    200,
	})
}

package main

import (
	"SeewoMitM/internal/log"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strconv"
)

// @title SeewoMitM管理API
// @version 1.0
// @description 用于为SeewoMitM的WebUI提供管理API
// @termsOfService http://swagger.io/terms/

// @license.name GNU General Public License v3.0
// @license.url https://www.gnu.org/licenses/gpl-3.0.html

// @host 127.0.0.1:11451
func LaunchManageServer(port int) error {
	engine := gin.New()

	engine.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	api := engine.Group("/api")
	{
		v1 := api.Group("/v1")
		{
			v1.GET("/ping", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "pong",
					"code":    200,
				})
			})

			v1.GET("/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "status",
					"code":    200,
				})
			})

			v1.GET("/config", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "config",
					"data":    GetScreensaverContent(),
					"code":    200,
				})
			})

			v1.GET("/getScreensaverPayload", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "getScreensaverPayload",
					"data":    GenScreensaverPayload(),
					"code":    200,
				})
			})

			v1.GET("/getScreensaverContent", GetScreensaverContentHandler)
			v1.POST("/getScreensaverContentByID", GetScreensaverContentByIDHandler)
			v1.POST("/addScreensaverContent", AddScreensaverContentHandler)
			v1.POST("/updateScreensaverContent", UpdateScreensaverContentHandler)
			v1.POST("/deleteScreensaverContent", DeleteScreensaverContentHandler)
		}
	}

	engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	err := engine.Run(":" + strconv.Itoa(port))
	if err != nil {
		log.WithFields(log.Fields{"type": "ManageServer"}).Error("failed to launch manage server")
		return err
	}
	return nil
}

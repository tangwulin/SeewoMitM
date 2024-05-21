package main

import (
	"SeewoMitM/internal/log"
	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"strconv"
)

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

package manage

import (
	"SeewoMitM/internal/log"
	"SeewoMitM/internal/screensaver"
	"github.com/gin-gonic/gin"
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
					"code":    "200",
				})
			})

			v1.GET("/status", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "status",
					"code":    "200",
				})
			})

			v1.GET("/config", func(c *gin.Context) {
				c.JSON(200, gin.H{
					"message": "config",
					"data":    "config",
					"code":    screensaver.GetScreensaverData(),
				})
			})
		}
	}

	err := engine.Run(":" + strconv.Itoa(port))
	if err != nil {
		log.WithFields(log.Fields{"type": "ManageServer"}).Error("failed to launch manage server")
		return err
	}
	return nil
}

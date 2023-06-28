package main

import (
	"github.com/algrt-hm/steamdeck-intermediary/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// playerMacAddressConst is the MAC address of the Squeezebox player we want to control
	const touchMacAddress = "00:04:20:23:a1:b5"
	const radioMacAddress = "00:04:20:2b:76:f6"

	gin.SetMode(gin.ReleaseMode)
	// instantiate server, includes logger and recovery middleware
	server := gin.Default()

	// Test / status endpoint
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// Simple group: lms
	lms := server.Group("/lms")
	{
		lms.GET("/play", func(c *gin.Context) {
			code, body := service.LmsPost(touchMacAddress, service.LmsPlay)
			c.String(code, body)
		})
		lms.GET("/pause", func(c *gin.Context) {
			code, body := service.LmsPost(touchMacAddress, service.LmsPause)
			c.String(code, body)
		})

		touch := lms.Group("/touch")
		{
			touch.GET("/voldown", func(c *gin.Context) {
				code, body := service.LmsPost(touchMacAddress, service.LmsVolumeDown)
				c.String(code, body)
			})
			touch.GET("/volup", func(c *gin.Context) {
				code, body := service.LmsPost(touchMacAddress, service.LmsVolumeUp)
				c.String(code, body)
			})
		}

		radio := lms.Group("/radio")
		{
			radio.GET("/voldown", func(c *gin.Context) {
				code, body := service.LmsPost(radioMacAddress, service.LmsVolumeDown)
				c.String(code, body)
			})
			radio.GET("/volup", func(c *gin.Context) {
				code, body := service.LmsPost(radioMacAddress, service.LmsVolumeUp)
				c.String(code, body)
			})
		}
	}

	server.Run(":8080")
}

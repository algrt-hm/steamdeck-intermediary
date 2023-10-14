package main

import (
	// service here contains the LmsPost function used to generate the JSON payloads and POST them to the LMS server
	"github.com/algrt-hm/steamdeck-intermediary/service"
	"github.com/gin-gonic/gin"
)

func main() {
	// playerMacAddressConst is the MAC address of the Squeezebox player we want to control
	// here by way of example I have two devices, a Squeezebox Touch and a Squeezebox Radio
	// each with a different MAC address, as one would expect
	const touchMacAddress = "00:04:20:23:a1:b5"
	const radioMacAddress = "00:04:20:2b:76:f6"

	// set gin to release mode i.e. not debug
	gin.SetMode(gin.ReleaseMode)
	// instantiate server, includes logger and recovery middleware
	server := gin.Default()

	// Test / status endpoint
	// we can hit this to check the server is running
	server.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "Hello World!",
		})
	})

	// Simple group: lms
	lms := server.Group("/lms")
	{
		// play and pause endpoints
		// these will affect both players because within the LMS configuration
		// the two are synch'd
		lms.GET("/play", func(c *gin.Context) {
			code, body := service.LmsPost(touchMacAddress, service.LmsPlay)
			c.String(code, body)
		})
		lms.GET("/pause", func(c *gin.Context) {
			code, body := service.LmsPost(touchMacAddress, service.LmsPause)
			c.String(code, body)
		})

		// decrease and increase volume respectively for the Squeezebox Touch
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

		// decrease and increase volume respectively for the Squeezebox Radio
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

	// run the server on port 8080
	server.Run(":8080")
}

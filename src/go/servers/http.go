package servers

import (
	"frank/src/go/config"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Server *gin.Engine
}

var Http *HttpServer

var server *gin.Engine

func NewHttpServer() {
	server = gin.Default()

	loadBaseRoutes()

	go func() {
		server.Run() // listen and serve on 0.0.0.0:8080
	}()
}

func AddRoute(routeType string, relativePath string, handlers ...gin.HandlerFunc) {
	server.Handle(routeType, relativePath, handlers...)
}

func loadBaseRoutes() {

	server.GET("/devices", func(c *gin.Context) {
		data := c.DefaultQuery("data", "slim")
		if data == "full" {
			devices := []config.Device{}
			for _, d := range config.ParsedConfig.Devices {
				d.Commands = config.GetCommandsByDeviceName(d.Name)
				devices = append(devices, d)
			}
			c.JSON(200, devices)
		} else {
			c.JSON(200, config.ParsedConfig.Devices)
		}

	})

	server.POST("/devices", func(c *gin.Context) {
		var device config.Device
		err := c.BindJSON(&device)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}

		config.AddDevice(device)
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	server.DELETE("/devices/:name", func(c *gin.Context) {
		name := c.Param("name")
		err := config.RemoveDevice(name)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"status": "failed", "error": err})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": "Deleted"})
	})

	server.GET("/commands", func(c *gin.Context) {
		c.JSON(200, config.ParsedConfig.Commands)
	})

	server.POST("/commands", func(c *gin.Context) {
		var command config.Command
		err := c.BindJSON(&command)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}

		config.AddCommand(command)
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	server.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "Frank the bot API",
		})
	})
}

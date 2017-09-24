package servers

import (
	"fmt"
	"net/http"

	"frank/src/go/config"
	"frank/src/go/models"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Server *gin.Engine
}

var Http *HttpServer

var server *gin.Engine

func NewHttpServer() {
	// leave debug for the moment
	// gin.SetMode(gin.ReleaseMode)
	server = gin.Default()

	server.Use(CORSMiddleware())

	loadBaseRoutes()

	go func() {
		server.Run() // listen and serve on 0.0.0.0:8080
	}()
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Max-Age", "86400")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")

		if c.Request.Method == "OPTIONS" {
			fmt.Println("OPTIONS")
			c.AbortWithStatus(200)
		} else {
			c.Next()
		}
	}
}

func AddRoute(routeType string, relativePath string, handlers ...gin.HandlerFunc) {
	server.Handle(routeType, relativePath, handlers...)
}

func loadBaseRoutes() {

	server.GET("/devices", func(c *gin.Context) {
		data := c.DefaultQuery("data", "slim")
		if data == "full" {
			devices := []models.Device{}
			for _, d := range config.ParsedConfig.Devices {
				d.Actions = config.GetActionsByDeviceName(d.Name)
				devices = append(devices, d)
			}
			c.JSON(200, devices)
		} else {
			c.JSON(200, config.ParsedConfig.Devices)
		}

	})

	server.POST("/devices", func(c *gin.Context) {
		var device models.Device
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
		var command models.Command
		err := c.BindJSON(&command)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}

		config.AddCommand(command)
		c.JSON(http.StatusCreated, gin.H{"status": "created"})
	})

	server.GET("/plugins", func(c *gin.Context) {
		c.JSON(200, config.GetAvailablePlugins())
	})

	server.GET("/info", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name": "Frank the bot API",
		})
	})
}

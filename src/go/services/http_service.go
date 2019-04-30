package services

import (
	"fmt"
	"net/http"
	"strconv"

	"frank/src/go/config"
	"frank/src/go/helpers/log"
	"frank/src/go/managers"
	"frank/src/go/models"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Server *gin.Engine
	config *models.HTTP
}

var defaultPort = 8080
var Http HttpServer

var server *gin.Engine

// NewHttpServer istantiates a new HTTP server and stores it into Http.
func NewHTTPServer(config *models.HTTP) {
	if config.Disabled == true || (models.HTTP{}) == *config {
		return
	}

	Http = HttpServer{}
	Http.config = config
	// leave debug for the moment
	gin.SetMode(gin.ReleaseMode)
	server = gin.Default()

	Http.Server = server
	server.Use(CORSMiddleware())

	loadBaseRoutes()
	port := ":" + strconv.Itoa(defaultPort)
	if config.Port != 0 {
		port = ":" + strconv.Itoa(config.Port)
	}
	go func() {
		log.Log.Debug("Starting HTTP/REST server at %s", port)
		server.Run(port) // listen and serve on 0.0.0.0:8080
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
	c := config.GetHTTP()
	if c.Disabled == true {
		log.Log.Warning("Cannot add route %s %s because http server is not running", routeType, relativePath)
		return
	}
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

	server.GET("/crons", func(c *gin.Context) {
		if len(config.ParsedConfig.Crons) == 0 {
			c.JSON(200, "[]")
			return
		}

		c.JSON(200, config.ParsedConfig.Crons)

	})

	server.POST("/ddns", func(c *gin.Context) {
		var ddns models.Ddns
		err := c.BindJSON(&ddns)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}

		config.SetDdns(ddns)
		if config.ParsedConfig.Ddns.Hostname != "" {
			LoadDdns(config.ParsedConfig.Ddns)
			go DdnsManager.SetIp()
		}
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

	server.POST("/action", func(c *gin.Context) {
		var action models.ActionRequest

		err := c.BindJSON(&action)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}

		go managers.ManageAction(action.Name, action.ExtraText)

		c.JSON(http.StatusCreated, gin.H{"status": "Launched"})

	})

	server.POST("/reading", func(c *gin.Context) {
		var reading models.ReadingRequest

		err := c.BindJSON(&reading)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"status": "failed", "error": err})
			return
		}
		fmt.Printf("%v+", reading)
		res, _ := managers.ManageReading(reading.Name)

		c.JSON(http.StatusCreated, res)

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

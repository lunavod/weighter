package main

import (
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tarm/serial"
	"log"
	"net"
	"weighter/scales"
)

func connectToScales() scales.Connection {
	config := GetConfig()
	var conn scales.Connection
	var err error

	if config.Scales.ConnectionType == "IP" {
		conn, _ = net.Dial("tcp", fmt.Sprintf("%s:%d", config.Scales.IP, config.Scales.Port))
	} else {
		c := &serial.Config{Name: config.Scales.COMPort, Baud: 115200}
		conn, err = serial.OpenPort(c)
		if err != nil {
			panic(err)
		}
	}

	return conn
}

func main() {
	config := GetConfig()

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           0,
	}))

	r.GET("/name", func(c *gin.Context) {
		conn := connectToScales()
		defer conn.Close()
		builtRequest := scales.BuildRequest([]byte{scales.Commands["CMD_GET_NAME"]})

		_, err := conn.Write(builtRequest)
		if err != nil {
			log.Fatal(err)
		}

		resp, _ := scales.ReadGetNameResponse(conn)
		c.JSON(200, gin.H{
			"result": resp,
		})
	})

	r.GET("/weight", func(c *gin.Context) {
		conn := connectToScales()
		defer conn.Close()
		builtRequest := scales.BuildRequest([]byte{scales.Commands["CMD_GET_MASSA"]})

		_, err := conn.Write(builtRequest)
		if err != nil {
			log.Fatal(err)
		}

		resp, _ := scales.ReadGetMassaResponse(conn)
		c.JSON(200, gin.H{
			"result": resp,
		})
	})

	addr := fmt.Sprintf("%s:%d", config.Server.IP, config.Server.Port)
	log.Printf("Starting server at %s\n", addr)
	log.Fatal(r.Run(addr))
}

package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/lunavod/weighter/scales"

	"github.com/tarm/serial"
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

func getWeight() uint32 {
	conn := connectToScales()
	defer conn.Close()
	builtRequest := scales.BuildRequest([]byte{scales.Commands["CMD_GET_MASSA"]})

	_, err := conn.Write(builtRequest)
	if err != nil {
		log.Fatal(err)
	}

	resp, _ := scales.ReadGetMassaResponse(conn)
	return resp.Weight
}

var ctx = context.Background()

func main() {
	config := GetConfig()
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.Redis.Addr,
		Password: config.Redis.Password,
		DB:       config.Redis.Db,
	})

	fmt.Println("Connected. Sending weight every 3 seconds")

	x := 1
	for {
		weight := getWeight()
		err := rdb.Set(ctx, "current_weight", int(weight)+x, 0).Err()
		if err != nil {
			panic(err)
		}
		x++
		time.Sleep(3 * time.Second)
	}
}

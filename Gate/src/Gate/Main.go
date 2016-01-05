/**
作者:Kyle Ding
模块：
说明：
创建时间：2015-12-20
**/
package main

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

import (
	. "GxMisc"
)

var r *rand.Rand

func start_server() {
	go clientRouter.Start(config.Host1 + ":" + strconv.Itoa(config.Port1)) //Start client tcp server.
	serverRouter.Start(config.Host2 + ":" + strconv.Itoa(config.Port2))    //Start server tcp server.
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("paramter error")
		fmt.Println("gxcenter <config-file-name>")
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := LoadGateConfig()
	if err != nil {
		Debug("load config fail, err: %s", err)
		return
	}

	InitLogger("gate")

	GeneratePidFile("gate", config.Id)

	r = rand.New(rand.NewSource(time.Now().UnixNano()))
	//
	err = ConnectRedis(config.RedisHost, config.RedisPort, config.RedisDb)
	if err != nil {
		Debug("connect redis fail, err: %s", err)
		return
	}

	gate_run()
	//
	start_server()
}

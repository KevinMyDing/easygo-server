package main

//中心服务器，游戏区服务器，提供每个区服的逻辑

import (
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"time"
)

import (
	. "GxMisc"
	. "GxNet"
)

var r *rand.Rand

func main() {
	//参数检查 参数为配置文件
	if len(os.Args) != 2 {
		fmt.Println("paramter error")
		fmt.Println("gxcenter <config-file-name>")
		return
	}

	//启动所有CPU核心
	runtime.GOMAXPROCS(runtime.NumCPU())

	//读取Center服务器配置
	err := LoadCenterConfig()

	//返回错误
	if err != nil {
		Debug("load config fail, err: %s", err)
		return
	}

	//Initialization logger.
	InitLogger("center")

	//Generate pid file.
	GeneratePidFile("center", config.Id)

	//Return a rand.
	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	//Connect redis database and return error.
	err = ConnectRedis(config.RedisHost, config.RedisPort, config.RedisDb)

	//Return error.
	if err != nil {
		Debug("connect redis fail, err: %s", err)
		return
	}

	//Initialization MySQL and return error.
	err = InitMysql(config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbDb, config.DbCharset)

	//Return error.
	if err != nil {
		Debug("connect mysql fail, err: %s", err)
		return
	}

	//Output debug message.
	Debug("connect mysql ok, host: %s:%d", config.DbHost, config.DbPort)

	//Sync role.
	SyncRole()

	err = ConnectAllGate(uint32(config.Id))
	if err != nil {
		Debug("ConnectAllGate fail, %s", err)
		return
	}
}

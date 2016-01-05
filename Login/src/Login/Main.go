package main

import (
	"errors"
	"fmt"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxStatic"
)

var server *GxTcpServer
var r *rand.Rand

func NewConn(conn *GxTcpConn) {
	Debug("new connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func DisConn(conn *GxTcpConn) {
	Debug("dis connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func NewMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new message, remote: %s", conn.Conn.RemoteAddr().String())
	conn.Send(msg)
	return errors.New("close")
}

func start_server() {
	server = NewGxTcpServer(NewConn, DisConn, NewMessage, true)
	server.RegisterClientCmd(CmdLogin, login)
	server.RegisterClientCmd(CmdRegister, register)
	server.Start(":" + strconv.Itoa(config.Port))
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("paramter error")
		fmt.Println("gxlogin <config-file-name>")
		return
	}
	runtime.GOMAXPROCS(runtime.NumCPU())

	err := LoadGateConfig()
	if err != nil {
		Debug("load config fail, err: %s", err)
		return
	}

	InitLogger("login")

	GeneratePidFile("login", 0)

	r = rand.New(rand.NewSource(time.Now().UnixNano()))

	//
	err = ConnectRedis(config.RedisHost, config.RedisPort, config.RedisDb)
	if err != nil {
		Debug("connect redis fail, err: %s", err)
		return
	}
	Debug("connect redis ok, host: %s:%d", config.RedisHost, config.RedisPort)

	err = InitMysql(config.DbUser, config.DbPwd, config.DbHost, config.DbPort, config.DbDb, config.DbCharset)
	if err != nil {
		Debug("connect mysql fail, err: %s", err)
		return
	}
	Debug("connect mysql ok, host: %s:%d", config.DbHost, config.DbPort)

	err = LoadPlayer()
	if err != nil {
		Debug("load player fail, err: %s", err)
		return
	}
	Debug("load player ok")
	//
	start_server()
	Debug("connect redis fail")
}

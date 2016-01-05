package main

import (
	"fmt"
	"os"
)

import (
	. "GxMisc"
)

type GateConfig struct {
	Id    int
	Host1 string
	Port1 int
	Host2 string
	Port2 int

	DbHost    string
	DbPort    int
	DbUser    string
	DbPwd     string
	DbDb      string
	DbCharset string

	RedisHost string
	RedisPort int
	RedisDb   int64
}

var config *GateConfig

func LoadGateConfig() error {
	err := LoadConfig(os.Args[1])
	if err != nil {
		return err
	}

	config = new(GateConfig)

	config.Id, _ = Config.Get("server").Get("id").Int()

	config.Host1, _ = Config.Get("server").Get("host1").String()
	config.Port1, _ = Config.Get("server").Get("port1").Int()
	config.Host2, _ = Config.Get("server").Get("host2").String()
	config.Port2, _ = Config.Get("server").Get("port2").Int()

	config.DbHost, _ = Config.Get("db").Get("host").String()
	config.DbPort, _ = Config.Get("db").Get("port").Int()
	config.DbUser, _ = Config.Get("db").Get("user").String()
	config.DbPwd, _ = Config.Get("db").Get("pwd").String()
	config.DbDb, _ = Config.Get("db").Get("db").String()
	config.DbCharset, _ = Config.Get("db").Get("charset").String()

	config.RedisHost, _ = Config.Get("redis").Get("host").String()
	config.RedisPort, _ = Config.Get("redis").Get("port").Int()
	config.RedisDb, _ = Config.Get("redis").Get("db").Int64()

	fmt.Println("=================config info=====================")
	fmt.Println("Id       : ", config.Id)
	fmt.Println("Host1    : ", config.Host1)
	fmt.Println("Port1    : ", config.Port1)
	fmt.Println("Host2    : ", config.Host2)
	fmt.Println("Port2    : ", config.Port2)
	fmt.Println("DbHost   : ", config.DbHost)
	fmt.Println("DbPort   : ", config.DbPort)
	fmt.Println("DbUser   : ", config.DbUser)
	fmt.Println("DbPwd    : ", config.DbPwd)
	fmt.Println("DbDb     : ", config.DbDb)
	fmt.Println("DbCharset: ", config.DbCharset)
	fmt.Println("RedisHost: ", config.RedisHost)
	fmt.Println("RedisPort: ", config.RedisPort)
	fmt.Println("RedisDb  : ", config.RedisDb)
	fmt.Println("=================================================")

	return nil
}

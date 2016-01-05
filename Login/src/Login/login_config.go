package main

import (
	"fmt"
	"os"
)

import (
	. "GxMisc"
)

type GateConfig struct {
	Port int

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

	config.Port, _ = Config.Get("port").Int()

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
	fmt.Println("Port     : ", config.Port)
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

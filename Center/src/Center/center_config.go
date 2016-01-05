/**
作者:Kyle Ding
模块：游戏服务器配置模块
说明：
创建时间：2015-12-12
**/
package main

import (
	"fmt"
	"os"
)

import (
	. "GxMisc"
)

type CenterConfig struct {
	Id int //主键

	DbHost    string //主机IP
	DbPort    int    //端口
	DbUser    string //用户名
	DbPwd     string //密码
	DbDb      string //数据库
	DbCharset string //字符集

	RedisHost string //Redis缓存主机IP
	RedisPort int    //Redis缓存端口
	RedisDb   int64  //Redis缓存数据库
}

var config *CenterConfig //CenterConfig结构体

// Load center config file, file type is json
func LoadCenterConfig() error {
	err := LoadConfig(os.Args[1]) //Get command first args and return error.
	if err != nil {
		return err //Return error.
	}

	config = new(CenterConfig) // Create CenterConfig struct.

	// Get config from command args
	config.Id, _ = Config.Get("id").Int()

	config.DbHost, _ = Config.Get("db").Get("host").String()
	config.DbPort, _ = Config.Get("db").Get("port").Int()
	config.DbUser, _ = Config.Get("db").Get("user").String()
	config.DbPwd, _ = Config.Get("db").Get("pwd").String()
	config.DbDb, _ = Config.Get("db").Get("db").String()
	config.DbCharset, _ = Config.Get("db").Get("charset").String()

	config.RedisHost, _ = Config.Get("redis").Get("host").String()
	config.RedisPort, _ = Config.Get("redis").Get("port").Int()
	config.RedisDb, _ = Config.Get("redis").Get("db").Int64()

	// Output to cmd.
	fmt.Println("=================config info=====================")
	fmt.Println("Id       : ", config.Id)
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

package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
	"os"
	"strconv"
	"time"
)

import (
	. "GxMisc"
	// . "GxProto"
	. "GxStatic"
)

var db *sql.DB
var rdClient *redis.Client

func connect_redis() bool {
	rdHost, _ := Config.Get("redis").Get("host").String()
	rdPort, _ := Config.Get("redis").Get("port").Int()
	rdDb, _ := Config.Get("redis").Get("db").Int64()
	rdClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", rdHost, rdPort),
		Password: "",   // no password set
		DB:       rdDb, // use default DB
	})
	return rdClient != nil
}

func connectDb() {
	mysqlHost, _ := Config.Get("db").Get("host").String()
	mysqlPort, _ := Config.Get("db").Get("port").Int()
	mysqlUser, _ := Config.Get("db").Get("user").String()
	mysqlPwd, _ := Config.Get("db").Get("pwd").String()
	mysqlDb, _ := Config.Get("db").Get("db").String()
	mysqlCharset, _ := Config.Get("db").Get("charset").String()

	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", mysqlUser, mysqlPwd, mysqlHost, mysqlPort, mysqlDb, mysqlCharset)

	fmt.Println("connect to", connInfo)

	var err error
	db, err = sql.Open("mysql", connInfo)
	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	// defer db.Close()
}

func systemTable() {
	connectDb()

	str := GenerateCreateSql(&Player{}, "")

	_, err := db.Exec(str)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("new system table ok")
}

func serverTable() {
	if len(os.Args) != 3 {
		fmt.Println("paramter error")
		fmt.Println("GxTool server-table <server-id>")
		return
	}

	connectDb()

	str := GenerateCreateSql(&Role{}, os.Args[2])

	_, err := db.Exec(str)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	str = GenerateIndexSql(&Role{}, os.Args[2], "1")

	_, err = db.Exec(str)
	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}
	fmt.Println("new server table ok")
}

func newServer() {
	if len(os.Args) != 6 {
		fmt.Println("paramter error")
		fmt.Println("GxTool new-server <server-id> <server-name> <server-status> <open-time>")
		return
	}

	if !connect_redis() {
		fmt.Println("connect redis fail")
		return
	}

	id, _ := strconv.Atoi(os.Args[2])
	status, _ := strconv.Atoi(os.Args[4])
	opents := StrToTime(os.Args[5])
	if opents == 0 {
		fmt.Println("open time format error")
		return
	}

	server := &GameServer{
		Id:     uint32(id),
		Name:   os.Args[3],
		Status: uint32(status),
		OpenTs: opents,
	}

	err := SaveGameServer(rdClient, server)
	if err != nil {
		Debug("SaveGameServer 1 error: %s", err)
		return
	}

	fmt.Println("new server ok, info:", server)
}

func test() {
	connectDb()
	for {

		// fmt.Println("===1===", db.Ping())

		rows, err := db.Query("SELECT Id FROM tb_player limit 1")
		defer rows.Close()
		if err != nil {
			fmt.Println("===2===", err)
		}

		for rows.Next() {
			//将行数据保存到record字典
			var id int
			err = rows.Scan(&id)
			fmt.Println("===3===", id)
		}
		time.Sleep(15 * time.Second)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("paramter error")
		fmt.Println("GxTool new-server <server-id> <server-name> <server-status> <open-time>")
		fmt.Println("GxTool system-table")
		fmt.Println("GxTool server-table <server-id>")
		return
	}

	LoadConfig("config.json")

	if os.Args[1] == "system-table" {
		systemTable()
		return
	} else if os.Args[1] == "server-table" {
		serverTable()
		return
	} else if os.Args[1] == "new-server" {
		newServer()
		return
	} else if os.Args[1] == "test" {
		test()
		return
	} else {
		fmt.Println("paramter error")
		fmt.Println("GxTool system-table")
		fmt.Println("GxTool server-table <server-id>")
	}
}

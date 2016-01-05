/**
作者:guangbo
模块：游戏服务器信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"gopkg.in/redis.v3"
	"strconv"
)

import (
	. "GxMisc"
)

var GameServerTableName = "h_game_server"

type GameServer struct {
	Id     uint32 //服务器id
	Name   string //服务器名称
	Status uint32 //服务器状态
	OpenTs int64  //服务器开放时间
}

func SaveGameServer(client *redis.Client, server *GameServer) error {
	buf, err := MsgToBuf(server)
	if err != nil {
		return err
	}

	client.HSet(GameServerTableName, strconv.Itoa(int(server.Id)), string(buf))

	return nil
}

func GetAllGameServer(client *redis.Client, servers *[]*GameServer) error {
	m := client.HGetAllMap(GameServerTableName)
	r, err := m.Result()
	if err != nil {
		return err
	}

	for _, v := range r {
		j, err2 := BufToMsg([]byte(v))
		if err2 != nil {
			return err2
		}
		server := new(GameServer)
		JsonToStruct(j, server)
		*servers = append(*servers, server)
	}
	return nil
}

/**
作者:guangbo
模块：玩家登录信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
	"strconv"
)

import (
	. "GxMisc"
)

type LoginInfo struct {
	PlayerName string `pk:"true"` //帐号名
	GateId     uint32 //当前网关id
	ConnId     uint32 //当前网关连接id
	BeginTs    int64  //登录时间
	EndTs      int64  //登出时间
	ServerId   uint32 //当前连接的游戏服务器
	RoleId     uint32 //当前登录的角色id
}

var GateLoginInfoTableName = "k_gate_login_info:"

func (info *LoginInfo) Set4Redis(client *redis.Client) error {
	SaveToRedis(client, info)
	return nil
}

func (info *LoginInfo) Get4Redis(client *redis.Client) error {
	return LoadFromRedis(client, info)
}

func (info *LoginInfo) Del4Redis(client *redis.Client) error {
	return DelFromRedis(client, info)
}

func (info *LoginInfo) Set4Mysql(db *sql.DB) error {
	return nil
}

func (info *LoginInfo) Get4Mysql(db *sql.DB) error {
	return nil
}

func (info *LoginInfo) Del4Mysql(db *sql.DB) error {
	return nil
}

func SaveGateLoginInfo(client *redis.Client, gateId uint32, connId uint32, palyerName string) {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	client.Set(key, palyerName, 0)
}

func GetGateLoginInfo(client *redis.Client, gateId uint32, connId uint32) string {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	return client.Get(key).Val()
}

func DelGateLoginInfo(client *redis.Client, gateId uint32, connId uint32) {
	key := GateLoginInfoTableName + strconv.Itoa(int(gateId)) + ":" + strconv.Itoa(int(connId))

	client.Del(key)
}

//掉线重连
func DisconnLogin(client *redis.Client, token string, info *LoginInfo) uint16 {
	playerName := CheckToken(client, token)
	if playerName == "" {
		return RetTokenError
	}
	oldInfo := new(LoginInfo)
	oldInfo.PlayerName = playerName
	oldInfo.Get4Redis(client)

	if info.GateId != oldInfo.GateId || info.ConnId != oldInfo.ConnId {
		info.PlayerName = oldInfo.PlayerName
		info.ServerId = oldInfo.ServerId
		info.RoleId = oldInfo.RoleId
		info.BeginTs = oldInfo.BeginTs
		DelGateLoginInfo(client, oldInfo.GateId, oldInfo.ConnId)
		SaveGateLoginInfo(client, info.GateId, info.ConnId, playerName)
		info.Set4Redis(client)
	}
	return RetSucc
}

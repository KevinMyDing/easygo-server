/**
作者:guangbo
模块：保存玩家最近登录服务器信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

import (
	. "GxMisc"
)

type LastLoginInfo struct {
	PlayerName   string `pk:"true"` //帐号名
	GameServerId uint32 //服务器id
	Ts           int64  `type:"time"` //登录时间
}

func (lli *LastLoginInfo) Set4Redis(client *redis.Client) error {
	SaveToRedis(client, lli)
	return nil
}

func (lli *LastLoginInfo) Get4Redis(client *redis.Client) error {
	return LoadFromRedis(client, lli)
}

func (lli *LastLoginInfo) Del4Redis(client *redis.Client) error {
	return DelFromRedis(client, lli)
}

func (lli *LastLoginInfo) Set4Mysql(db *sql.DB) error {
	return nil
}

func (lli *LastLoginInfo) Get4Mysql(db *sql.DB) error {
	return nil
}

func (lli *LastLoginInfo) Del4Mysql(db *sql.DB) error {
	return nil
}

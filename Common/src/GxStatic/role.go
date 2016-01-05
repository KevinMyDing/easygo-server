/**
作者:guangbo
模块：游戏角色信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
	"strconv"
	"sync"
)

import (
	. "GxMisc"
)

var RoleListTableName = "l_role_list:"
var RoleNameListTableName = "s_role_name:"
var RoleIdTableName = "k_role_id"
var RoleCreateList = "s_role_create:"

var roleIdMutex *sync.Mutex

type Role struct {
	Id         uint32 `pk:"true"`          //角色id
	PlayerName string `len:"32" index:"1"` //所属帐号名
	Name       string `len:"32"`           //角色名
	VocationId uint32 //职业
	Vip        uint32 //vip等级
	Expr       uint64 //累计经验
	Money      uint64 //金币
	Gold       uint64 //元宝
	Power      uint64 //体力
	PowerTs    uint64 `type:"time"` //体力恢复时间
}

func init() {
	roleIdMutex = new(sync.Mutex)
}

func NewRoleID(client *redis.Client) uint32 {
	roleIdMutex.Lock()
	defer roleIdMutex.Unlock()

	if !client.Exists(RoleIdTableName).Val() {
		client.Set(RoleIdTableName, "10000000", 0)

	}
	return uint32(client.Incr(RoleIdTableName).Val())
}

func (role *Role) Set4Redis(client *redis.Client) error {
	SaveToRedis(client, role)
	return nil
}

func (role *Role) Get4Redis(client *redis.Client) error {
	return LoadFromRedis(client, role)
}

func (role *Role) Del4Redis(client *redis.Client) error {
	return DelFromRedis(client, role)
}

func (role *Role) Set4Mysql(db *sql.DB, serverId uint32) error {
	r := new(Role)
	r.Id = role.Id
	var str string
	if r.Get4Mysql(db, serverId) == nil {
		str = GenerateUpdateSql(role, strconv.FormatUint(uint64(serverId), 10))
	} else {
		str = GenerateInsertSql(role, strconv.FormatUint(uint64(serverId), 10))
	}

	_, err := db.Exec(str)
	if err != nil {
		fmt.Println("set error", err)
		return err
	}
	return nil
}

func (role *Role) Get4Mysql(db *sql.DB, serverId uint32) error {
	str := GenerateSelectSql(role, strconv.FormatUint(uint64(serverId), 10))

	rows, err := db.Query(str)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		err = rows.Scan(&role.PlayerName, &role.Name, &role.VocationId, &role.Vip, &role.Expr, //
			&role.Money, &role.Gold, &role.Power, &role.PowerTs)
		return nil
	}
	return errors.New("null")
}

func (role *Role) Del4Mysql(db *sql.DB, serverId uint32) error {
	str := GenerateDeleteSql(role, strconv.FormatUint(uint64(serverId), 10))
	_, err := db.Exec(str)
	if err != nil {
		return err
	}
	return nil
}

//玩家角色id列表.玩家登录服务器时候,需要拉取所有角色信息
func SaveRoleId4RoleList(client *redis.Client, playerName string, gameServerId uint32, roleid uint32) {
	serverId := strconv.Itoa(int(gameServerId))
	id := strconv.Itoa(int(roleid))
	client.LPush(RoleListTableName+playerName+":"+serverId, id)
}

func GetRoleIds4RoleList(client *redis.Client, playerName string, gameServerId uint32) []string {
	return client.LRange(RoleListTableName+playerName+":"+strconv.Itoa(int(gameServerId)), 0, -1).Val()
}

//创建角色时候检查角色名是否冲突
func SaveRoleName(client *redis.Client, gameServerId uint32, name string) bool {
	return client.SAdd(RoleNameListTableName+strconv.Itoa(int(gameServerId)), name).Val() == 1
}

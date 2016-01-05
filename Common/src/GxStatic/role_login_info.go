/**
作者:guangbo
模块：保存角色最近登录信息模块
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

var RoleLoginTableName = "h_role_Login:"

type RoleLoginInfo struct {
	RoleId uint32 //角色id
	Del    uint32 //是否已经清除了缓存
	Ts     int64  `type:"time"` //登录时间
}

func SaveRoleLogin(client *redis.Client, serverId uint32, info *RoleLoginInfo) error {
	buf, err := MsgToBuf(info)
	if err != nil {
		return err
	}

	client.HSet(RoleLoginTableName+strconv.Itoa(int(serverId)), strconv.Itoa(int(info.RoleId)), string(buf))

	return nil
}

func GetAllRoleLogin(client *redis.Client, serverId uint32, infos *[]*RoleLoginInfo) error {
	m := client.HGetAllMap(RoleLoginTableName + strconv.Itoa(int(serverId)))
	r, err := m.Result()
	if err != nil {
		return err
	}

	for _, v := range r {
		j, err2 := BufToMsg([]byte(v))
		if err2 != nil {
			return err2
		}
		info := new(RoleLoginInfo)
		JsonToStruct(j, info)
		*infos = append(*infos, info)
	}
	return nil
}

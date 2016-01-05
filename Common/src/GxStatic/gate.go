/**
作者:guangbo
模块：网关信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"strconv"

	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

import (
	. "GxMisc"
)

var GateInfoTableName string = "h_gate_info"

type GateInfo struct {
	Id    uint32 //网关id
	Host1 string //网关外网ip
	Port1 uint32 //网关外网端口
	Host2 string //网关内网ip
	Port2 uint32 //网关内网端口
	Count uint32 //当前连接数
	Ts    int64  //信息更新时间
}

func SaveGate(client *redis.Client, gate *GateInfo) error {
	buf, err := MsgToBuf(gate) //Convert to json
	if err != nil {
		return err
	}

	client.HSet(GateInfoTableName, strconv.Itoa(int(gate.Id)), string(buf)) //Save data to redis.

	return nil
}

func GetAllGate(client *redis.Client, gates *[]*GateInfo) error {
	m := client.HGetAllMap(GateInfoTableName)
	r, err := m.Result()
	if err != nil {
		return err
	}
	for _, v := range r {
		j, err2 := BufToMsg([]byte(v))
		if err2 != nil {
			return err2
		}
		gate := new(GateInfo)
		JsonToStruct(j, gate)
		*gates = append(*gates, gate)
	}
	return nil
}

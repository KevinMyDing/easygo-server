/**
作者:guangbo
模块：角色信息管理模块
说明：
创建时间：2015-11-2
**/
package main

import (
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v3"
)

import (
	. "GxMisc"
	. "GxProto"
	. "GxStatic"
)

func FillSelectRoleRsp(client *redis.Client, roleId uint32, serverId uint32, rsp *SelectRoleRsp) (string, error) {
	rsp.Role = new(RoleCommonInfo)
	playername, err := FillRoleCommonInfo(client, roleId, serverId, rsp.GetRole())
	if err != nil {
		return "", err
	}
	return playername, nil
}

func FillCreateRoleRsp(client *redis.Client, roleId uint32, serverId uint32, rsp *CreateRoleRsp) error {
	rsp.Role = new(RoleCommonInfo)
	_, err := FillRoleCommonInfo(client, roleId, serverId, rsp.GetRole())
	if err != nil {
		return err
	}
	return nil
}

func FillRoleCommonInfo(client *redis.Client, roleId uint32, serverId uint32, role *RoleCommonInfo) (string, error) {
	r := new(Role)
	r.Id = roleId
	err := r.Get4Redis(client)
	if err != nil {
		err = r.Get4Mysql(Db, serverId)
		if err != nil {
			return "", err
		}
		r.Set4Redis(client)
	}

	role.Id = proto.Uint32(r.Id)
	role.Name = proto.String(r.Name)
	role.VocationId = proto.Uint32(r.VocationId)
	role.Vip = proto.Uint32(r.Vip)
	role.Expr = proto.Uint64(r.Expr)
	role.Money = proto.Uint64(r.Money)
	role.Gold = proto.Uint64(r.Gold)
	role.Power = proto.Uint64(r.Power)
	role.Powerts = proto.Uint64(r.PowerTs)
	return r.PlayerName, nil
}

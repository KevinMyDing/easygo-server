/**
作者:guangbo
模块：角色初始化管理模块
说明：
创建时间：2015-11-2
**/
package main

import (
	"gopkg.in/redis.v3"
	"strconv"
	"time"
)

import (
	. "GxStatic"
)

func CreateRole(client *redis.Client, playerName string, roleName string, vocationId uint32) uint32 {
	id := initRoleComm(client, playerName, roleName, vocationId)

	return id
}

func initRoleComm(client *redis.Client, playerName string, roleName string, vocationId uint32) uint32 {
	role := &Role{
		Id:         NewRoleID(client),
		PlayerName: playerName,
		Name:       roleName,
		VocationId: vocationId,
		Vip:        0,
		Expr:       0,
		Money:      0,
		Gold:       0,
		Power:      20,
		PowerTs:    uint64(time.Now().Unix()),
	}
	role.Set4Redis(client)

	client.SAdd(RoleCreateList+strconv.Itoa(config.Id), strconv.Itoa(int(role.Id)))
	return role.Id
}

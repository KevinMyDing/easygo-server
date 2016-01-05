/**
作者:guangbo
模块：角色消息处理模块
说明：
创建时间：2015-11-2
**/
package main

import (
	"github.com/golang/protobuf/proto"
	"strconv"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

func init() {
	RegisterMessageCallback(CmdHeartbeat, ClientHeartbeatCallback)
	RegisterMessageCallback(CmdGetRoleList, GetRoleListCallback)
	RegisterMessageCallback(CmdSelectRole, SelectRoleCallback)
	RegisterMessageCallback(CmdCreateRole, CreateRoleCallback)
}

func ClientHeartbeatCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {

}

func GetRoleListCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req GetRoleListReq
	var rsp GetRoleListRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}
	if req.Info == nil || req.GetInfo().Token == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	playerName := CheckToken(rdClient, req.GetInfo().GetToken())
	if playerName == "" {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetTokenError, nil)
		return
	}

	ids := GetRoleIds4RoleList(rdClient, playerName, uint32(config.Id))
	for i := 0; i < len(ids); i++ {
		id, _ := strconv.Atoi(ids[i])

		role := new(Role)
		role.Id = uint32(id)
		err = role.Get4Redis(rdClient)
		if err != nil {
			Debug("role %d is not existst", id)
			continue
		}
		rsp.Roles = append(rsp.Roles, &RoleCommonInfo{
			Id:         proto.Uint32(role.Id),
			Name:       proto.String(role.Name),
			Expr:       proto.Uint64(role.Expr),
			VocationId: proto.Uint32(role.VocationId),
		})
	}

	info.PlayerName = playerName
	info.BeginTs = time.Now().Unix()
	info.ServerId = uint32(config.Id)
	SaveGateLoginInfo(rdClient, info.GateId, info.ConnId, playerName)
	info.Set4Redis(rdClient)

	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &rsp)
}

func SelectRoleCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req SelectRoleReq
	var rsp SelectRoleRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}

	if req.RoleId == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	if req.Info != nil && req.GetInfo().Token != nil {
		//重新重连
		ret := DisconnLogin(rdClient, req.GetInfo().GetToken(), info)
		if ret != RetSucc {
			SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), ret, nil)
			return
		}
	}

	playerName, err1 := FillSelectRoleRsp(rdClient, req.GetRoleId(), uint32(config.Id), &rsp)
	if err1 != nil {
		Debug("role %d is not existst", req.GetRoleId())
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleNotExists, nil)
		return
	}

	//登陆角色不属于当前账号
	if info.PlayerName != playerName {
		Debug("role %d is not player: %s 's role", req.GetRoleId(), info.PlayerName)
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}

	info.RoleId = req.GetRoleId()
	info.Set4Redis(rdClient)

	SaveRoleLogin(rdClient, uint32(config.Id), &RoleLoginInfo{
		RoleId: info.RoleId,
		Del:    0,
		Ts:     time.Now().Unix(),
	})

	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &rsp)

}

func CreateRoleCallback(conn *GxTcpConn, info *LoginInfo, msg *GxMessage) {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req CreateRoleReq
	var rsp CreateRoleRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return
	}

	if req.Name == nil || req.VocationId == nil {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetMsgFormatError, nil)
		return
	}

	if req.Info != nil && req.GetInfo().Token != nil {
		//重新重连
		ret := DisconnLogin(rdClient, req.GetInfo().GetToken(), info)
		if ret != RetSucc {
			SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), ret, nil)
			return
		}
	}

	ids := GetRoleIds4RoleList(rdClient, info.PlayerName, info.ServerId)
	if len(ids) > 0 {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleExists, nil)
		return
	}

	//检查角色名冲突
	if !SaveRoleName(rdClient, info.ServerId, req.GetName()) {
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleNameConflict, nil)
		return
	}

	//init role info
	info.RoleId = CreateRole(rdClient, info.PlayerName, req.GetName(), req.GetVocationId())
	info.Set4Redis(rdClient)

	SaveRoleId4RoleList(rdClient, info.PlayerName, uint32(config.Id), info.RoleId)
	//
	err = FillCreateRoleRsp(rdClient, info.RoleId, uint32(config.Id), &rsp)
	if err != nil {
		Debug("role %d is not existst", info.RoleId)
		SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetRoleNotExists, nil)
		return
	}
	//
	SendPbMessage(conn, false, msg.GetId(), msg.GetCmd(), msg.GetSeq(), RetSucc, &rsp)
}

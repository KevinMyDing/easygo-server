package main

import (
	"errors"
	"github.com/golang/protobuf/proto"
	"gopkg.in/redis.v3"
	"time"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

func fillLoginRsp(rdClient *redis.Client, player *Player, rsp *LoginServerRsp) {
	lli := new(LastLoginInfo)
	lli.PlayerName = player.Username
	lastLogin := lli.Get4Redis(rdClient) != nil

	rsp.Info = &LoginRspInfo{
		Token: proto.String(player.SaveToken(rdClient)),
	}

	var gates []*GateInfo
	GetAllGate(rdClient, &gates)
	var gate *GateInfo = nil
	for i := 0; i < len(gates); i++ {
		if (time.Now().Unix() - gates[i].Ts) > 10 {
			continue
		}
		if gate == nil {
			gate = gates[i]
		} else {
			if gate.Count > gates[i].Count {
				gate = gates[i]
			}
		}
	}
	if gate != nil {
		rsp.GetInfo().Host = proto.String(gate.Host1)
		rsp.GetInfo().Port = proto.Uint32(gate.Port1)
	}

	var servers []*GameServer
	GetAllGameServer(rdClient, &servers)
	for i := 0; i < len(servers); i++ {
		Debug("server-id: %d, name: %s", servers[i].Id, servers[i].Name)
		var lastts int64 = 0
		if lastLogin && lli.GameServerId == servers[i].Id {
			lastts = lli.Ts
		}
		rsp.GetInfo().Srvs = append(rsp.GetInfo().Srvs, &GameSrvInfo{
			Index:  proto.Uint32(servers[i].Id),
			Name:   proto.String(servers[i].Name),
			Status: proto.Uint32(servers[i].Status),
			Lastts: proto.Uint32(uint32(lastts)),
		})
	}

}

func login(conn *GxTcpConn, msg *GxMessage) error {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req LoginServerReq
	var rsp LoginServerRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		Debug("UnpackagePbmsg error")
		return errors.New("close")
	}

	player := FindPlayer(req.GetRaw().GetUsername())
	if player == nil {
		Debug("user is not exists, username: %s", req.GetRaw().GetUsername())
		SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetUserNotExists, nil)
		return errors.New("close")
	}

	if !VerifyPassword(player, req.GetRaw().GetPwd()) {
		SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetPwdError, nil)
		return errors.New("close")
	}

	Debug("old user: %s login from %s", req.GetRaw().GetUsername(), conn.Conn.RemoteAddr().String())
	fillLoginRsp(rdClient, player, &rsp)

	SendPbMessage(conn, false, 0, CmdLogin, msg.GetSeq(), RetSucc, &rsp)
	return errors.New("close")
}

func register(conn *GxTcpConn, msg *GxMessage) error {
	rdClient := PopRedisClient()
	defer PushRedisClient(rdClient)

	var req LoginServerReq
	var rsp LoginServerRsp
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		Debug("UnpackagePbmsg error")
		return errors.New("close")
	}

	player := FindPlayer(req.GetRaw().GetUsername())
	if player != nil {
		Debug("user has been exists, username: %s", req.GetRaw().GetUsername())
		SendPbMessage(conn, false, 0, CmdRegister, msg.GetSeq(), RetUserExists, nil)
		return errors.New("close")
	}

	player = NewPlayer(rdClient, req.Raw.GetUsername(), req.GetRaw().GetPwd(), uint32(req.GetPt()))
	err = AddPlayer(player)
	if err != nil {
		Debug("user has been exists, username: %s", req.GetRaw().GetUsername())
		SendPbMessage(conn, false, 0, CmdRegister, msg.GetSeq(), RetUserExists, nil)
		return errors.New("close")
	}

	fillLoginRsp(rdClient, player, &rsp)

	Debug("new user: %s login from %s", req.Raw.GetUsername(), conn.Conn.RemoteAddr().String())
	SendPbMessage(conn, false, 0, CmdRegister, msg.GetSeq(), RetSucc, &rsp)

	return errors.New("close")
}

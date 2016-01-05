package main

import (
	"github.com/golang/protobuf/proto"
)
import (
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

var LoginServerAddr = "127.0.0.1:9000"
var GateServerAddr = "127.0.0.1:13000"

var username = "guang"
var pwd = "guang123"
var token = "3174a3d9204acb90576712e1703004d6"

var serverId = 1

var counter *Counter

func test_login(conn *GxTcpConn) {
	var req LoginServerReq
	var rsp LoginServerRsp

	req.Raw = &PlayerRaw{
		Username: proto.String(username),
		Pwd:      proto.String(pwd),
	}
	SendPbMessage(conn, false, 0, CmdLogin, uint16(counter.Genarate()), 0, &req)
	//
	msg, err := conn.Recv()
	if err != nil {
		Debug("recv error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s", rsp.String())
	err = msg.UnpackagePbmsg(&rsp)
	if err != nil {
		Debug("unpackage error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s\n\t%s", msg.String(), rsp.String())
	Debug("servername: %s", rsp.GetInfo().GetSrvs()[0].GetName())
}

func test_register(conn *GxTcpConn) {
	var req LoginServerReq
	var rsp LoginServerRsp

	req.Raw = &PlayerRaw{
		Username: proto.String(username),
		Pwd:      proto.String(pwd),
	}

	pt := PlatformType_GC_PT_91_ASSISTANT
	req.Pt = &pt
	SendPbMessage(conn, false, 0, CmdRegister, uint16(counter.Genarate()), 0, &req)
	//
	msg, err := conn.Recv()
	if err != nil {
		Debug("recv error, %s", err)
		return
	}
	err = msg.UnpackagePbmsg(&rsp)
	if err != nil {
		Debug("unpackage error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s\n\t%s", msg.String(), rsp.String())
}

func test_get_role_list(conn *GxTcpConn) {
	var req GetRoleListReq
	var rsp GetRoleListRsp
	req.Info = &RequestInfo{
		Token: proto.String(token),
	}

	SendPbMessage(conn, false, uint32(serverId), CmdGetRoleList, uint16(counter.Genarate()), 0, &req)
	//
	msg, err := conn.Recv()
	if err != nil {
		Debug("recv error, %s", err)
		return
	}
	err = msg.UnpackagePbmsg(&rsp)
	if err != nil {
		Debug("unpackage error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s\n\t%s", msg.String(), rsp.String())
}

func test_select_role(conn *GxTcpConn) {
	var req SelectRoleReq
	var rsp SelectRoleRsp
	req.Info = &RequestInfo{
		Token: proto.String(token),
	}
	req.RoleId = proto.Uint32(10000004)

	SendPbMessage(conn, false, uint32(serverId), CmdSelectRole, uint16(counter.Genarate()), 0, &req)
	//
	msg, err := conn.Recv()
	if err != nil {
		Debug("recv error, %s", err)
		return
	}
	err = msg.UnpackagePbmsg(&rsp)
	if err != nil {
		Debug("unpackage error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s\n\t%s", msg.String(), rsp.String())
}

func test_create_role(conn *GxTcpConn) {
	var req CreateRoleReq
	var rsp CreateRoleRsp

	req.Info = &RequestInfo{
		Token: proto.String(token),
	}
	req.Name = proto.String(username)
	req.VocationId = proto.Uint32(1)

	SendPbMessage(conn, false, uint32(serverId), CmdCreateRole, uint16(counter.Genarate()), 0, &req)
	//
	msg, err := conn.Recv()
	if err != nil {
		Debug("recv error, %s", err)
		return
	}
	err = msg.UnpackagePbmsg(&rsp)
	if err != nil {
		Debug("unpackage error, %s", err)
		return
	}
	Debug("recv buff msg, info: %s\n\t%s", msg.String(), rsp.String())
}

func main() {
	InitLogger("client")

	counter = NewCounter()

	conn := NewTcpConn()
	err := conn.Connect(GateServerAddr)
	if err != nil {
		Debug("new connnect, remote: %s", err)
		return
	}
	defer conn.Conn.Close()

	// test_login(conn)
	// test_register(conn)
	test_get_role_list(conn)
	test_select_role(conn)
	// test_create_role(conn)
	return
}

/**
作者:Kyle Ding
模块：网关内网连接管理模块
说明：
创建时间：2015-12-20
**/
package main

import (
	"sync"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
	. "GxProto"
	. "GxStatic"
)

type ServersInfo map[string]*GxTcpConn

var serverRouter *GxTcpServer
var mutex *sync.Mutex
var CmdServers map[uint16]ServersInfo //保存所有注册cmd的服务器
var IdServers map[uint32]*GxTcpConn

func init() {
	mutex = new(sync.Mutex)
	CmdServers = make(map[uint16]ServersInfo)
	IdServers = make(map[uint32]*GxTcpConn)

	serverRouter = NewGxTcpServer(ServerNewConn, ServerDisConn, ServerRawMessage, false)
	serverRouter.RegisterClientCmd(CmdServerConnectGate, ServerConnectGateCallback)
}

func ServerNewConn(conn *GxTcpConn) {
	Debug("new server connnect, remote: %s", conn.Conn.RemoteAddr().String())
}

func ServerDisConn(conn *GxTcpConn) {
	Debug("dis server connnect, remote: %s", conn.Conn.RemoteAddr().String())

	mutex.Lock()
	for i := 0; i < len(conn.Data); i++ {
		if conn.Id == 0 {
			delete(CmdServers[uint16(conn.Data[i].(uint32))], conn.Conn.RemoteAddr().String())
		} else {
			delete(IdServers, conn.Id)
		}

	}
	mutex.Unlock()
}

func ServerRawMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new server message, remote: %s %s", conn.Conn.RemoteAddr().String(), msg.String())

	client := clientRouter.FindConnById(msg.GetId())
	if client == nil {
		Debug("msg cannot find target, remote: %s, msg: %s", conn.Conn.RemoteAddr().String(),
			msg.String())
		return nil
	}

	client.Send(msg)
	if msg.GetMask(MessageMaskDisconn) {
		client.Conn.Close()
	}
	return nil
}

func ServerConnectGateCallback(conn *GxTcpConn, msg *GxMessage) error {
	//register server
	var req ServerConnectGateReq
	err := msg.UnpackagePbmsg(&req)
	if err != nil {
		SendPbMessage(conn, false, 0, msg.GetCmd(), msg.GetSeq(), RetFail, nil)
		return err
	}

	mutex.Lock()
	if req.GetId() == 0 {
		//公用服务器
		conn.Data = make([]interface{}, len(req.Cmds))

		for i := 0; i < len(req.Cmds); i++ {
			//保存自己处理的消息cmd
			conn.Data[i] = req.Cmds[i]

			cmd := uint16(req.Cmds[i])
			cmdinfo, ok := CmdServers[cmd]
			if !ok {
				CmdServers[cmd] = make(ServersInfo)
				cmdinfo, _ = CmdServers[cmd]
			}
			cmdinfo[conn.Conn.RemoteAddr().String()] = conn
		}
	} else {
		//游戏区服务器
		IdServers[req.GetId()] = conn
	}
	mutex.Unlock()

	SendPbMessage(conn, false, 0, msg.GetCmd(), msg.GetSeq(), RetSucc, nil)
	return nil
}

func GetServerByCmd(cmd uint16) (string, *GxTcpConn) {
	mutex.Lock()
	defer mutex.Unlock()

	info, ok := CmdServers[cmd]
	if ok {
		count := len(info)
		var v []string
		for _, s := range info {
			v = append(v, s.Conn.RemoteAddr().String())
		}
		remote := v[r.Intn(count)]
		conn, ok1 := info[remote]
		if ok1 {
			return remote, conn
		}
	}
	return "", nil
}

func GetServerById(id uint32) *GxTcpConn {
	mutex.Lock()
	defer mutex.Unlock()

	conn, ok := IdServers[id]
	if ok {
		return conn
	}
	return nil
}

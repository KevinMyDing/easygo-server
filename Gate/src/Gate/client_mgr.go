package main

import (
	"errors"
)

import (
	. "GxMessage"
	. "GxMisc"
	. "GxNet"
)

var clientRouter *GxTcpServer

func init() {
	clientRouter = NewGxTcpServer(ClientNewConn, ClientDisConn, ClientRawMessage, true)
}

func ClientNewConn(conn *GxTcpConn) {
	Debug("new client connnect, remote: %s", conn.Conn.RemoteAddr().String())
	addClient()
}

func ClientDisConn(conn *GxTcpConn) {
	Debug("dis client connnect, remote: %s", conn.Conn.RemoteAddr().String())
	subClient()
}

func ClientRawMessage(conn *GxTcpConn, msg *GxMessage) error {
	Debug("new client message, remote: %s %s", conn.Conn.RemoteAddr().String(), msg.String())
	//保存目的服务器id
	serverId := msg.GetId()

	msg.SetId(conn.Id)

	_, server := GetServerByCmd(msg.GetCmd())
	if server != nil {
		server.Send(msg)
		return nil
	}

	server = GetServerById(serverId)
	if server == nil {
		Debug("Can not find a server, info: %s", msg.String())
		return errors.New("Can not find a server")
	}
	server.Send(msg)
	return nil
}

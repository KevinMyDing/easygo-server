/**
作者:guangbo
模块：消息命令字模块
说明：
创建时间：2015-10-30
**/
package GxStatic

const (
	//login
	CmdRegister = 1000
	CmdLogin    = 1001

	//gate
	CmdHeartbeat   = 2000
	CmdGetRoleList = 2001
	CmdCreateRole  = 2002
	CmdSelectRole  = 2003

	//server
	CmdServerConnectGate = 3000
)

type ServerStatus int

const (
	ServerStatusHot ServerStatus = iota
	ServerStatusNew              //1
	ServerStatusMaintain
)

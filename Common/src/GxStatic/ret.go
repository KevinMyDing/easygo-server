/**
作者:guangbo
模块：返回值
说明：
创建时间：2015-10-30
**/
package GxStatic

const (
	RetSucc uint16 = iota
	RetFail        //1
	RetMessageNotSupport
	RetMsgFormatError  //2
	RetPwdError        //3
	RetUserNotExists   //4
	RetUserExists      //5
	RetServerNotExists //6
	//
	RetTokenError       //7
	RetRoleNotExists    //8
	RetRoleExists       //9
	RetRoleNameConflict //10
	RetNotLogin         //11
)

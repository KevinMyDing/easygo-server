/**
作者:guangbo
模块：mysql连接池
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Db *sql.DB

func InitMysql(user string, pwd string, host string, port int, dbs string, charset string) error {
	var err error
	connInfo := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s", user, pwd, host, port, dbs, charset)
	Db, err = sql.Open("mysql", connInfo)
	if err != nil {
		return err
	}
	Db.SetMaxOpenConns(128)
	Db.SetMaxIdleConns(64)
	Db.Ping()

	return nil
}

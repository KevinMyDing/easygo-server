/**
作者:guangbo
模块：
说明：
创建时间：2015-10-30
**/
package GxMisc

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
)

type GxModule interface {
	Set4Redis(client *redis.Client) error
	Get4Redis(client *redis.Client) error
	Del4Redis(client *redis.Client) error
	Set4Mysql(db *sql.DB) error
	Get4Mysql(db *sql.DB) error
	Del4Mysql(db *sql.DB) error
}

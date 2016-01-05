/**
作者:guangbo
模块：帐号信息模块
说明：
创建时间：2015-10-30
**/
package GxStatic

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/redis.v3"
	"io"
	// "strconv"
	"errors"
	"sync"
	"time"
)

import (
	. "GxMisc"
)

var salt1 = "f4g*h(j"
var salt2 = "1^2&4*(d)"

var PlayerIdTableName = "k_player_id"
var PlayerTokenTableName = "k_player_token:"
var PlayerNameListTableName = "s_player_name:"

var idMutex *sync.Mutex

type Player struct {
	Id       uint32 //帐号id
	Username string `pk:"true"` //帐号名
	Password string //帐号密码
	CreateTs uint64 `type:"time"` //创建帐号时间
	Platform uint32 //所属平台
}

func init() {
	idMutex = new(sync.Mutex)
}

func newPalayerID(client *redis.Client) uint32 {
	idMutex.Lock()
	defer idMutex.Unlock()

	if !client.Exists(PlayerIdTableName).Val() {
		client.Set(PlayerIdTableName, "100000", 0)

	}
	return uint32(client.Incr(PlayerIdTableName).Val())
}

func NewPlayer(client *redis.Client, username string, password string, platform uint32) *Player {
	player := &Player{
		Id:       newPalayerID(client),
		Username: username,
		Password: generatePassward(username, password),
		CreateTs: uint64(time.Now().Unix()),
		Platform: platform,
	}
	return player
}

func (player *Player) Set4Redis(client *redis.Client) error {
	SaveToRedis(client, player)
	return nil
}

func (player *Player) Get4Redis(client *redis.Client) error {
	return LoadFromRedis(client, player)
}

func (player *Player) Del4Redis(client *redis.Client) error {
	return DelFromRedis(client, player)
}

func (player *Player) Set4Mysql(db *sql.DB) error {
	r := new(Player)
	r.Id = player.Id
	var str string
	if r.Get4Mysql(db) == nil {
		str = GenerateUpdateSql(player, "")
	} else {
		str = GenerateInsertSql(player, "")
	}

	_, err := db.Exec(str)
	if err != nil {
		fmt.Println("set error", err)
		return err
	}
	return nil
}

func (player *Player) Get4Mysql(db *sql.DB) error {
	str := GenerateSelectSql(player, "")

	rows, err := db.Query(str)
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer rows.Close()

	if rows.Next() {
		rows.Scan(&player.Id, &player.Password, &player.CreateTs, &player.Platform)

		fmt.Println("get one")
		return nil
	}
	fmt.Println("get null")
	return errors.New("null")
}

func (player *Player) Del4Mysql(db *sql.DB) error {
	str := GenerateDeleteSql(player, "")
	_, err := db.Exec(str)
	if err != nil {
		return err
	}
	return nil
}

//生成密码函数
func generatePassward(username string, password string) string {
	h := md5.New()

	io.WriteString(h, salt1)
	io.WriteString(h, username)
	io.WriteString(h, salt2)
	io.WriteString(h, password)

	return fmt.Sprintf("%x", h.Sum(nil))
}

func VerifyPassword(player *Player, password string) bool {
	return player.Password == generatePassward(player.Username, password)
}

//生成token函数
func (player *Player) SaveToken(client *redis.Client) string {
	h := md5.New()
	io.WriteString(h, player.Username)
	io.WriteString(h, fmt.Sprintf("%ld", time.Now().Unix()))
	io.WriteString(h, "2@#RR#R@")
	token := fmt.Sprintf("%x", h.Sum(nil))

	client.Set(PlayerTokenTableName+token, player.Username, time.Hour*1)
	return token
}

func CheckToken(client *redis.Client, token string) string {
	key := PlayerTokenTableName + token
	if !client.Exists(key).Val() {
		return ""
	}
	return client.Get(key).Val()
}

//创建角色时候检查角色名是否冲突
func CheckPlayerNameConflict(client *redis.Client, playerName string) bool {
	return client.SIsMember(PlayerNameListTableName, playerName).Val()
}

func SavePlayerName(client *redis.Client, playerName string) {
	client.SAdd(PlayerNameListTableName, playerName)
}

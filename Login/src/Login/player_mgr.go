package main

import (
	"container/list"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"sync"
	"time"
)

import (
	. "GxMisc"
	. "GxStatic"
)

type Players map[string]*Player

var players Players
var playersMutex *sync.Mutex

var sqlList *list.List
var sqlMutex *sync.Mutex

func init() {
	sqlList = list.New()
	sqlMutex = new(sync.Mutex)

	players = make(Players)
	playersMutex = new(sync.Mutex)
}

func PushPlayerSql(sql string) {
	sqlMutex.Lock()
	defer sqlMutex.Unlock()

	sqlList.PushBack(sql)
}

func PopPlayerSql() string {
	sqlMutex.Lock()
	defer sqlMutex.Unlock()

	if sqlList.Len() == 0 {
		return ""
	}

	sql := sqlList.Front().Value.(string)
	sqlList.Remove(sqlList.Front())
	return sql
}

func LoadPlayer() error {
	str := GenerateSelectAllSql(&Player{}, "")
	Debug("load player, sql: %d", str)
	rows, err := Db.Query(str)
	defer rows.Close()
	if err != nil {
		return err
	}

	n := 0
	for rows.Next() {
		player := new(Player)
		err = rows.Scan(&player.Id, &player.Username, &player.Password, &player.CreateTs, &player.Platform)
		players[player.Username] = player
		n++
	}
	Debug("load player, count: %d", n)

	go func() {
		for {
			sql := PopPlayerSql()
			if sql == "" {
				time.Sleep(time.Millisecond * 10)
				continue
			}
			_, err := Db.Exec(sql)
			if err != nil {
				fmt.Println(err, ",sql:", sql)
			}
		}

	}()

	return nil
}

func FindPlayer(name string) *Player {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	player, ok := players[name]
	if ok {
		return player
	}

	return nil
}

func AddPlayer(player *Player) error {
	playersMutex.Lock()
	defer playersMutex.Unlock()

	_, ok := players[player.Username]
	if ok {
		return errors.New("player name is exists")
	}

	players[player.Username] = player
	PushPlayerSql(GenerateInsertSql(player, ""))
	return nil
}

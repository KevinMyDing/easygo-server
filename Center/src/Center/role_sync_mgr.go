/**
作者:guangbo
模块：角色同步模块
说明：定时将缓存数据同步到mysql中
创建时间：2015-11-2
**/
package main

import (
	"strconv"
	"time"

	"gopkg.in/redis.v3"
)

import (
	. "GxMisc"
	. "GxStatic"
)

//Sync role.
func SyncRole() {
	go SyncCreateRole()
	go DailySyncRole()
}

//Create created role.
func SyncCreateRole() {
	rdClient := PopRedisClient()                    //Get Redis client.
	defer PushRedisClient(rdClient)                 //Insert redis client.
	key := RoleCreateList + strconv.Itoa(config.Id) //Convert config.id to string type.

	for {
		//每五分钟检查是否有新数据
		t := time.NewTicker(5 * time.Minute)
		select {
		case <-t.C:
			for {
				if rdClient.SCard(key).Val() == 0 {
					break
				}

				roleId := rdClient.SRandMember(key).Val()
				rdClient.SRem(key, roleId)
				Debug("sync create role: %s", roleId)

				//common
				id, _ := strconv.Atoi(roleId)
				saveRole(rdClient, uint32(id))
			}
		}
	}
}

//Daily sync role.
func DailySyncRole() {
	rdClient := PopRedisClient()    //Get redis client.
	defer PushRedisClient(rdClient) //Insert reids client.

	for {
		now := time.Now()
		tvl := NextTime(4, 0, 0) - now.Unix()
		t := time.NewTicker(time.Duration(tvl) * time.Second)
		select {
		case <-t.C:
			var infos []*RoleLoginInfo
			GetAllRoleLogin(rdClient, uint32(config.Id), &infos)
			n1 := len(infos)
			n2 := 0
			for i := 0; i < len(infos); i++ {
				saveRole(rdClient, infos[i].RoleId)

				if (now.Unix() > infos[i].Ts) && ((now.Unix() - infos[i].Ts) >= int64(7*24*time.Hour.Seconds())) {
					n2++
					delRoleCache(infos[i].RoleId)
					//
					infos[i].Del = 1
					SaveRoleLogin(rdClient, uint32(config.Id), infos[i])
				}
			}
			Debug("Daily Sync Role, time: %s, sync-role-count: %d, del-role-cache-count: %d, next-ts", //
				TimeToStr(now.Unix()), n1, n2, TimeToStr(NextTime(4, 0, 0)))
		}
	}
}

func saveRole(client *redis.Client, roleId uint32) {
	r := new(Role)
	r.Id = roleId
	r.Get4Redis(client)
	r.Set4Mysql(Db, uint32(config.Id))
}

func delRoleCache(roleId uint32) {
	//常用信息，装备等不用清除
}

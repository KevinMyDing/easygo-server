/**
作者:guangbo
模块：网关信息更新模块
说明：
创建时间：2015-10-30
**/
package main

import (
	"sync"
	"time"
)

import (
	. "GxMisc"
	. "GxStatic"
)

var Self *GateInfo
var countMutex *sync.Mutex
var t *time.Ticker

func gate_run() {
	countMutex = new(sync.Mutex)        //Create mutex.
	t = time.NewTicker(3 * time.Second) //Create ticker in 3s.

	Self = new(GateInfo) //Create GateInfo.

	//Get config.
	id, _ := Config.Get("server").Get("id").Int()
	Self.Id = uint32(id)
	Self.Host1, _ = Config.Get("server").Get("host1").String()
	port1, _ := Config.Get("server").Get("port1").Int()
	Self.Port1 = uint32(port1)
	Self.Host2, _ = Config.Get("server").Get("host2").String()
	port2, _ := Config.Get("server").Get("port2").Int()
	Self.Port2 = uint32(port2)
	Self.Count = 0
	go func() {
		for {
			select {
			case <-t.C:
				//定时更新自己的信息到缓存中
				go func() {
					rdClient := PopRedisClient()
					defer PushRedisClient(rdClient)

					countMutex.Lock() //Mutex Lock.
					defer countMutex.Unlock()
					Self.Ts = time.Now().Unix()
					SaveGate(rdClient, Self) //Save gate data to redis.
				}()
			}
		}

	}()
}

func addClient() {
	countMutex.Lock()
	defer countMutex.Unlock()
	Self.Count++
}

func subClient() {
	countMutex.Lock()
	defer countMutex.Unlock()
	Self.Count--
}

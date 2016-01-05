/**
作者:Kyle Ding
模块：redis连接池
说明：
创建时间：2015-12-20
**/
package GxMisc

import (
	"container/list"
	"errors"
	"fmt"
	"sync"

	"gopkg.in/redis.v3"
)

var reidsClients *list.List
var reidsMutex *sync.Mutex

var redisHost string
var redisPort int
var redisDb int64

var redisCount int

func init() {
	reidsClients = list.New()
	reidsMutex = new(sync.Mutex)
	redisCount = 4
}

func ConnectRedis(host string, port int, db int64) error {
	redisHost = host
	redisPort = port
	redisDb = db

	for i := 0; i < redisCount; i++ {
		rdClient := redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
			Password: "",      // no password set
			DB:       redisDb, // use default DB
		})
		if rdClient == nil {
			return errors.New("connect redis fail")
		}
		reidsClients.PushBack(rdClient)
	}
	return nil
}

//Get redis client.
func PopRedisClient() *redis.Client {
	reidsMutex.Lock()         //Lock
	defer reidsMutex.Unlock() //Unlock
	if reidsClients.Len() == 0 {
		for i := 0; i < redisCount; i++ {
			rdClient := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", redisHost, redisPort),
				Password: "",      // no password set
				DB:       redisDb, // use default DB
			})
			if rdClient == nil {
				return nil
			}
			reidsClients.PushBack(rdClient) //Inserts a new element  with value  at the back of list.
		}
		redisCount += redisCount
	}

	client := reidsClients.Front().Value.(*redis.Client) //Get first client.
	reidsClients.Remove(reidsClients.Front())            //Remove it.
	return client
}

func PushRedisClient(client *redis.Client) {
	reidsMutex.Lock()
	defer reidsMutex.Unlock()

	reidsClients.PushBack(client) //Insert element.
}

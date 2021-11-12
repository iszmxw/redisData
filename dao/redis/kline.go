package redis

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"time"
)

// 创建redis key
var (
	ErrorRedisDataIsNull = errors.New("name does not exist")
	ErrorGetDataFail     = errors.New("name does not exist")
)

func CreateOrChangeKline(key string, value interface{}) {
	fullKey := getRedisKey(key)
	//修改过期时间为60s
	err := rdb.Set(fullKey, value, 60*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
	}
}

func CreateRedisData(key string, value interface{}) {
	fullKey := getRedisKey(key)
	err := rdb.Set(fullKey, value, 600*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
	}
}

func CreateHistoryKline(key string, value interface{})  {
	fullKey := getHistoryKline(key)
	err := rdb.Set(fullKey, value, 120*time.Second).Err()
	//log.Println("redis finish create or change")
	if err != nil {
		log.Println(err)
	}
}

func GetKlineHistory(key string) (string, error)  {
	fullKey := getHistoryKline(key)
	val, err := rdb.Get(fullKey).Result()
	fmt.Println(fullKey)
	if err == redis.Nil {
		log.Println("key does not exist")
		return "", redis.Nil
	} else if err != nil {
		log.Printf("get name failed, err:%v\n", err)
		return "", ErrorGetDataFail
	} else {
		//log.Println("name", val)
		return val, nil
	}
}


func GetKline(key string) (string, error) {
	fullKey := getRedisKey(key)
	//fmt.Println(fullKey)
	val, err := rdb.Get(fullKey).Result()
	if err == redis.Nil {
		log.Println("key does not exist")
		return "", redis.Nil
	} else if err != nil {
		log.Printf("get name failed, err:%v\n", err)
		return "", ErrorGetDataFail
	} else {
		//log.Println("name", val)
		return val, nil
	}
}

//修改redis key
//func ChangeKline(key, value string) {
//
//}

//根据key获取值

// ExistKey 判断key是否存在
func ExistKey(key string) bool {
	result, err := rdb.Exists(key).Result()
	if err != nil {
		return true
	}
	if result == 1 {
		return true
	}
	if result == 0 {
		return false
	}
	return true
}

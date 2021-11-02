package redis

import (
	"errors"
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
	err := rdb.Set(fullKey, value, 600*time.Second).Err()
	log.Println(err)
}

func GetKline(key string) (string, error) {
	fullKey := getRedisKey(key)
	val, err := rdb.Get(fullKey).Result()
	if err == redis.Nil {
		log.Println("name does not exist")
		return "", ErrorRedisDataIsNull
	} else if err != nil {
		log.Printf("get name failed, err:%v\n", err)
		return "", ErrorGetDataFail
	} else {
		log.Println("name", val)
		return val, nil
	}
}

//修改redis key
//func ChangeKline(key, value string) {
//
//}

//根据key获取值

//判断key是否存在

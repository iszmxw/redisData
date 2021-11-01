package redis

import (
	"log"
	"time"
)

// 创建redis key

func CreateKline(key string,value interface{})  {
	fullKey := getRedisKey(key)
	err := rdb.Set(fullKey,value,60*time.Second).Err()
	log.Println(err)
}

//修改redis key

func ChangeKline(key ,value string)  {

}

//根据key获取值

func GetKline(key ,value string)  {

}

//判断key是否存在

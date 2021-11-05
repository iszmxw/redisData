package redis

const (
	Prefix = "kline:" // 项目key前缀
)

// 给redis key加上前缀
func getRedisKey(key string) string {
	return Prefix + key
}

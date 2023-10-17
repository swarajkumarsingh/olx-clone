package conf

import (
	"fmt"
	"os"
	"time"
)

/*
Redis Configurations
*/

const FreedomRedisTTL = time.Hour * 24

func GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
}

func getRedisAddrTemp() string {
	return "host.docker.internal:6379"
}

var RedisConf = map[string]interface{}{
	"Addr":     getRedisAddrTemp(),
	"SSL":      ENV == ENV_PROD,
	"Username": os.Getenv("REDIS_USER"),
	"Password": os.Getenv("REDIS_PASSWORD"),
}
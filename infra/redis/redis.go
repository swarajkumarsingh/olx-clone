/*
To do any of the things below, first import the package

import "kyc-orchestrator/infra/redis"

--------------------------------------------------------
Set key and value of type string
--------------------------------------------------------

err = Set(key, value)
if err != nil {
    log.Println("error in inserting")
}

-----------------------------------------------------------
Set key and value of struct type
-----------------------------------------------------------

err = SetStruct(key, obj, ttl)
if err != nil {
    log.Println("error in inserting")
}

--------------------------------------------
Get value of type string (non struct obj)
--------------------------------------------

val, err = Get(key)
if err != nil {
    log.Println("error in fetching")
}

--------------------------------------------
Get value of type string (struct obj)
--------------------------------------------

val, err = Get(key)
if err != nil {
    log.Println("error in fetching")
}

rawIn := json.RawMessage(val)
bytes, err := rawIn.MarshalJSON()
if err != nil {
    return nil, err
}

// obj is an instance of a struct
err = json.Unmarshal(bytes, &obj)
if err != nil {
    return nil, err
}

*/

package redis

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"olx-clone/conf"
	"time"

	redis "github.com/go-redis/redis/v8"
)

var (
	ctx = context.Background()
	rdb *redis.Client
)

func init() {

	enableSSL, _ := conf.RedisConf["SSL"].(bool)
	endpoint, _ := conf.RedisConf["Addr"].(string)
	userName, _ := conf.RedisConf["Username"].(string)
	password, _ := conf.RedisConf["Password"].(string)
	redisOption := &redis.Options{
		Addr:         endpoint,
		DialTimeout:  10 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		PoolSize:     20,
		Username:     userName,
		Password:     password,
		PoolTimeout:  30 * time.Second,
	}
	if enableSSL {
		redisOption.TLSConfig = &tls.Config{}
	}
	rdb = redis.NewClient(redisOption)
}

var suffix = "-kycapp"

func Set(key string, value string, ttl time.Duration) error {
	key = key + suffix
	err := rdb.Set(ctx, key, value, ttl).Err()

	if err != nil {
		return err
	}

	return nil
}

func GetTTL(key string) (time.Duration, error) {
	ttl, err := rdb.TTL(ctx, key).Result()
	return ttl, err
}

func SetStruct(key string, obj interface{}, ttl time.Duration) error {
	key = key + suffix
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err1 := rdb.Set(ctx, key, string(json), ttl).Err()
	if err1 != nil {
		return err1
	}

	return nil
}

func SetStructWithLongTTL(key string, obj interface{}) error {
	key = key + suffix
	json, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	err1 := rdb.Set(ctx, key, string(json), conf.FreedomRedisTTL).Err()
	if err1 != nil {
		return err1
	}

	return nil
}

func Get(key string) (string, error) {
	key = key + suffix
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

func Delete(key string) error {
	key = key + suffix
	_, err := rdb.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}

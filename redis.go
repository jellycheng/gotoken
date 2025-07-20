package gotoken

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/jellycheng/gosupport"
	"time"
)

type RedisCfg struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Db       string `json:"db"`
	Prefix   string `json:"prefix"`
}

type MyRedisClient struct {
	rdb *redis.Client
	cfg RedisCfg
}

func (m MyRedisClient) GetCfg() RedisCfg {
	return m.cfg
}

func (m MyRedisClient) GetRedisClient() *redis.Client {
	return m.rdb
}

var ctx = context.Background()
var rdbObjMap = make(map[string]*MyRedisClient)

func NewRedisClient(cfg RedisCfg) *MyRedisClient {
	myRedis := &MyRedisClient{
		cfg: cfg,
	}
	k := gosupport.Md5V1(fmt.Sprintf("%s%s%s%s", cfg.Host, cfg.Port, cfg.Username, cfg.Password))
	if r, ok := rdbObjMap[k]; ok {
		return r
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Username: cfg.Username,
		Password: cfg.Password,
		DB:       gosupport.Str2Int(cfg.Db),
	})
	myRedis.rdb = rdb
	rdbObjMap[k] = myRedis
	return myRedis
}

func SetKeyValue(myRedis *MyRedisClient, key, value string, expiration time.Duration) error {
	tmpKey := fmt.Sprintf("%s%s", myRedis.cfg.Prefix, key)
	err := myRedis.rdb.Set(ctx, tmpKey, value, expiration).Err()
	return err
}

func GetKeyValue(myRedis *MyRedisClient, key string) string {
	tmpKey := fmt.Sprintf("%s%s", myRedis.cfg.Prefix, key)
	val, _ := myRedis.rdb.Get(ctx, tmpKey).Result()
	return val
}

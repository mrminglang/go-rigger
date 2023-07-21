package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type RedisConfig struct {
	RedisIp          string
	RedisPort        int32
	RedisPassword    string
	RedisDbIndex     int    // 默认数据库
	RedisPrefix      string // 默认前缀
	ResisMaxIdle     int    // 最大空闲链接数量
	RedisMaxActive   int    // 最大链接数，0表示并发不限制
	RedisIdleTimeout string // 最大空闲时间，用完连接后 n 秒回收到链接池
}

// GetRedisConfig redis config
func GetRedisConfig() *RedisConfig {
	var redis = &RedisConfig{}
	common.DBredis = storage.Storage.FullPath(common.DBredis) // 单元测试需要使用绝对路径
	glog.Channel("redis").Infoln("GetRedisConfig common path::", common.DBredis)
	isExist, _ := files.IsFileExist(common.DBredis)
	if !isExist {
		glog.Channel("db").Errorln("GetRedisConfig is not exist...")
		return redis
	}
	info, err := ioutil.ReadFile(common.DBredis)
	if err != nil {
		return redis
	}
	err = json.Unmarshal(info, redis)
	return redis
}

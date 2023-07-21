package gredis

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/mrminglang/go-rigger/config"
	"github.com/zhan3333/glog"
	"strconv"
	"strings"
	"time"
)

type RedisConn struct {
	pool          *redis.Pool
	bConnected    bool
	defaultPrefix string
}

var defaultRedisConn *RedisConn

// GetRedisCon 获取redis链接对象
func GetRedisCon() *RedisConn {
	return defaultRedisConn
}

// CloseRedisCon 关闭redis链接对象
func CloseRedisCon() error {
	return defaultRedisConn.DisconnectServer()
}

// Init 初始化Redis链接
func Init() {
	glog.Channel("redis").Infoln("gredis connect start....")
	dbConfig := config.GetRedisConfig()
	glog.Channel("redis").Infoln("gredis config::", dbConfig)

	if defaultRedisConn == nil {
		defaultRedisConn = new(RedisConn)
		err := defaultRedisConn.ConnectServer(dbConfig)
		if err != nil {
			glog.Channel("redis").Infoln("gredis connect error::", err.Error())
			return
		}
		glog.Channel("redis").Infoln("gredis connect success....")
	} else {
		glog.Channel("redis").Infoln("gredis connect is exist....")
	}
}

// ConnectServer 链接reids链接池
func (base *RedisConn) ConnectServer(config *config.RedisConfig) (err error) {
	if err = base.DisconnectServer(); err != nil {
		glog.Channel("redis").Errorln("disconnect from original base server failed")
	}
	duration, _ := time.ParseDuration(config.RedisIdleTimeout)

	base.pool = &redis.Pool{
		MaxIdle:     config.ResisMaxIdle,   // 最大空闲链接数量
		MaxActive:   config.RedisMaxActive, // 最大链接数
		IdleTimeout: duration,              // 最大空闲时间
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			con, err := redis.Dial("tcp", config.RedisIp+":"+strconv.FormatInt(int64(config.RedisPort), 10))
			if err != nil {
				glog.Channel("redis").Errorln("connect redis server Dial error:", err.Error())
				return nil, err
			}
			if config.RedisPassword != "" {
				if _, err = con.Do("AUTH", config.RedisPassword); err != nil {
					glog.Channel("redis").Errorln("connect redis server AUTH error:", err.Error())
					return nil, err
				}
			}
			_, err = con.Do("SELECT", config.RedisDbIndex)
			if err != nil {
				glog.Channel("redis").Errorln("connect redis server SELECT error:", err.Error())
				return nil, err
			}
			return con, nil
		},
	}

	base.bConnected = true
	base.defaultPrefix = config.RedisPrefix
	return
}

// DisconnectServer 释放redis链接池中所有链接
func (base *RedisConn) DisconnectServer() (err error) {
	if base.bConnected {
		err = base.pool.Close()
		if err != nil {
			glog.Channel("redis").Errorln("disconnect from base server failed")
			return
		}

		base.bConnected = false
		glog.Channel("redis").Infoln("disconnect from base server ok")
		return
	}

	return
}

// GetCon 获取单个链接
func (base *RedisConn) GetCon() redis.Conn {
	return base.pool.Get()
}

// Get 取key值
func (base *RedisConn) Get(key string, prefix string) (string, error) {
	if !base.bConnected || base == nil {
		return "", errors.New("connect redis is not exist")
	}

	key = base.CompleteKey(key, prefix)
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()
	glog.Channel("redis").Infoln("redis get key::", key)
	val, err := redis.String(con.Do("GET", key))
	if err != nil {
		glog.Channel("redis").Errorln("redis get key:" + key + "err:" + err.Error())
		return "", err
	} else {
		if val == "" {
			glog.Channel("redis").Errorln("redis get key:" + key + "is nil")
			return "", nil
		}
		return val, nil
	}
}

// Set 存在key-val
func (base *RedisConn) Set(key string, val string, prefix string) error {
	if !base.bConnected || base == nil {
		return errors.New("connect redis is not exist")
	}

	key = base.CompleteKey(key, prefix)
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()
	glog.Channel("redis").Infoln("redis set key::", key)
	_, err := redis.String(con.Do("set", key, val))
	if err != nil {
		glog.Channel("redis").Errorln("redis set key:" + key + ",value:" + val + "err:" + err.Error())
		return err
	}
	return nil
}

// Del 删除key-val
func (base *RedisConn) Del(key string, prefix string) bool {
	if !base.bConnected || base == nil {
		return false
	}

	key = base.CompleteKey(key, prefix)
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()
	glog.Channel("redis").Infoln("redis del key::", key)
	val, err := redis.Bool(con.Do("del", key))
	if err != nil {
		glog.Channel("redis").Errorln("redis del key:" + key + "err:" + err.Error())
		return false
	}
	return val
}

// HSet 保存key val
func (base *RedisConn) HSet(mapName string, key string, val string) (err error) {
	if !base.bConnected || base == nil {
		return errors.New("connect is not exist")
	}

	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()
	_, err = con.Do("HSET", mapName, key, val)
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("hSet mapName:%s key:%s val:%s err:%s", mapName, key, val, err.Error()))
		return
	}
	return
}

// HGet 得到key val
func (base *RedisConn) HGet(mapName string, key string) (val string, err error) {
	if !base.bConnected || base == nil {
		return "", errors.New("connect is not exist")
	}

	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	val, err = redis.String(con.Do("HGET", mapName, key))
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("hGet mapName:%s key:%s err:%s", mapName, key, err.Error()))
		return
	}
	return
}

// HDel 删除key val
func (base *RedisConn) HDel(mapName string, key string) bool {
	if !base.bConnected || base == nil {
		return false
	}

	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	_, err := con.Do("HDEL", mapName, key)
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("hDel mapName:%s key:%s err:%s", mapName, key, err.Error()))
		return false
	}

	return true
}

// HGetAll 获取key所有字段
func (base *RedisConn) HGetAll(mapName string) (maps map[string]string, err error) {
	if !base.bConnected || base == nil {
		return nil, errors.New("connect is not exist")
	}

	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	maps, err = redis.StringMap(con.Do("HGETALL", mapName))
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("HGetAll mapName:%s err:%s", mapName, err.Error()))
		return
	}
	return
}

// SAdd 添加member
func (base *RedisConn) SAdd(key string, member string) (err error) {
	if !base.bConnected || base == nil {
		return errors.New("connect is not exist")
	}

	if !strings.HasPrefix(key, base.defaultPrefix) {
		key = base.defaultPrefix + key
	}
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	_, err = con.Do("sadd", key, member)
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("SAdd key:%s member:%s err:%s", key, member, err.Error()))
		return
	}

	return
}

// SMembers 获取所有member
func (base *RedisConn) SMembers(key string) (res []string, err error) {
	if !base.bConnected || base == nil {
		return nil, errors.New("connect is not exist")
	}
	if !strings.HasPrefix(key, base.defaultPrefix) {
		key = base.defaultPrefix + key
	}
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()
	res, err = redis.Strings(con.Do("smembers", key))
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("SMembers key:%s err:%s", key, err.Error()))
		return
	}
	return
}

// SCard 获取key中member数量
func (base *RedisConn) SCard(key string) (num int, err error) {
	if !base.bConnected || base == nil {
		return 0, errors.New("connect is not exist")
	}

	if !strings.HasPrefix(key, base.defaultPrefix) {
		key = base.defaultPrefix + key
	}

	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	num, err = redis.Int(con.Do("scard", key))
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("SCard key:%s err:%s", key, err.Error()))
		return
	}
	return
}

// GetAllKey 获取所有key
func (base *RedisConn) GetAllKey(prefix string) (keys []string, err error) {
	if !base.bConnected || base == nil {
		return nil, errors.New("connect is not exist")
	}
	if prefix == "" {
		prefix = base.defaultPrefix
	}
	glog.Channel("redis").Infoln("GetAllKey prefix::", prefix)
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	keys, err = redis.Strings(con.Do("keys", prefix+"*"))
	if err != nil {
		glog.Channel("redis").Errorln(fmt.Sprintf("GetAllKey prefix:%s err:%s", prefix, err.Error()))
		return
	}
	return
}

// SetLock 设置分布式锁
func (base *RedisConn) SetLock(key string, val string, timeSecond int, prefix string) bool {
	glog.Channel("redis").Infoln("set lock " + key + "," + val)

	key = base.CompleteKey(key, prefix)
	con := base.GetCon()
	defer func() {
		_ = con.Close()
	}()

	_, err := con.Do("set", key, val, "EX", timeSecond, "NX")
	if err != nil {
		glog.Channel("redis").Errorln("set key:" + key + ",value:" + val + "分布式锁异常:" + err.Error())
		glog.Channel("redis").Infoln("分布式锁异常描述::" + err.Error())
		return false
	} else {
		glog.Channel("redis").Infoln("设置分布式锁成功...")
		return true
	}
}

// ReleaseLock 释放分布式锁
func (base *RedisConn) ReleaseLock(key string, val string, prefix string) bool {
	redisVal, _ := base.Get(key, prefix)
	if redisVal == val {
		res := base.Del(key, prefix)
		glog.Channel("redis").Infoln("释放分布式锁成功...")
		return res
	} else {
		return false
	}
}

// CompleteKey 获取完整的key
func (base *RedisConn) CompleteKey(key, prefix string) string {
	if prefix != "" {
		key = prefix + key
	} else {
		if !strings.HasPrefix(key, base.defaultPrefix) {
			key = base.defaultPrefix + key
		}
	}

	return key
}

package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type MongoConfig struct {
	DBHost          string
	DBPort          int32
	DBUser          string
	DBPWd           string
	DBName          string
	DBType          string
	DBBatchSize     int
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime string
}

// GetMongoConfig Mongo config
func GetMongoConfig() *MongoConfig {
	var mongo = &MongoConfig{}
	common.Mongo = storage.Storage.FullPath(common.Mongo) // 单元测试需要使用绝对路径
	glog.Channel("mongo").Infoln("GetMongoConfig common path::", common.Mongo)
	isExist, _ := files.IsFileExist(common.Mongo)
	if !isExist {
		glog.Channel("mongo").Errorln("GetMongoConfig is not exist...")
		return mongo
	}
	info, err := ioutil.ReadFile(common.Mongo)
	if err != nil {
		return mongo
	}
	err = json.Unmarshal(info, mongo)
	return mongo
}

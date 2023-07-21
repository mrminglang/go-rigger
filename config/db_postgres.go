package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type PostgresqlConfig struct {
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

func GetPostgresqlConfig() *PostgresqlConfig {
	var postgresql = &PostgresqlConfig{}
	common.DBpostgres = storage.Storage.FullPath(common.DBpostgres) // 单元测试需要使用绝对路径
	glog.Channel("db").Infoln("GetPostgresqlConfig common path::", common.DBpostgres)
	isExist, _ := files.IsFileExist(common.DBpostgres)
	if !isExist {
		glog.Channel("db").Errorln("GetPostgresqlConfig is not exist...")
		return postgresql
	}
	info, err := ioutil.ReadFile(common.DBpostgres)
	if err != nil {
		return postgresql
	}
	err = json.Unmarshal(info, postgresql)
	return postgresql
}

package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type MySQLConfig struct {
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

// GetMySQLConfig mysql config
func GetMySQLConfig() *MySQLConfig {
	var mysql = &MySQLConfig{}
	common.DBmysql = storage.Storage.FullPath(common.DBmysql) // 单元测试需要使用绝对路径
	glog.Channel("db").Infoln("GetMySQLConfig common path::", common.DBmysql)
	isExist, _ := files.IsFileExist(common.DBmysql)
	if !isExist {
		glog.Channel("db").Errorln("GetMySQLConfig is not exist...")
		return mysql
	}
	info, err := ioutil.ReadFile(common.DBmysql)
	if err != nil {
		return mysql
	}
	err = json.Unmarshal(info, mysql)
	return mysql
}

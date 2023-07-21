package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type EnvConfig struct {
	GinPort int32
}

// GetEnvConfig env config
func GetEnvConfig() *EnvConfig {
	var env = &EnvConfig{}

	common.DBmysql = storage.Storage.FullPath(common.Env) // 单元测试需要使用绝对路径
	glog.Channel("db").Infoln("GetEnvConfig common path::", common.Env)
	isExist, _ := files.IsFileExist(common.Env)
	if !isExist {
		glog.Channel("db").Errorln("GetEnvConfig is not exist...")
		return env
	}
	info, err := ioutil.ReadFile(common.Env)
	if err != nil {
		return env
	}
	err = json.Unmarshal(info, env)
	return env
}

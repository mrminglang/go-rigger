package config

import (
	"encoding/json"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/mrminglang/tools/files"
	"github.com/zhan3333/glog"
	"io/ioutil"
)

type RabbitMQConfig struct {
	MqIp             string
	MqPort           string
	MqUser           string
	MqPassword       string
	MqMonitorTimeout int
	MqExchange       string
	MqKing           string // 交换机
}

// GetRabbitMQConfig MQ config
func GetRabbitMQConfig() *RabbitMQConfig {
	var mq = &RabbitMQConfig{}
	common.DBrabbitMQ = storage.Storage.FullPath(common.DBrabbitMQ) // 单元测试需要使用绝对路径
	glog.Channel("mq").Infoln("GetRabbitMQConfig common path::", common.DBrabbitMQ)
	isExist, _ := files.IsFileExist(common.DBrabbitMQ)
	if !isExist {
		glog.Channel("mq").Errorln("GetRabbitMQConfig is not exist...")
		return mq
	}
	info, err := ioutil.ReadFile(common.DBrabbitMQ)
	if err != nil {
		return mq
	}
	err = json.Unmarshal(info, mq)
	return mq
}

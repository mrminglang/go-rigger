package config

import (
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/tools/paths"
	"github.com/zhan3333/glog"
)

var Logging = struct {
	Channels map[string]glog.Log
	Default  string
}{
	Channels: map[string]glog.Log{
		"default": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/def/def.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 默认日志
		"db": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/db/db.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 数据库日志
		"redis": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/redis/redis.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 数据库日志
		"gin": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/gin/gin.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 数据库日志
		"mq": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/mq/mq.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 数据库日志
		"mongo": {
			Driver: glog.DAILY,
			Path:   paths.FullStoragePath(common.Storage, "logs/mongo/mongo.log"),
			Level:  glog.DebugLevel,
			Days:   7,
		}, // 数据库日志
	},
	Default: "default",
}

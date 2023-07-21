package logging

import (
	"fmt"
	"github.com/mrminglang/go-rigger/common"
	"github.com/mrminglang/go-rigger/config"
	"github.com/zhan3333/glog"
	"time"
)

var disableBootScreenLog bool

func DisableBootScreenLog() {
	disableBootScreenLog = true
}

// LogBootInfo 详情
func LogBootInfo(info string) {
	datetime := time.Now().Format(common.Layout)
	glog.Def().Infof("[%s] boot: %s", datetime, info)
	if !disableBootScreenLog {
		fmt.Printf("[%s] boot: %s\n", datetime, info)
	}
}

// LogBootPanic 宕机
func LogBootPanic(msg string, err error) {
	datetime := time.Now().Format(common.Layout)
	if !disableBootScreenLog {
		fmt.Printf("[%s] boot: %s: %+v\n", datetime, msg, err)
	}
	glog.Def().WithError(err).Panicf("[%s] boot: %s: %+v", datetime, msg, err)
}

// LogBoot 日志启动
func LogBoot() {
	glog.DefLogChannel = config.Logging.Default
	glog.LogConfigs = config.Logging.Channels
	glog.LoadChannels()
	LogBootInfo("logging module init")
}

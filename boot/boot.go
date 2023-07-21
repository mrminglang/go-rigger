package boot

import (
	"github.com/mrminglang/go-rigger/connect/gmongo"
	"github.com/mrminglang/go-rigger/connect/logging"
	"github.com/mrminglang/go-rigger/storage"
	"github.com/zhan3333/glog"
	"os"
	"path/filepath"
)

func Boot() {
	// 启动日志
	logging.DisableBootScreenLog()
	logging.LogBoot()
	glog.Def().Infof("boot start......")

	// 初始化Storage
	bootStorage()

	// 注册数据库连接
	//gmysql.Init()
	//gpostgres.Init()
	//gmysql.InitGorm()
	//gpostgres.InitGorm()
	gmongo.Init()

	// 注册Redis链接
	//gredis.Init()

	// 注册rabbitMQ链接
	//grabbitmq.Init()

	// 注册gin服务
	//cmd.RunServer()

	glog.Def().Infof("boot success......")
}

func Destroy() {
	//_ = gmysql.Close()
	//_ = gpostgres.Close()

	//_ = gredis.CloseRedisCon()
	//_ = grabbitmq.CloseRabbitCon()

	gmongo.Close()
	glog.Close()
	glog.Def().Infof("boot close......")
}

// RootPath 获取项目路径
func RootPath() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return ""
		}
		dir = parent
	}
}

func bootStorage() {
	storage.Init(RootPath())
	glog.Def().Infof("storage module init")
}

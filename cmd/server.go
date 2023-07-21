package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrminglang/go-rigger/config"
	"github.com/mrminglang/go-rigger/middleware"
	"github.com/zhan3333/glog"
)

func RunServer() {
	baseServer := "0.0.0.0:" + fmt.Sprintf("%d", config.GetEnvConfig().GinPort)
	glog.Channel("gin").Infoln("start server...\r\ngo：http://" + baseServer)
	engine := gin.Default()
	engine.Use(middleware.Cors())

	err := engine.Run(baseServer)
	if err != nil {
		glog.Channel("gin").Infoln("engine.Run(baseServer) error::", err.Error())
	}
}

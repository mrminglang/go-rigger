package gmongo

import (
	"fmt"
	"github.com/mrminglang/go-rigger/config"
	"github.com/zhan3333/glog"
	"gopkg.in/mgo.v2"
)

var MongoDBConn *DBOperation

type DBOperation struct {
	Session *mgo.Session
	DBName  string
}

// Init 初始化mongo数据库连接 支持数据库连接池
func Init() {
	glog.Channel("mongo").Infoln("mongo connect start....")
	dbConfig := config.GetMongoConfig()
	glog.Channel("mongo").Infoln("mongo dbConfig::", dbConfig)
	// mongodb://host:port/database
	dsn := fmt.Sprintf("mongodb://%s:%d/%s", dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
	session, err := mgo.Dial(dsn)
	if err != nil {
		glog.Channel("mongo").Errorln(fmt.Sprintf("mongo mgo.Dial %s fire error:::%s", dsn, err.Error()))
		return
	}
	// 登录认证
	if dbConfig.DBUser != "" && dbConfig.DBPWd != "" {
		err = session.Login(&mgo.Credential{
			Username: dbConfig.DBUser,
			Password: dbConfig.DBPWd,
		})
		if err != nil {
			glog.Channel("mongo").Errorln(fmt.Sprintf("mongo mgo.Dial %s Login error:::%s", dsn, err.Error()))
			return
		}
	}

	//session.SetMode(mgo.Eventual) //该模式最大的特点就是不会缓存连接
	session.SetPoolLimit(dbConfig.MaxOpenConns) //设置连接缓冲池的最大值

	if MongoDBConn == nil {
		MongoDBConn = new(DBOperation)
		MongoDBConn.Session = session
		MongoDBConn.DBName = dbConfig.DBName
		glog.Channel("mongo").Infoln("gpostgres connect success....", dsn)
	}
}

// Close 关闭 postgresql 数据库连接
func Close() {
	MongoDBConn.Session.Close()
}

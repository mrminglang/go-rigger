package gmysql

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/mrminglang/go-rigger/config"
	"github.com/zhan3333/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DbMySQLConn     *sql.DB
	DbMySQLGormConn *gorm.DB
	DbBatchSize     int
)

// Init 初始化MySQL数据库
func Init() {
	glog.Channel("db").Infoln("gmysql connect start....")
	dbConfig := config.GetMySQLConfig()
	glog.Channel("db").Infoln("gmysql config::", dbConfig)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true",
		dbConfig.DBUser, dbConfig.DBPWd, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)
	db, err := sql.Open(dbConfig.DBType, dsn)
	if err != nil {
		glog.Channel("db").Errorln(fmt.Sprintf("gmysql sql.open %s fire error::%s", dsn, err.Error()))
		return
	}

	// 设置最大开放连接数，注意该值为小于0或等于0指的是无限制连接数
	db.SetMaxOpenConns(dbConfig.MaxOpenConns)

	// 设置空闲连接数，小于等于0表示无限制
	db.SetMaxIdleConns(dbConfig.MaxIdleConns)

	// 重复使用连接的最大时间字符串解析为time.Duration类型
	duration, _ := time.ParseDuration(dbConfig.ConnMaxLifetime)
	// 设置可以重复使用连接的最大时间，设置为0，表示没有最大生存期，并且连接会被重用
	db.SetConnMaxLifetime(duration)

	// 创建一个具有5秒超时期限的上下文。
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	//使用PingContext()建立到数据库的新连接，并传入上下文信息，连接超时就返回
	err = db.PingContext(ctx)
	if err != nil {
		glog.Channel("db").Errorln(fmt.Sprintf("gmysql db.PingContext fire error::%s", err.Error()))
		return
	}

	//err = db.Ping()
	//if err != nil {
	//	logrus.Errorln(fmt.Sprintf("gmysql db.Ping() fire error::%s", err.Error()))
	//	return
	//}

	if DbMySQLConn == nil {
		DbMySQLConn = db

		glog.Channel("db").Infoln("gmysql connect success....", dsn)
	}
}

// Close 关闭MySQL数据库连接
func Close() error {
	return DbMySQLConn.Close()
}

// InitGorm 初始化数据库连接 gorm
func InitGorm() {
	glog.Channel("db").Infoln("gmysql connect start....")
	dbConfig := config.GetMySQLConfig()
	glog.Channel("db").Infoln("gmysql dbConfig::", dbConfig)

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=Local",
		dbConfig.DBUser, dbConfig.DBPWd, dbConfig.DBHost, dbConfig.DBPort, dbConfig.DBName)

	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.New(glog.Channel("db"), logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log level
			Colorful:      false,       // 彩色打印
		}),
		DisableForeignKeyConstraintWhenMigrating: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "public.",
			SingularTable: true, //表名后面不加s
		}})
	if err != nil {
		glog.Channel("db").Errorln("gmysql gorm.Open() error::", err.Error())
		return
	}
	sqlDB, err := conn.DB()
	if err != nil {
		glog.Channel("db").Errorln("gmysql conn.DB() error::", err.Error())
		return
	}

	sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
	sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
	duration, _ := time.ParseDuration(dbConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxLifetime(duration)

	if DbMySQLGormConn == nil {
		DbMySQLGormConn = conn
		DbBatchSize = dbConfig.DBBatchSize
		glog.Channel("db").Infoln("gmysql connect success....")
	}
}

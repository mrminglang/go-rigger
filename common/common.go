package common

import "github.com/mrminglang/tools/paths"

var (
	Dir        string = "netmcset/"
	DBpostgres string = Dir + "postgres.json"
	DBmysql    string = Dir + "mysql.json"
	DBredis    string = Dir + "redis.json"
	DBrabbitMQ string = Dir + "rabbitmq.json"
	Mongo      string = Dir + "mongo.json"
	Env        string = Dir + "env.json"
	Storage    string = paths.FullPath("storage")
)

const Layout = "2006-01-02 15:04:05"

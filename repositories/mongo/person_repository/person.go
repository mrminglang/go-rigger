package person_repository

import (
	"github.com/mrminglang/go-rigger/connect/gmongo"
	"github.com/mrminglang/go-rigger/repositories/models"
	"github.com/zhan3333/glog"
)

type person struct{}

func NewPerson() *person {
	return &person{}
}

// CreatePerson 创建用户
func (u *person) CreatePerson(person *models.Person) (err error) {
	table := models.Person{}
	tableName := table.TableName()
	dbName := gmongo.MongoDBConn.DBName
	glog.Channel("mongo").Errorln("{mongo CreatePerson DBName^tableName^person}", dbName, tableName, person)
	err = gmongo.MongoDBConn.Session.DB(dbName).C(tableName).Insert(person)
	if err != nil {
		glog.Channel("mongo").Errorln("mongo CreatePerson error::", err.Error())
		return
	}

	return
}

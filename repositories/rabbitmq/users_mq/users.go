package users_mq

import (
	"github.com/mrminglang/go-rigger/repositories/rabbitmq"
	"github.com/zhan3333/glog"
)

type UserInfo struct {
	UserID   string
	UserName string
}

// SendUserToMQ 发送数据 生产者
func SendUserToMQ(topic string, user *UserInfo) bool {
	//b, _ := json.Marshal(user)
	b := "1234456"
	rabbitmq.SendDataToMQ(topic, []byte(b))
	return true
}

// RecvUserFromMQ 消费数据
func RecvUserFromMQ(topic string) bool {
	rabbitmq.RecieveDataFromMQ(topic, "json", OnHandlerUserInfo)
	return true
}

func OnHandlerUserInfo(user *UserInfo) error {
	glog.Channel("mq").Infoln("user info::", user)
	return nil
}

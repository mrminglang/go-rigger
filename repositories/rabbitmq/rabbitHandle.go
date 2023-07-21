package rabbitmq

import (
	"github.com/mrminglang/go-rigger/connect/grabbitmq"
	"github.com/zhan3333/glog"
)

// SendDataToMQ 发送数据到MQ
func SendDataToMQ(topic string, msg []byte) bool {
	mq := grabbitmq.GetRabbitCon()
	if mq.Conn == nil {
		glog.Channel("mq").Errorln("SendDataToMQ rabbitmq链接对象不存在...")
		return false
	}
	glog.Channel("mq").Infoln("SendDataToMQ topic::" + topic + " msq::" + string(msg))
	mq.GetSendRabbit(topic).SendMsg(topic, topic, msg)
	return true
}

// RecieveDataFromMQ 接收数据从MQ
func RecieveDataFromMQ(topic, protocol string, handler interface{}) bool {
	mq := grabbitmq.GetRabbitCon()
	if mq.Conn == nil {
		glog.Channel("mq").Errorln("SendDataToMQ rabbitmq链接对象不存在...")
		return false
	}

	mq.GetReceiveRabbit(topic, protocol, handler)
	return true
}

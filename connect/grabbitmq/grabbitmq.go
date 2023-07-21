package grabbitmq

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/mrminglang/go-rigger/config"
	"github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	"github.com/zhan3333/glog"
	"reflect"
	"strconv"
	"time"
)

var defaultRabbitConn *RabbitConn

// RabbitConn 链接对象结构体
type RabbitConn struct {
	Conn        *amqp.Connection // 连接
	Exchange    string
	King        string
	MonitorTime int
}

type Handler struct {
	Method  reflect.Value //对应方法指针
	ReqType reflect.Type  //一般是实现接口proto.Message
	CtxType reflect.Type  //一般是context
}

// RabbitChannel channel对象结构体
type RabbitChannel struct {
	Topic     string      //
	Handler   *Handler    //业务处理对象
	Cancel    bool        //是否退出监听
	RabbitCon *RabbitConn //rabbitmq链接对象
}

// GetRabbitCon 获取rabbit链接对象
func GetRabbitCon() *RabbitConn {
	return defaultRabbitConn
}

// CloseRabbitCon 关闭rabbit链接对象
func CloseRabbitCon() error {
	return defaultRabbitConn.Conn.Close()
}

// Init 初始化 rabbitMQ 链接
func Init() *RabbitConn {
	glog.Channel("mq").Infoln("grabbitMQ connect start....")
	if defaultRabbitConn != nil {
		glog.Channel("mq").Infoln("grabbitMQ connect is exist....")
	}
	// 获取配置消息
	conf := config.GetRabbitMQConfig()
	glog.Channel("mq").Infoln("grabbitMQ config::", conf)

	defaultRabbitConn = new(RabbitConn)
	defaultRabbitConn.StartMqService(conf)
	glog.Channel("mq").Infoln("grabbitMQ connect success....")

	return defaultRabbitConn
}

// StartMqService 开启链接和监听MQ
func (mq *RabbitConn) StartMqService(config *config.RabbitMQConfig) {
	mq.CreateConnection(config)
	mq.Reconnect(config)
	_ = mq.CrateExchange(config.MqExchange, config.MqKing)
}

// CreateConnection 创建MQ链接
func (mq *RabbitConn) CreateConnection(config *config.RabbitMQConfig) *amqp.Connection {
	// 	amqp://guest:guest@mqIp:mqPort/
	dns := "amqp://" + config.MqUser + ":" + config.MqPassword + "@" + config.MqIp + ":" + config.MqPort
	conn, err := amqp.Dial(dns)
	if err != nil {
		glog.Channel("mq").Errorln("rabbitmq CreateConnection 链接错误..." + err.Error())
		return nil
	}

	mq.Conn = conn
	mq.Exchange = config.MqExchange
	mq.King = config.MqKing
	mq.MonitorTime = config.MqMonitorTimeout

	return conn
}

// Reconnect 处理断网重新链接MQ
func (mq *RabbitConn) Reconnect(config *config.RabbitMQConfig) {
	go func() {
		for {
			if mq.Conn == nil || mq.Conn.IsClosed() {
				conn := mq.CreateConnection(config)
				if conn == nil {
					glog.Channel("mq").Infoln("reconnect 监听线程重连rabbitMQ失败...")
				} else {
					glog.Channel("mq").Infoln("reconnect 监听线程重连rabbitMQ成功...")
				}
			}
			time.Sleep(time.Duration(mq.MonitorTime) * time.Second)
		}
	}()
}

// CrateExchange 创建交换器
func (mq *RabbitConn) CrateExchange(exchangeName, king string) error {
	if mq.Conn == nil {
		glog.Channel("mq").Infoln("CrateExchange rabbitmq链接对象不存在...")
		return errors.New("CrateExchange rabbitmq链接对象不存在")
	}

	channel, err := mq.Conn.Channel()
	if channel != nil {
		defer channel.Close()
	}
	if err != nil {
		glog.Channel("mq").Errorln("CrateExchange Channel 创建交换器异常 error::" + err.Error())
		return errors.New("CrateExchange Channel 创建交换器异常")
	}

	// king交换机类型
	// topic-话题模式;fanout-广播模式;direct-路由模式;
	err = channel.ExchangeDeclare(exchangeName, king, true, false, false, false, nil)
	if err != nil {
		glog.Channel("mq").Errorln("CrateExchange ExchangeDeclare 创建交换器异常 error::" + err.Error())
		return errors.New("CrateExchange ExchangeDeclare 创建交换器异常")
	}
	return nil
}

// GetSendRabbit 获取一个新的Rabbit对象发送消息
func (mq *RabbitConn) GetSendRabbit(topic string) *RabbitChannel {
	rabbit := &RabbitChannel{ //创建rabbit功能对象
		Topic:     topic,
		RabbitCon: mq,
	}
	_, _ = rabbit.CreateQueue(topic)
	_ = rabbit.CreateQueueBind(topic, topic, mq.Exchange)
	return rabbit
}

// GetReceiveRabbit 获取一个新的Rabbit对象接收消息
func (mq *RabbitConn) GetReceiveRabbit(topic string, protocol string, handler interface{}) *RabbitChannel {
	rabbit := &RabbitChannel{ //创建rabbit功能对象
		Topic:     topic,
		RabbitCon: mq,
	}
	//rabbit.msgHandler = handler //业务数据处理器绑定到rabbit功能对象中
	if typ := reflect.TypeOf(handler); typ.Kind() == reflect.Func {
		h := &Handler{
			Method: reflect.ValueOf(handler),
		}
		switch typ.NumIn() {
		case 1:
			h.ReqType = typ.In(0)
		case 2:
			h.CtxType = typ.In(0)
			h.ReqType = typ.In(1)
		}
		rabbit.Handler = h
	}

	_, _ = rabbit.CreateQueue(topic)
	_ = rabbit.CreateQueueBind(topic, topic, mq.Exchange)
	rabbit.RecvMsg(topic, protocol)

	return rabbit
}

// CreateQueue 创建队列处理
func (rabbit *RabbitChannel) CreateQueue(queueName string) (*amqp.Queue, error) {
	if rabbit.RabbitCon.Conn == nil {
		glog.Channel("mq").Infoln("CreateQueue rabbitmq链接对象不存在...")
		return nil, errors.New("CreateQueue rabbitmq链接对象不存在")
	}

	channel, err := rabbit.RabbitCon.Conn.Channel()
	if err != nil {
		glog.Channel("mq").Infoln("CreateQueue Channel 创建队列异常 error::" + err.Error())
		return nil, errors.New("CreateQueue Channel 创建队列异常")
	}
	if channel != nil {
		//defer channel.Close()
	}

	queue, err := channel.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		glog.Channel("mq").Errorln("CreateQueue QueueDeclare 创建队列异常 error::" + err.Error())
		return nil, errors.New("CreateQueue QueueDeclare 创建队列异常")
	}

	return &queue, nil
}

// DeleteQueue 删除队列处理
func (rabbit *RabbitChannel) DeleteQueue(queueName string) bool {
	if rabbit.RabbitCon.Conn == nil {
		glog.Channel("mq").Errorln("DeleteQueue 删除队列异常，rabbitmq链接对象不存在...")
		return false
	}

	channel, err := rabbit.RabbitCon.Conn.Channel()
	if channel != nil {
		defer channel.Close()
	}
	if err != nil {
		glog.Channel("mq").Errorln("DeleteQueue Channel 删除队列异常 error::", err.Error())
		return false
	}

	_, err = channel.QueueDelete(queueName, false, false, false)
	if err != nil {
		logrus.Infoln("DeleteQueue QueueDelete 删除队列异常 error::", err.Error())
		return false
	}

	glog.Channel("mq").Errorln("DeleteQueue 删除队列成功...")
	return true
}

// CreateQueueBind 创建一个队列与交换器绑定
func (rabbit *RabbitChannel) CreateQueueBind(queueName string, routingKey string, exchange string) error {
	if rabbit.RabbitCon.Conn == nil {
		glog.Channel("mq").Infoln("CreateQueueBind rabbitmq链接对象不存在...")
		return errors.New("CreateQueueBind rabbitmq链接对象不存在")
	}

	channel, err := rabbit.RabbitCon.Conn.Channel()
	if err != nil {
		glog.Channel("mq").Infoln("CreateQueueBind Channel 创建绑定异常 error::" + err.Error())
		return errors.New("CreateQueueBind Channel 创建绑定异常")
	}
	if channel != nil {
		//defer channel.Close()
	}

	err = channel.QueueBind(queueName, routingKey, exchange, false, nil)
	if err != nil {
		glog.Channel("mq").Infoln("CreateQueueBind QueueBind 创建绑定异常 error::" + err.Error())
		return errors.New("CreateQueueBind QueueBind 创建绑定异常")
	}

	return nil
}

// SendMsg 发送数据 生产者
func (rabbit *RabbitChannel) SendMsg(exchangeName string, routingKey string, msg []byte) {
	glog.Channel("mq").Infoln("SendMsg 开始发送消息......")
	go func() {
		for {
			isOK := false
			glog.Channel("mq").Infoln("SendMsg 发消息线程启动......")
			channel, err := rabbit.RabbitCon.Conn.Channel()
			if channel != nil {
				defer channel.Close()
			}

			if err == nil {
				pub := amqp.Publishing{
					Headers:      amqp.Table{},
					DeliveryMode: amqp.Persistent,
					ContentType:  "application/json",
					Body:         msg,
				}
				pub.Headers["Content-Type"] = "application/json"
				pub.Headers["Micro-Topic"] = routingKey
				err = channel.Publish(exchangeName, routingKey, false, false, pub)
				if err != nil {
					glog.Channel("mq").Errorln("SendMsg Publish 发送消息到rabbitmq异常 error::" + err.Error())
				} else {
					glog.Channel("mq").Errorln("SendMsg Publish 发送消息到rabbitmq成功......")
					isOK = true
				}
			} else {
				glog.Channel("mq").Infoln("生产者=>>>通道异常，尝试重新链接 error::" + err.Error())
			}

			if isOK {
				break
			}
		}
	}()
}

// RecvMsg 消费数据 消费者
func (rabbit *RabbitChannel) RecvMsg(queueName, protocol string) {
	go func() {
		for {
			if rabbit.RabbitCon.Conn != nil {
				//重新获取信道
				channel, err := rabbit.RabbitCon.Conn.Channel()
				if err == nil {
					msg, err := channel.Consume(queueName, "", false, false, false, false, nil)
					if err == nil {
						for {
							//取通道数据
							data, isColse := <-msg
							glog.Channel("mq").Infoln("RecvMsg data Body::", data.Body)
							if !isColse {
								glog.Channel("mq").Infoln("chan关闭，取数据异常....")
								break
							} else {
								glog.Channel("mq").Infoln("是否取消订阅:" + strconv.FormatBool(rabbit.Cancel))
								//收到消息体处理业务
								doErr := rabbit.Handler.Do(rabbit.Topic, protocol, data.Body)
								if doErr == nil {
									_ = data.Ack(false) //false接收到消息确认处理,队列删除数据
								} else {
									_ = data.Nack(false, true)
								}
							}

							//是否取消订阅
							if rabbit.Cancel {
								channel.Close()
								glog.Channel("mq").Infoln("退出监听")
								break
							}
						}
					} else {
						glog.Channel("mq").Infoln("Consume 消费端=>>>获取通道错误..." + err.Error())
					}
				} else {
					glog.Channel("mq").Infoln("Channel 消费端=>>>获取通道错误..." + err.Error())
				}
			}

			//是否取消订阅
			if rabbit.Cancel {
				glog.Channel("mq").Infoln("退出监听...")
				break
			}

			glog.Channel("mq").Infoln("消费端=>>>获取通道错误...")
			time.Sleep(5 * time.Second)
		}
	}()
}

// Do 收到消息体处理业务
func (h *Handler) Do(_ string, protocol string, data []byte) error {
	values := []reflect.Value{}

	ctx := context.Background()
	ctxV := reflect.ValueOf(ctx)
	values = append(values, ctxV)

	reqPtr := reflect.New(h.ReqType.Elem())
	//req := reqPtr.Elem()
	if protocol == "json" {
		_ = json.Unmarshal(data, reqPtr.Interface())
	} else {
		_ = proto.Unmarshal(data, reqPtr.Interface().(proto.Message))
	}

	values = append(values, reqPtr)

	result := h.Method.Call(values)
	errIF := result[0].Interface()
	if errIF == nil {
		return nil
	} else {
		return errIF.(error)
	}
}

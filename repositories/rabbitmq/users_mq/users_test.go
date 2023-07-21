package users_mq_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/repositories/rabbitmq/users_mq"
	"github.com/mrminglang/tools/uuids"
	genid "github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	boot.Boot()
	m.Run()
}

func TestSendUserToMQ(t *testing.T) {
	topic := "mqTopic1"
	user := &users_mq.UserInfo{
		UserID:   uuids.GetRandowUUID(),
		UserName: genid.NewGeneratorData().Name,
	}
	ok := users_mq.SendUserToMQ(topic, user)
	assert.True(t, ok)

	ok1 := users_mq.RecvUserFromMQ(topic)
	assert.True(t, ok1)

	//time.Sleep(5 * time.Second)
	//boot.Destroy()
}

func TestRecvUserFromMQ(t *testing.T) {
	topic := "mqTopic"
	ok := users_mq.RecvUserFromMQ(topic)
	assert.True(t, ok)

	time.Sleep(10 * time.Second)

	boot.Destroy()
}

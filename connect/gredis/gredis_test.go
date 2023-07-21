package gredis_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/connect/gredis"
	"github.com/mrminglang/tools/dumps"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMain(m *testing.M) {
	boot.Boot()
	m.Run()
}

func TestRedisBase_Set(t *testing.T) {
	key := "key123"
	val := "value123"
	err := gredis.GetRedisCon().Set(key, val, "")
	assert.Nil(t, err)
	boot.Destroy()
}

func TestRedisBase_Get(t *testing.T) {
	key := "key123"
	val, err := gredis.GetRedisCon().Get(key, "")
	assert.Nil(t, err)
	dumps.Dump(val)
	boot.Destroy()
}

func TestRedisBase_Del(t *testing.T) {
	key := "key123"
	ok := gredis.GetRedisCon().Del(key, "")
	assert.True(t, ok)

	boot.Destroy()
}

func TestRedisBase_GetAllKey(t *testing.T) {
	keys, err := gredis.GetRedisCon().GetAllKey("")
	assert.Nil(t, err)
	dumps.Dump(keys)
	boot.Destroy()
}

func TestRedisBase_SetLock(t *testing.T) {
	ok := gredis.GetRedisCon().SetLock("key100", "value100", 60, "")
	assert.True(t, ok)

	val, err := gredis.GetRedisCon().Get("key100", "")
	assert.Nil(t, err)
	dumps.Dump(val)

	boot.Destroy()
}
func TestRedisBase_ReleaseLock(t *testing.T) {
	ok := gredis.GetRedisCon().ReleaseLock("key100", "value100", "")
	assert.True(t, ok)
	boot.Destroy()
}

package users_repository_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/repositories/models"
	"github.com/mrminglang/go-rigger/repositories/postgresql/users_repository"
	"github.com/mrminglang/tools/dumps"
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

func TestUsers_IsExistUserName(t *testing.T) {
	userName := "张三三"
	ok := users_repository.NewUsers().IsExistUserName(userName)
	assert.True(t, ok)

	boot.Destroy()
}

func TestUsers_CreateUsers(t *testing.T) {
	user := &models.Users{
		UserID:    uuids.GetRandowUUID(),
		UserName:  "张三",
		Phone:     genid.NewGeneratorData().PhoneNum,
		CreatedAt: time.Now(),
	}
	err := users_repository.NewUsers().CreateUsers(user)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_BatchCreateUsers(t *testing.T) {
	users := make([]*models.Users, 0)
	for i := 0; i < 100; i++ {
		user := &models.Users{
			UserID:    uuids.GetRandowUUID(),
			UserName:  genid.NewGeneratorData().Name,
			Phone:     genid.NewGeneratorData().PhoneNum,
			CreatedAt: time.Now(),
		}
		users = append(users, user)
	}

	err := users_repository.NewUsers().BatchCreateUsers(users)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_UpdateUsers(t *testing.T) {
	user := &models.Users{
		UserID:    "2a6b4103e25041e9af56875d6c6d1927",
		UserName:  "张三三",
		UpdatedAt: time.Now(),
	}
	err := users_repository.NewUsers().UpdateUsers(user)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_DeleteUsers(t *testing.T) {
	userID := "b650a100cc2c42f38e1d9418d7c072bf"
	err := users_repository.NewUsers().DeleteUsers(userID)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_QueryUsers(t *testing.T) {
	whereMaps := map[string]string{
		//"user_name": "张三",
		"order": "created_at DESC",
	}

	total, users, err := users_repository.NewUsers().QueryUsers(0, 10, whereMaps)
	assert.Nil(t, err)
	dumps.Dump(total)
	dumps.Dump(users)

	boot.Destroy()
}

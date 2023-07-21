package users_repository_test

import (
	"github.com/mrminglang/go-rigger/boot"
	"github.com/mrminglang/go-rigger/repositories/models"
	"github.com/mrminglang/go-rigger/repositories/mysql/users_repository"
	"github.com/mrminglang/tools/dumps"
	"github.com/mrminglang/tools/uuids"
	"github.com/srlemon/gen-id"
	"github.com/stretchr/testify/assert"
	"github.com/zhan3333/glog"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	//boot.Boot()
	m.Run()
}

func TestUsers_IsExistUserName(t *testing.T) {
	userName := "张三三"
	dumps.Dump(userName)
	//ok := users_repository.NewUsers().IsExistUserName(userName)
	//assert.True(t, ok)

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

	for i := 0; i < 100000; i++ {
		user := &models.Users{
			UserID:    uuids.GetRandowUUID(),
			UserName:  genid.NewGeneratorData().Name,
			Phone:     genid.NewGeneratorData().PhoneNum,
			CreatedAt: time.Now(),
		}
		users = append(users, user)
	}
	start := time.Now()
	err := users_repository.NewUsers().BatchCreateUsers(users)
	assert.Nil(t, err)
	end := time.Now()
	glog.Def().Infoln("start time::", start)
	glog.Def().Infoln("end time::", end)
	boot.Destroy()
}

func TestUsers_UpdateUsers(t *testing.T) {
	user := &models.Users{
		UserID:    "bf88ef7589984a7db15349aa01925d8d",
		UserName:  "张三三",
		UpdatedAt: time.Now(),
	}
	err := users_repository.NewUsers().UpdateUsers(user)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_DeleteUsers(t *testing.T) {
	userID := "cf77cc7967a648da88ecdfdb1d36e046"
	err := users_repository.NewUsers().DeleteUsers(userID)
	assert.Nil(t, err)

	boot.Destroy()
}

func TestUsers_QueryUsers(t *testing.T) {
	whereMaps := map[string]string{
		//"user_name": "张三",
		"order": "created_at ASC",
	}

	total, users, err := users_repository.NewUsers().QueryUsers(0, 10, whereMaps)
	assert.Nil(t, err)
	dumps.Dump(total)
	dumps.Dump(users)

	boot.Destroy()
}

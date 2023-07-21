package users_repository

import (
	"github.com/mrminglang/go-rigger/connect/gmysql"
	"github.com/mrminglang/go-rigger/repositories/models"
	"gorm.io/gorm"
)

type user struct{}

func NewUsers() *user {
	return &user{}
}

// IsExistUserName 判断用户名称是否存在
func (u *user) IsExistUserName(userName string) bool {
	if gmysql.DbMySQLGormConn.Select("user_id").Where("user_name = ?", userName).First(&models.Users{}).RowsAffected <= 0 {
		return false
	}

	return true
}

// CreateUsers 创建用户
func (u *user) CreateUsers(users *models.Users) (err error) {
	err = gmysql.DbMySQLGormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&users).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

// BatchCreateUsers 批量分块创建用户
func (u *user) BatchCreateUsers(users []*models.Users) (err error) {
	err = gmysql.DbMySQLGormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.CreateInBatches(users, gmysql.DbBatchSize).Error; err != nil {
			return err
		}
		return nil
	})

	return
}

// UpdateUsers 更新用户数据
func (u *user) UpdateUsers(users *models.Users) (err error) {
	err = gmysql.DbMySQLGormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(&users).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

// DeleteUsers 删除用户数据
func (u *user) DeleteUsers(userID string) (err error) {
	err = gmysql.DbMySQLGormConn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&models.Users{}, "user_id = ?", userID).Error; err != nil {
			return err
		}
		return nil
	})
	return
}

// QueryUsers 分页查询用户
func (u *user) QueryUsers(
	beginIndex int,
	count int,
	whereMaps map[string]string,
) (total int64, users []*models.Users, err error) {
	query := gmysql.DbMySQLGormConn.Model(&users)

	// 用户名称
	if whereMaps["user_name"] != "" {
		query = query.Where("user_name LIKE ?", "%"+whereMaps["user_name"]+"%")
	}

	// 手机号
	if whereMaps["phone"] != "" {
		query = query.Where("phone LIKE ?", "%"+whereMaps["phone"]+"%")
	}

	// 排序
	if whereMaps["order"] != "" {
		query = query.Order(whereMaps["order"])
	}

	err = query.Count(&total).Offset(beginIndex).Limit(count).Find(&users).Error
	if err != nil {
		return
	}

	return
}

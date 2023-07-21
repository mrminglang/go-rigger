package models

import "time"

type Users struct {
	UserID    string    `gorm:"column:user_id;type:char(50);primary_key;comment:'用户ID'" json:"userID"`
	UserName  string    `gorm:"column:user_name;type:char(100);default:'';comment:'用户名称'" json:"userName"`
	Phone     string    `gorm:"column:phone;type:char(11);index:phone;comment:'手机号'" json:"phone"`
	CreatedAt time.Time `gorm:"column:created_at;comment:'创建时间'" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;comment:'更新时间'" json:"updatedAt"`
}

// TableName sets the insert table name for this struct type
func (r Users) TableName() string {
	return "users"
}

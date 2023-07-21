package models

import "time"

type ResInfo struct {
	ResID    int64     `gorm:"column:RES_ID;type:int(10);primaryKey" json:"RES_ID"`
	DeclDate time.Time `gorm:"column:DECL_DATE;type:datetime;index:IN_RES_INFO_WRI" json:"DECL_DATE"`
	ResTitle string    `gorm:"column:RES_TITLE;type:char(300)" json:"RES_TITLE"`
}

// TableName sets the insert table name for this struct type
func (r ResInfo) TableName() string {
	return "UPCENTER.RES_INFO"
}

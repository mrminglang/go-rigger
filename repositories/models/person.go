package models

import "time"

type Person struct {
	Name     string    `bson:"name"`
	Age      int       `bson:"age"`
	Email    string    `bson:"email"`
	CreateAt time.Time `bson:"create_at"`
}

// TableName sets the insert table name for this struct type
func (r Person) TableName() string {
	return "person"
}

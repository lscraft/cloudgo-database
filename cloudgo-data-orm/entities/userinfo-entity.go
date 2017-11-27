package entities

import (
	"time"
)

// UserInfo .
type UserInfo struct {
	UID        int `xorm:"notnull autoincr pk"` //语义标签
	UserName   string
	DepartName string
	CreateAt   time.Time `xorm:"created"`
}

// NewUserInfo .
func NewUserInfo(u UserInfo) *UserInfo {
	if len(u.UserName) == 0 {
		panic("UserName shold not null!")
	}
	return &u
}

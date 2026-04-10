package mysqlServer

import "gorm.io/gorm"

// ---------- 模型定义（新增模型在此文件追加） ----------

// User 用户表模型
type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);not null"            json:"-"`
	Nickname string `gorm:"column:nickname;type:varchar(64)"                      json:"nickname"`
	Email    string `gorm:"column:email;type:varchar(128)"                        json:"email"`
	Status   int    `gorm:"column:status;type:tinyint;default:1"                  json:"status"` // 1=正常 0=禁用
}

func (User) TableName() string { return "user" }

// ---------- 模型注册（新增模型在 init 中追加） ----------

func init() {
	RegisterModel(
		&User{},
	)
}

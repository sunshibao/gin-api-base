package model

import (
	"gin-api/server"

	"gorm.io/gorm"
)

// User 用户表模型
type User struct {
	gorm.Model
	Username string `gorm:"column:username;type:varchar(64);uniqueIndex;not null" json:"username"`
	Password string `gorm:"column:password;type:varchar(255);not null"            json:"-"`
	Nickname string `gorm:"column:nickname;type:varchar(64)"                      json:"nickname"`
	Email    string `gorm:"column:email;type:varchar(128)"                        json:"email"`
	Status   int    `gorm:"column:status;type:tinyint;default:1"                  json:"status"` // 1=正常 0=禁用
}

func (User) TableName() string {
	return "user"
}

// AutoMigrate 自动建表
func AutoMigrate() error {
	return server.GetDB().AutoMigrate(&User{})
}

// --- 数据访问方法 ---

// CreateUser 创建用户
func CreateUser(user *User) error {
	return server.GetDB().Create(user).Error
}

// GetUserByID 根据 ID 查询用户
func GetUserByID(id uint) (*User, error) {
	var user User
	err := server.GetDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名查询
func GetUserByUsername(username string) (*User, error) {
	var user User
	err := server.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户
func UpdateUser(user *User) error {
	return server.GetDB().Save(user).Error
}

// DeleteUser 删除用户（软删除）
func DeleteUser(id uint) error {
	return server.GetDB().Delete(&User{}, id).Error
}

// ListUsers 分页查询用户列表
func ListUsers(page, size int) ([]User, int64, error) {
	var users []User
	var total int64

	db := server.GetDB().Model(&User{})
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&users).Error
	return users, total, err
}

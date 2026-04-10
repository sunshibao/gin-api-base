package model

import mysqlServer "gin-api/server/mysql"

// --- 数据访问方法 ---

// CreateUser 创建用户
func CreateUser(user *mysqlServer.User) error {
	return mysqlServer.GetDB().Create(user).Error
}

// GetUserByID 根据 ID 查询用户
func GetUserByID(id uint) (*mysqlServer.User, error) {
	var user mysqlServer.User
	err := mysqlServer.GetDB().First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetUserByUsername 根据用户名查询
func GetUserByUsername(username string) (*mysqlServer.User, error) {
	var user mysqlServer.User
	err := mysqlServer.GetDB().Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser 更新用户
func UpdateUser(user *mysqlServer.User) error {
	return mysqlServer.GetDB().Save(user).Error
}

// DeleteUser 删除用户（软删除）
func DeleteUser(id uint) error {
	return mysqlServer.GetDB().Delete(&mysqlServer.User{}, id).Error
}

// ListUsers 分页查询用户列表
func ListUsers(page, size int) ([]mysqlServer.User, int64, error) {
	var users []mysqlServer.User
	var total int64

	db := mysqlServer.GetDB().Model(&mysqlServer.User{})
	db.Count(&total)

	offset := (page - 1) * size
	err := db.Offset(offset).Limit(size).Order("id DESC").Find(&users).Error
	return users, total, err
}

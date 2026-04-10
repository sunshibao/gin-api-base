package service

import (
	"crypto/sha256"
	"errors"
	"fmt"

	"gin-api/model"
	"gin-api/pkg/jwtutil"
	mysqlServer "gin-api/server/mysql"

	"gorm.io/gorm"
)

// --- 请求/响应结构体 ---

type RegisterReq struct {
	Username string `json:"username" binding:"required,min=3,max=32"`
	Password string `json:"password" binding:"required,min=6,max=64"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type LoginReq struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResp struct {
	Token    string `json:"token"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type UpdateUserReq struct {
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
}

type UserInfo struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Email    string `json:"email"`
	Status   int    `json:"status"`
}

// --- 业务逻辑 ---

// Register 用户注册
func Register(req *RegisterReq) error {
	// 检查用户名是否已存在
	_, err := model.GetUserByUsername(req.Username)
	if err == nil {
		return errors.New("用户名已存在")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("查询用户失败: %w", err)
	}

	user := &mysqlServer.User{
		Username: req.Username,
		Password: hashPassword(req.Password),
		Nickname: req.Nickname,
		Email:    req.Email,
		Status:   1,
	}
	return model.CreateUser(user)
}

// Login 用户登录
func Login(req *LoginReq) (*LoginResp, error) {
	user, err := model.GetUserByUsername(req.Username)
	if err != nil {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Password != hashPassword(req.Password) {
		return nil, errors.New("用户名或密码错误")
	}

	if user.Status != 1 {
		return nil, errors.New("账号已被禁用")
	}

	token, err := jwtutil.GenerateToken(user.ID, user.Username)
	if err != nil {
		return nil, fmt.Errorf("生成 token 失败: %w", err)
	}

	return &LoginResp{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

// GetUser 获取用户信息
func GetUser(id uint) (*UserInfo, error) {
	user, err := model.GetUserByID(id)
	if err != nil {
		return nil, errors.New("用户不存在")
	}
	return toUserInfo(user), nil
}

// UpdateUser 更新用户信息
func UpdateUser(id uint, req *UpdateUserReq) error {
	user, err := model.GetUserByID(id)
	if err != nil {
		return errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	return model.UpdateUser(user)
}

// DeleteUser 删除用户
func DeleteUser(id uint) error {
	return model.DeleteUser(id)
}

// ListUsers 用户列表
func ListUsers(page, size int) ([]UserInfo, int64, error) {
	if page < 1 {
		page = 1
	}
	if size < 1 || size > 100 {
		size = 10
	}

	users, total, err := model.ListUsers(page, size)
	if err != nil {
		return nil, 0, err
	}

	list := make([]UserInfo, len(users))
	for i, u := range users {
		list[i] = *toUserInfo(&u)
	}
	return list, total, nil
}

// --- 工具函数 ---

func hashPassword(password string) string {
	h := sha256.Sum256([]byte(password))
	return fmt.Sprintf("%x", h)
}

func toUserInfo(u *mysqlServer.User) *UserInfo {
	return &UserInfo{
		ID:       u.ID,
		Username: u.Username,
		Nickname: u.Nickname,
		Email:    u.Email,
		Status:   u.Status,
	}
}

package handler

import (
	"strconv"

	"gin-api/middleware"
	"gin-api/pkg/resp"
	"gin-api/service"

	"github.com/gin-gonic/gin"
)

// Register 用户注册
// POST /api/register
func Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := service.Register(&req); err != nil {
		resp.Fail(c, 1001, err.Error())
		return
	}
	resp.OK(c, "注册成功")
}

// Login 用户登录
// POST /api/login
func Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	data, err := service.Login(&req)
	if err != nil {
		resp.Fail(c, 1002, err.Error())
		return
	}
	resp.OK(c, data)
}

// GetUser 获取用户信息
// GET /api/user/:id
func GetUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, 400, "无效的用户 ID")
		return
	}

	data, err := service.GetUser(uint(id))
	if err != nil {
		resp.Fail(c, 404, err.Error())
		return
	}
	resp.OK(c, data)
}

// GetCurrentUser 获取当前登录用户信息
// GET /api/user/me
func GetCurrentUser(c *gin.Context) {
	userID := middleware.GetUserID(c)
	data, err := service.GetUser(userID)
	if err != nil {
		resp.Fail(c, 404, err.Error())
		return
	}
	resp.OK(c, data)
}

// UpdateUser 更新用户信息
// PUT /api/user/:id
func UpdateUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, 400, "无效的用户 ID")
		return
	}

	var req service.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		resp.Fail(c, 400, "参数错误: "+err.Error())
		return
	}

	if err := service.UpdateUser(uint(id), &req); err != nil {
		resp.Fail(c, 1003, err.Error())
		return
	}
	resp.OK(c, "更新成功")
}

// DeleteUser 删除用户
// DELETE /api/user/:id
func DeleteUser(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		resp.Fail(c, 400, "无效的用户 ID")
		return
	}

	if err := service.DeleteUser(uint(id)); err != nil {
		resp.Fail(c, 1004, err.Error())
		return
	}
	resp.OK(c, "删除成功")
}

// ListUsers 用户列表
// GET /api/users?page=1&size=10
func ListUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "10"))

	list, total, err := service.ListUsers(page, size)
	if err != nil {
		resp.Fail(c, 500, "查询失败")
		return
	}
	resp.OK(c, resp.Page{List: list, Total: total, Page: page, Size: size})
}

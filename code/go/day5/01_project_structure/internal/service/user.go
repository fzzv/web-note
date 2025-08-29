// Package service 提供业务逻辑服务
package service

import (
	"errors"
	"time"
)

// 预定义错误
var (
	ErrUserNotFound = errors.New("用户不存在")
	ErrInvalidInput = errors.New("无效输入")
)

// User 用户模型
type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UserService 用户服务
type UserService struct {
	// 在实际应用中，这里会有数据库连接或仓库层
	users []User
}

// NewUserService 创建新的用户服务
func NewUserService() *UserService {
	// 模拟一些用户数据
	now := time.Now()
	users := []User{
		{
			ID:        1,
			Name:      "张三",
			Email:     "zhangsan@example.com",
			CreatedAt: now.Add(-24 * time.Hour),
			UpdatedAt: now.Add(-24 * time.Hour),
		},
		{
			ID:        2,
			Name:      "李四",
			Email:     "lisi@example.com",
			CreatedAt: now.Add(-12 * time.Hour),
			UpdatedAt: now.Add(-12 * time.Hour),
		},
		{
			ID:        3,
			Name:      "王五",
			Email:     "wangwu@example.com",
			CreatedAt: now.Add(-6 * time.Hour),
			UpdatedAt: now.Add(-6 * time.Hour),
		},
	}

	return &UserService{
		users: users,
	}
}

// GetAllUsers 获取所有用户
func (s *UserService) GetAllUsers() ([]User, error) {
	// 在实际应用中，这里会查询数据库
	// 这里只是返回模拟数据
	return s.users, nil
}

// GetUserByID 根据 ID 获取用户
func (s *UserService) GetUserByID(id int) (*User, error) {
	// 在实际应用中，这里会查询数据库
	for _, user := range s.users {
		if user.ID == id {
			// 返回副本，避免外部修改
			userCopy := user
			return &userCopy, nil
		}
	}

	return nil, ErrUserNotFound
}

// CreateUser 创建新用户
func (s *UserService) CreateUser(name, email string) (*User, error) {
	// 输入验证
	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	// 检查邮箱是否已存在
	for _, user := range s.users {
		if user.Email == email {
			return nil, errors.New("邮箱已存在")
		}
	}

	// 创建新用户
	now := time.Now()
	newUser := User{
		ID:        len(s.users) + 1, // 简单的 ID 生成策略
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}

	// 添加到用户列表
	s.users = append(s.users, newUser)

	return &newUser, nil
}

// UpdateUser 更新用户信息
func (s *UserService) UpdateUser(id int, name, email string) (*User, error) {
	// 输入验证
	if name == "" || email == "" {
		return nil, ErrInvalidInput
	}

	// 查找用户
	for i, user := range s.users {
		if user.ID == id {
			// 更新用户信息
			s.users[i].Name = name
			s.users[i].Email = email
			s.users[i].UpdatedAt = time.Now()

			// 返回更新后的用户
			userCopy := s.users[i]
			return &userCopy, nil
		}
	}

	return nil, ErrUserNotFound
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	// 查找用户
	for i, user := range s.users {
		if user.ID == id {
			// 删除用户
			s.users = append(s.users[:i], s.users[i+1:]...)
			return nil
		}
	}

	return ErrUserNotFound
}

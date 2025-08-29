// Package handler 提供 HTTP 请求处理器
package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"day5/01_project_structure/internal/service"
	"day5/01_project_structure/pkg/logger"
)

// UserHandler 用户相关的 HTTP 处理器
type UserHandler struct {
	userService *service.UserService
	logger      logger.Logger
}

// NewUserHandler 创建新的用户处理器
func NewUserHandler(userService *service.UserService, logger logger.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUsers 获取用户列表
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	h.logger.Info("获取用户列表请求")

	users, err := h.userService.GetAllUsers()
	if err != nil {
		h.logger.Error("获取用户列表失败", "error", err)
		h.writeError(w, http.StatusInternalServerError, "内部服务器错误")
		return
	}

	h.writeJSON(w, http.StatusOK, map[string]interface{}{
		"users": users,
		"count": len(users),
	})
}

// GetUser 根据 ID 获取单个用户
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.writeError(w, http.StatusMethodNotAllowed, "方法不允许")
		return
	}

	// 从 URL 路径中提取用户 ID
	// URL 格式: /users/{id}
	path := strings.TrimPrefix(r.URL.Path, "/users/")
	if path == "" {
		h.writeError(w, http.StatusBadRequest, "缺少用户 ID")
		return
	}

	userID, err := strconv.Atoi(path)
	if err != nil {
		h.writeError(w, http.StatusBadRequest, "无效的用户 ID")
		return
	}

	h.logger.Info("获取用户请求", "userID", userID)

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		if err == service.ErrUserNotFound {
			h.writeError(w, http.StatusNotFound, "用户不存在")
			return
		}
		h.logger.Error("获取用户失败", "userID", userID, "error", err)
		h.writeError(w, http.StatusInternalServerError, "内部服务器错误")
		return
	}

	h.writeJSON(w, http.StatusOK, user)
}

// writeJSON 写入 JSON 响应
func (h *UserHandler) writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("JSON 编码失败", "error", err)
	}
}

// writeError 写入错误响应
func (h *UserHandler) writeError(w http.ResponseWriter, status int, message string) {
	h.writeJSON(w, status, map[string]string{
		"error": message,
	})
}

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// ValidateTokenRequest 验证令牌请求结构
type ValidateTokenRequest struct {
	Token string `json:"token"`
}

// UserInfoResponse 用户信息响应结构
type UserInfoResponse struct {
	User struct {
		ID       uint   `json:"id"`
		Username string `json:"username"`
		Email    string `json:"email"`
	} `json:"user"`
}

// ValidateToken 向认证服务请求验证JWT令牌
func ValidateToken(token, authServiceURL string) (*UserInfoResponse, error) {
	// 构建认证服务的验证端点URL
	validateURL := fmt.Sprintf("http://%s/api/auth/verify", authServiceURL)

	// 创建验证请求
	reqBody := ValidateTokenRequest{
		Token: token,
	}
	reqJSON, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	// 发送请求到认证服务
	resp, err := http.Post(validateURL, "application/json", bytes.NewBuffer(reqJSON))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("validation failed with status code: %d", resp.StatusCode)
	}

	// 解析响应
	var userInfo UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}

// GetUserFromAuthService 从认证服务获取用户信息（备用方案：直接调用认证服务的/user端点）
func GetUserFromAuthService(token, authServiceURL string) (*UserInfoResponse, error) {
	// 构建请求URL
	userURL := fmt.Sprintf("http://%s/api/auth/me", authServiceURL)

	// 创建请求
	req, err := http.NewRequest("GET", userURL, nil)
	if err != nil {
		return nil, err
	}

	// 设置Authorization头
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info, status code: %d", resp.StatusCode)
	}

	// 解析响应
	var userInfo UserInfoResponse
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
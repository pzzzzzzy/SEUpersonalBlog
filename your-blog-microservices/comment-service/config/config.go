package config

import (
	"os"
)

type Config struct {
	Port        string
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	ConsulAddr  string
	ServiceName string
	ServiceID   string
	ServicePort string
	AuthService string
}

// LoadConfig 从环境变量加载配置
func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8083"),
		DBHost:      getEnv("DB_HOST", "mysql"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "blog"),
		ConsulAddr:  getEnv("CONSUL_ADDR", "consul:8500"),
		ServiceName: getEnv("SERVICE_NAME", "comment-service"),
		ServiceID:   getEnv("SERVICE_ID", "comment-service-1"),
		ServicePort: getEnv("SERVICE_PORT", "8083"),
		AuthService: getEnv("AUTH_SERVICE", "auth-service"),
	}
}

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
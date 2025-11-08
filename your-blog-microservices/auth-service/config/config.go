package config

import "os"

type Config struct {
	Port       string
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	JWTSecret  string
	ConsulAddr string
	ServiceName string
	ServiceID   string
	ServicePort string
}

func LoadConfig() *Config {
	return &Config{
		Port:        getEnv("PORT", "8081"),
		DBHost:      getEnv("DB_HOST", "mysql"),
		DBPort:      getEnv("DB_PORT", "3306"),
		DBUser:      getEnv("DB_USER", "root"),
		DBPassword:  getEnv("DB_PASSWORD", "password"),
		DBName:      getEnv("DB_NAME", "blog"),
		JWTSecret:   getEnv("JWT_SECRET", "your-secret-key"),
		ConsulAddr:  getEnv("CONSUL_ADDR", "consul:8500"),
		ServiceName: "auth-service",
		ServiceID:   getEnv("SERVICE_ID", "auth-service-1"),
		ServicePort: getEnv("PORT", "8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
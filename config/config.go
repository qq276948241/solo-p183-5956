package config

import "os"

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBCharset  string
	ServerPort string
	AuthToken  string
}

func Load() *Config {
	return &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),
		DBPort:     getEnv("DB_PORT", "3306"),
		DBUser:     getEnv("DB_USER", "clinic"),
		DBPassword: getEnv("DB_PASSWORD", "clinic123"),
		DBName:     getEnv("DB_NAME", "clinic"),
		DBCharset:  getEnv("DB_CHARSET", "utf8mb4"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
		AuthToken:  getEnv("AUTH_TOKEN", "clinic-secret-token-2024"),
	}
}

func (c *Config) DSN() string {
	return c.DBUser + ":" + c.DBPassword + "@tcp(" + c.DBHost + ":" + c.DBPort + ")/" + c.DBName + "?charset=" + c.DBCharset + "&parseTime=True&loc=Local"
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

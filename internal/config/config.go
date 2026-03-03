package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerAddr      string
	ShutdownTimeout time.Duration
	LogInterval     time.Duration

	MongoURI string
	MongoDB  string

	RedisAddr     string
	RedisPassword string
	RedisDB       int
	LogTTLSeconds int
}

func Default() Config {
	return Config{
		ServerAddr:      getEnv("SERVER_ADDR", ":8081"),
		ShutdownTimeout: 5 * time.Second,
		LogInterval:     200 * time.Millisecond,

		MongoURI: getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDB:  getEnv("MONGO_DB", "gostudy"),

		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvInt("REDIS_DB", 0),
		LogTTLSeconds: getEnvInt("LOG_TTL_SECONDS", 86400),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

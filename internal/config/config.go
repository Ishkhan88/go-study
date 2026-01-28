package config

import "time"

type Config struct {
	ServerAddr      string
	ShutdownTimeout time.Duration
	LogInterval     time.Duration
}

func Default() Config {
	return Config{
		ServerAddr:      ":8081",
		ShutdownTimeout: 5 * time.Second,
		LogInterval:     200 * time.Millisecond,
	}
}

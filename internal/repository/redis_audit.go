package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client
var auditTTL time.Duration

// InitRedis нужно вызвать при старте сервиса (рядом с InitMongo).
func InitRedis(ctx context.Context, addr, password string, db int, ttlSeconds int) error {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	pctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := redisClient.Ping(pctx).Err(); err != nil {
		return fmt.Errorf("redis ping: %w", err)
	}

	if ttlSeconds <= 0 {
		ttlSeconds = 86400
	}
	auditTTL = time.Duration(ttlSeconds) * time.Second
	return nil
}

func CloseRedis(ctx context.Context) error {
	if redisClient == nil {
		return nil
	}
	return redisClient.Close()
}

func mustRedis() {
	if redisClient == nil {
		panic("redis is not initialized: call repository.InitRedis() on startup")
	}
}

// AuditEvent пишет событие в Redis Stream и продлевает TTL ключа стрима.
func AuditEvent(entity string, entityID int, action string, payload any) {
	mustRedis()

	// stream key: audit:<entity>:<id>
	key := fmt.Sprintf("audit:%s:%d", entity, entityID)

	b, _ := json.Marshal(payload)

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	// XADD
	_ = redisClient.XAdd(ctx, &redis.XAddArgs{
		Stream: key,
		Values: map[string]any{
			"ts":     time.Now().UTC().Format(time.RFC3339Nano),
			"action": action, // create/update/delete
			"json":   string(b),
		},
	}).Err()

	// TTL на весь стрим (обновляем, чтобы история жила N секунд с последнего изменения)
	_ = redisClient.Expire(ctx, key, auditTTL).Err()
}

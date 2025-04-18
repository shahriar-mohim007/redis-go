package main

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()
var rdb *redis.Client

func InitRedis() {
	rdb = redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
		Password: getEnv("REDIS_PASSWORD", "redis1234"),
	})
}

func UpdatePlayerScore(playerID string, score int) error {
	return rdb.ZIncrBy(ctx, "leaderboard", float64(score), playerID).Err()
}

func GetTopPlayers(n int64) ([]redis.Z, error) {
	return rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, n-1).Result()
}

func PeriodicSyncWorker() {
	ticker := time.NewTicker(1 * time.Second)
	for range ticker.C {
		_ = SyncLeaderboardToDB()
	}
}

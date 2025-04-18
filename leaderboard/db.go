package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5433"),
		getEnv("DB_USER", "postgres"),
		getEnv("DB_PASSWORD", "postgres"),
		getEnv("DB_NAME", "leaderboard_db"),
	)
	db, err = sql.Open("postgres", dsn)
	if err != nil {
		log.Fatal("Error connecting to DB:", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		log.Fatal("DB ping failed:", err)
	}
}

func SyncLeaderboardToDB() error {
	leaderboard := rdb.ZRevRangeWithScores(ctx, "leaderboard", 0, -1)
	results, err := leaderboard.Result()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	for _, item := range results {
		player := item.Member.(string)
		score := int(item.Score)

		_, err := tx.Exec(
			"INSERT INTO leaderboard (player_id, score, updated_at) VALUES ($1, $2, NOW()) ON CONFLICT (player_id) DO UPDATE SET score=$2, updated_at=NOW()",
			player, score,
		)
		if err != nil {
			return err
		}

		log.Printf("Inserted/Updated player_id=%s with score=%d into DB", player, score)
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	if err := rdb.Del(ctx, "leaderboard").Err(); err != nil {
		return err
	}

	log.Println("✅ Leaderboard sync completed successfully.")

	return nil
}

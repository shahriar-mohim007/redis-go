package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	_ "github.com/lib/pq"
)

const (
	POSTGRES_USER     = "postgres"
	POSTGRES_PASSWORD = "postgres"
	POSTGRES_DB       = "url_shortener"
	POSTGRES_HOST     = "localhost"
	POSTGRES_PORT     = "5433"
)

var ctx = context.Background()

func newRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "redis1234",
		DB:       0,
	})
}

func newPostgresClient() (*sql.DB, error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		POSTGRES_USER, POSTGRES_PASSWORD, POSTGRES_DB, POSTGRES_HOST, POSTGRES_PORT)
	return sql.Open("postgres", connStr)
}

func getOriginalURLFromDB(db *sql.DB, shortURL string) (string, error) {
	var originalURL string
	err := db.QueryRow("SELECT original_url FROM urls WHERE short_url = $1", shortURL).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil
		}
		return "", err
	}
	return originalURL, nil
}

func getURLWithCache(redisClient *redis.Client, db *sql.DB, shortURL string) (string, error) {

	val, err := redisClient.Get(ctx, shortURL).Result()
	if errors.Is(err, redis.Nil) {
		fmt.Println("Cache miss for URL:", shortURL)

		val, err = getOriginalURLFromDB(db, shortURL)
		if err != nil {
			return "", err
		}
		if val != "" {
			redisClient.Set(ctx, shortURL, val, 10*time.Minute)
		}
	} else if err != nil {
		return "", err
	} else {
		fmt.Println("Cache hit for URL:", shortURL)
	}
	return val, nil
}

func main() {

	redisClient := newRedisClient()
	postgresClient, err := newPostgresClient()
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL:", err)
	}
	defer postgresClient.Close()

	shortURL := "42v8nK"
	start := time.Now()
	originalURL, err := getURLWithCache(redisClient, postgresClient, shortURL)
	if err != nil {
		log.Fatal("Error fetching URL:", err)
	}
	if originalURL != "" {
		fmt.Printf("Original URL: %s\n", originalURL)
	} else {
		fmt.Println("No URL found for short URL:", shortURL)
	}
	fmt.Printf("First call took %s\n", time.Since(start))

	start = time.Now()
	originalURL, err = getURLWithCache(redisClient, postgresClient, shortURL)
	if err != nil {
		log.Fatal("Error fetching URL:", err)
	}
	fmt.Printf("Original URL: %s\n", originalURL)
	fmt.Printf("Second call took %s\n", time.Since(start))
}

package gibrun

import (
	"context"
	"os"
	"strconv"
	"time"

	"github.com/arielfikru/gibrun"
)

var client *gibrun.Client

func Init() {
	addr := os.Getenv("CACHE_HOST")
	port := os.Getenv("CACHE_PORT")
	pass := os.Getenv("CACHE_PASSWORD")
	dbStr := os.Getenv("CACHE_DB")

	if port == "" {
		port = "6379" // Redis default port
	}

	db, _ := strconv.Atoi(dbStr)

	client = gibrun.New(gibrun.Config{
		Addr:     addr + ":" + port,
		Password: pass,
		DB:       db,
	})
}

// Gib (Store) wrapper
func Gib(ctx context.Context, key string, value interface{}) error {
	// Default TTL 24 hours
	return client.Gib(ctx, key).Value(value).TTL(24 * time.Hour).Exec()
}

// Run (Retrieve) wrapper
func Run(ctx context.Context, key string, dest interface{}) (bool, error) {
	return client.Run(ctx, key).Bind(dest)
}

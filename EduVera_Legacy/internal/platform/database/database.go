package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DB adalah instance global koneksi database
var DB *pgxpool.Pool

// Connect menginisialisasi koneksi ke PostgreSQL
func Connect(databaseURL string) *pgxpool.Pool {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	config, err := pgxpool.ParseConfig(databaseURL)
	if err != nil {
		log.Fatalf("Gagal parsing konfigurasi DB: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	// Tes Ping
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Database tidak merespon: %v", err)
	}

	fmt.Println("✅ Terhubung ke Database PostgreSQL!")
	DB = pool // Assign ke variabel global
	return pool
}

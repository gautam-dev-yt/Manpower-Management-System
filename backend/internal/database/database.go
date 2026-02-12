package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"manpower-backend/internal/config"
)

type Service interface {
	Health() map[string]string
	Close() error
	GetPool() *pgxpool.Pool
}

type service struct {
	pool *pgxpool.Pool
}

var dbInstance *service

func New(cfg *config.DBConfig) Service {
	// Reuse existing instance if available (singleton pattern)
	if dbInstance != nil {
		return dbInstance
	}

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode)

	poolConfig, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("Unable to parse database config: %v", err)
	}

	// Set connection pool settings
	poolConfig.MaxConns = 10
	poolConfig.MinConns = 2
	poolConfig.MaxConnLifetime = time.Hour
	poolConfig.MaxConnIdleTime = 30 * time.Minute

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	// Verify connection
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("Unable to connect to database: %v\nMake sure your password in .env is correct!", err)
	}

	log.Println("Connected to PostgreSQL database successfully!")

	dbInstance = &service{pool: pool}
	return dbInstance
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.pool.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// Get stats
	poolStats := s.pool.Stat()
	stats["open_connections"] = fmt.Sprintf("%d", poolStats.TotalConns())
	stats["in_use"] = fmt.Sprintf("%d", poolStats.AcquiredConns())
	stats["idle"] = fmt.Sprintf("%d", poolStats.IdleConns())

	return stats
}

func (s *service) Close() error {
	log.Println("Closing database connection pool...")
	s.pool.Close()
	return nil
}

func (s *service) GetPool() *pgxpool.Pool {
	return s.pool
}

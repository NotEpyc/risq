package database

import (
	"database/sql"
	"fmt"
	"time"

	"risq_backend/pkg/logger"
	"risq_backend/pkg/migrations"

	_ "github.com/lib/pq"
)

type DB struct {
	conn *sql.DB
}

func New(host, port, user, password, dbname, sslmode string) (*DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	logger.Info("Connecting to database...")
	conn, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}

	// Test the connection
	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Configure connection pool
	conn.SetMaxOpenConns(25)
	conn.SetMaxIdleConns(25)
	conn.SetConnMaxLifetime(5 * time.Minute)

	db := &DB{conn: conn}

	// Run migrations
	migrator := migrations.NewMigrator(conn)
	if err := migrator.RunMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	logger.Info("Database connected and migrated successfully")
	return db, nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func (db *DB) GetConn() *sql.DB {
	return db.conn
}

// Health check
func (db *DB) Ping() error {
	return db.conn.Ping()
}

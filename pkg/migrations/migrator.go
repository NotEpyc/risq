package migrations

import (
	"database/sql"
	"fmt"

	"risq_backend/pkg/logger"
)

// Migration represents a database migration
type Migration struct {
	ID   int
	Name string
	Up   func(*sql.DB) error
	Down func(*sql.DB) error
}

// Migrator handles database migrations
type Migrator struct {
	db         *sql.DB
	migrations []Migration
}

// NewMigrator creates a new migrator instance
func NewMigrator(db *sql.DB) *Migrator {
	migrator := &Migrator{
		db:         db,
		migrations: []Migration{},
	}

	// Register migrations
	migrator.registerMigrations()
	return migrator
}

func (m *Migrator) registerMigrations() {
	// Migration 1: Create all tables from clean schema
	m.migrations = append(m.migrations, Migration{
		ID:   1,
		Name: "Create initial clean schema",
		Up:   m.migration001Up,
		Down: m.migration001Down,
	})
}

// RunMigrations executes all pending migrations
func (m *Migrator) RunMigrations() error {
	logger.Info("Running database migrations...")

	// Create migrations table if it doesn't exist
	err := m.createMigrationsTable()
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations, err := m.getAppliedMigrations()
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Run pending migrations
	for _, migration := range m.migrations {
		if !contains(appliedMigrations, migration.ID) {
			logger.Infof("Applying migration %d: %s", migration.ID, migration.Name)

			err := migration.Up(m.db)
			if err != nil {
				return fmt.Errorf("failed to apply migration %d: %w", migration.ID, err)
			}

			err = m.recordMigration(migration.ID)
			if err != nil {
				return fmt.Errorf("failed to record migration %d: %w", migration.ID, err)
			}

			logger.Infof("Successfully applied migration %d", migration.ID)
		}
	}

	logger.Info("All migrations completed successfully")
	return nil
}

func (m *Migrator) createMigrationsTable() error {
	query := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
		);
	`
	_, err := m.db.Exec(query)
	return err
}

func (m *Migrator) getAppliedMigrations() ([]int, error) {
	query := "SELECT version FROM schema_migrations ORDER BY version"
	rows, err := m.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var migrations []int
	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return nil, err
		}
		migrations = append(migrations, id)
	}

	return migrations, nil
}

func (m *Migrator) recordMigration(id int) error {
	query := "INSERT INTO schema_migrations (version, applied_at) VALUES ($1, NOW())"
	_, err := m.db.Exec(query, id)
	return err
}

// Migration 001: Create clean schema from SQL file
func (m *Migrator) migration001Up(db *sql.DB) error {
	// Execute the clean schema SQL directly
	schemaSQL := `
-- Drop all tables if they exist (clean slate)
DROP TABLE IF EXISTS risk_evolutions CASCADE;
DROP TABLE IF EXISTS decisions CASCADE;
DROP TABLE IF EXISTS risk_profiles CASCADE;
DROP TABLE IF EXISTS startups CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'founder',
    startup_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Startups table
CREATE TABLE startups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    industry TEXT NOT NULL,
    sector TEXT,
    funding_stage TEXT NOT NULL,
    location TEXT NOT NULL,
    founded_date TIMESTAMP WITH TIME ZONE,
    team_size INTEGER DEFAULT 1,
    website TEXT,
    business_model TEXT,
    target_market TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Risk profiles table
CREATE TABLE risk_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    overall_risk TEXT NOT NULL DEFAULT 'medium',
    market_risk DECIMAL(3,2) DEFAULT 0.5,
    technical_risk DECIMAL(3,2) DEFAULT 0.5,
    financial_risk DECIMAL(3,2) DEFAULT 0.5,
    team_risk DECIMAL(3,2) DEFAULT 0.5,
    regulatory_risk DECIMAL(3,2) DEFAULT 0.5,
    risk_factors TEXT[],
    recommendations TEXT[],
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Decisions table
CREATE TABLE decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    category TEXT NOT NULL,
    urgency TEXT DEFAULT 'medium',
    status TEXT DEFAULT 'speculating',
    context TEXT,
    ai_analysis JSONB,
    final_choice TEXT,
    reasoning TEXT,
    confirmed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Foreign key constraints
ALTER TABLE users ADD CONSTRAINT fk_users_startup_id 
    FOREIGN KEY (startup_id) REFERENCES startups(id) ON DELETE SET NULL;

-- Indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_startup_id ON users(startup_id);
CREATE INDEX idx_startups_id ON startups(id);
CREATE INDEX idx_risk_profiles_startup_id ON risk_profiles(startup_id);
CREATE INDEX idx_decisions_startup_id ON decisions(startup_id);
CREATE INDEX idx_decisions_status ON decisions(status);`

	// Execute the SQL
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("failed to execute schema: %w", err)
	}

	return nil
}

func (m *Migrator) migration001Down(db *sql.DB) error {
	// Drop all tables
	tables := []string{"decisions", "risk_profiles", "startups", "users", "schema_migrations"}

	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}

// Helper function
func contains(slice []int, item int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

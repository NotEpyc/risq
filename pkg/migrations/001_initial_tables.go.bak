package migrations

import (
	"database/sql"
)

func migration001Up(db *sql.DB) error {
	// Create users table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			email TEXT NOT NULL UNIQUE,
			name TEXT NOT NULL,
			password TEXT NOT NULL,
			role TEXT DEFAULT 'founder',
			startup_id UUID,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create startups table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS startups (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name TEXT NOT NULL,
			description TEXT,
			industry TEXT NOT NULL,
			funding_stage TEXT NOT NULL,
			location TEXT NOT NULL,
			founded_date TIMESTAMPTZ,
			team_size INTEGER DEFAULT 1,
			website TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create decisions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS decisions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			startup_id UUID NOT NULL,
			description TEXT NOT NULL,
			category TEXT NOT NULL,
			context TEXT,
			timeline TEXT,
			budget DECIMAL,
			status TEXT DEFAULT 'speculative',
			previous_risk_score DECIMAL,
			projected_risk_score DECIMAL,
			risk_delta DECIMAL,
			confidence DECIMAL,
			suggestions TEXT[],
			reasoning TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			confirmed_at TIMESTAMPTZ
		);
	`)
	if err != nil {
		return err
	}

	// Create risk_profiles table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_profiles (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			startup_id UUID NOT NULL,
			score DECIMAL NOT NULL,
			level TEXT NOT NULL,
			confidence DECIMAL,
			factors TEXT[],
			suggestions TEXT[],
			reasoning TEXT,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
			updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create risk_evolutions table
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS risk_evolutions (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			startup_id UUID NOT NULL,
			score DECIMAL NOT NULL,
			level TEXT NOT NULL,
			trigger TEXT NOT NULL,
			created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP
		);
	`)
	if err != nil {
		return err
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_users_email ON users (email);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_decisions_startup_id ON decisions (startup_id);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_risk_profiles_startup_id ON risk_profiles (startup_id);`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_risk_evolutions_startup_id ON risk_evolutions (startup_id);`)
	if err != nil {
		return err
	}

	return nil
}

func migration001Down(db *sql.DB) error {
	// Drop tables in reverse order
	tables := []string{"risk_evolutions", "risk_profiles", "decisions", "startups", "users"}

	for _, table := range tables {
		_, err := db.Exec("DROP TABLE IF EXISTS " + table + " CASCADE")
		if err != nil {
			return err
		}
	}

	return nil
}

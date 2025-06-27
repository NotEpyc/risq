package migrations

const CreateInitialTables = `
-- Drop all tables first to ensure clean state
DROP TABLE IF EXISTS risk_evolutions CASCADE;
DROP TABLE IF EXISTS decisions CASCADE;
DROP TABLE IF EXISTS risk_profiles CASCADE;
DROP TABLE IF EXISTS startups CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'user',
    startup_id UUID,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Startups table
CREATE TABLE startups (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    industry TEXT NOT NULL,
    funding_stage TEXT NOT NULL CHECK (funding_stage IN ('idea', 'pre_seed', 'seed', 'series_a', 'series_b', 'series_c', 'ipo')),
    location TEXT NOT NULL,
    founded_date TIMESTAMP WITH TIME ZONE NOT NULL,
    team_size INTEGER NOT NULL DEFAULT 1,
    website TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Risk profiles table
CREATE TABLE risk_profiles (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL DEFAULT 0.0,
    level TEXT NOT NULL CHECK (level IN ('low', 'medium', 'high', 'critical')),
    confidence DECIMAL(5,2) NOT NULL DEFAULT 0.0,
    factors TEXT[] DEFAULT '{}',
    suggestions TEXT[] DEFAULT '{}',
    reasoning TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Risk evolutions table (for tracking risk changes over time)
CREATE TABLE risk_evolutions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    score DECIMAL(5,2) NOT NULL DEFAULT 0.0,
    level TEXT NOT NULL CHECK (level IN ('low', 'medium', 'high', 'critical')),
    trigger TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Decisions table
CREATE TABLE decisions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    startup_id UUID NOT NULL REFERENCES startups(id) ON DELETE CASCADE,
    description TEXT NOT NULL,
    category TEXT NOT NULL CHECK (category IN ('strategic', 'operational', 'financial', 'technical', 'marketing', 'hiring', 'product', 'other')),
    context TEXT,
    timeline TEXT,
    budget DECIMAL(12,2) DEFAULT 0.0,
    status TEXT NOT NULL DEFAULT 'pending' CHECK (status IN ('pending', 'approved', 'rejected', 'implemented', 'cancelled')),
    previous_risk_score DECIMAL(5,2) DEFAULT 0.0,
    projected_risk_score DECIMAL(5,2) DEFAULT 0.0,
    risk_delta DECIMAL(5,2) DEFAULT 0.0,
    confidence DECIMAL(5,2) DEFAULT 0.0,
    suggestions TEXT[] DEFAULT '{}',
    reasoning TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    confirmed_at TIMESTAMP WITH TIME ZONE
);

-- Add foreign key constraint for users.startup_id
ALTER TABLE users ADD CONSTRAINT fk_users_startup_id FOREIGN KEY (startup_id) REFERENCES startups(id) ON DELETE SET NULL;

-- Create indexes for better performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_startup_id ON users(startup_id);
CREATE INDEX idx_risk_profiles_startup_id ON risk_profiles(startup_id);
CREATE INDEX idx_risk_evolutions_startup_id ON risk_evolutions(startup_id);
CREATE INDEX idx_decisions_startup_id ON decisions(startup_id);
CREATE INDEX idx_decisions_status ON decisions(status);
CREATE INDEX idx_decisions_category ON decisions(category);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Create triggers for updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_startups_updated_at BEFORE UPDATE ON startups FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_risk_profiles_updated_at BEFORE UPDATE ON risk_profiles FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_decisions_updated_at BEFORE UPDATE ON decisions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
`

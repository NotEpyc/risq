-- ===============================================
-- Risk Assessment Backend - Clean Database Schema
-- ===============================================

-- Drop all tables if they exist (clean slate)
DROP TABLE IF EXISTS risk_evolutions CASCADE;
DROP TABLE IF EXISTS decisions CASCADE;
DROP TABLE IF EXISTS risk_profiles CASCADE;
DROP TABLE IF EXISTS startups CASCADE;
DROP TABLE IF EXISTS users CASCADE;

-- Create UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ===============================================
-- USERS TABLE
-- ===============================================
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

-- ===============================================
-- STARTUPS TABLE
-- ===============================================
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

-- ===============================================
-- RISK PROFILES TABLE
-- ===============================================
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

-- ===============================================
-- DECISIONS TABLE
-- ===============================================
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

-- ===============================================
-- FOREIGN KEY CONSTRAINTS
-- ===============================================
ALTER TABLE users ADD CONSTRAINT fk_users_startup_id 
    FOREIGN KEY (startup_id) REFERENCES startups(id) ON DELETE SET NULL;

-- ===============================================
-- INDEXES FOR PERFORMANCE
-- ===============================================
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_startup_id ON users(startup_id);
CREATE INDEX idx_startups_id ON startups(id);
CREATE INDEX idx_risk_profiles_startup_id ON risk_profiles(startup_id);
CREATE INDEX idx_decisions_startup_id ON decisions(startup_id);
CREATE INDEX idx_decisions_status ON decisions(status);

-- ===============================================
-- INSERT SAMPLE DATA (OPTIONAL - FOR TESTING)
-- ===============================================
-- This section can be uncommented for testing purposes

/*
-- Sample user
INSERT INTO users (id, email, name, password_hash, role) VALUES 
    ('11111111-1111-1111-1111-111111111111', 'test@example.com', 'Test User', '$2a$10$dummy.hash.for.testing', 'founder');

-- Sample startup
INSERT INTO startups (id, name, description, industry, funding_stage, location, team_size) VALUES 
    ('22222222-2222-2222-2222-222222222222', 'Test Startup', 'A sample startup for testing', 'Technology', 'seed', 'San Francisco', 3);

-- Link user to startup
UPDATE users SET startup_id = '22222222-2222-2222-2222-222222222222' WHERE id = '11111111-1111-1111-1111-111111111111';

-- Sample risk profile
INSERT INTO risk_profiles (startup_id, overall_risk, market_risk, technical_risk, financial_risk, team_risk, regulatory_risk) VALUES 
    ('22222222-2222-2222-2222-222222222222', 'medium', 0.6, 0.4, 0.7, 0.3, 0.2);
*/

-- ===============================================
-- VERIFICATION QUERIES
-- ===============================================
-- Uncomment these to verify the schema was created correctly

-- SELECT 'Schema created successfully!' as status;
-- SELECT table_name FROM information_schema.tables WHERE table_schema = 'public' ORDER BY table_name;
-- \d users;
-- \d startups;
-- \d risk_profiles;
-- \d decisions;

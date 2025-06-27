# Event-Driven Architecture Setup Guide

## 🎯 Overview

This guide explains the comprehensive event-driven architecture implementation for the Risk Assessment Backend. The system now uses NATS for event streaming and real-time processing of startup onboarding, market validation, risk analysis, and context storage.

## 🔧 Architecture Components

### Core Event Flow
```
Startup Onboarding → Market Validation → Risk Analysis → Context Storage → RAG
```

### Event Types and Handlers

| Event | Subject | Handler | Description |
|-------|---------|---------|-------------|
| `StartupOnboardedEvent` | `startup.onboarded` | MarketValidationHandler | Triggers market data collection |
| `MarketValidationRequestedEvent` | `market.validation.requested` | MarketValidationHandler | Fetches external market data |
| `MarketValidatedEvent` | `market.validated` | RiskAnalysisHandler | Triggers risk analysis |
| `RiskAnalysisRequestedEvent` | `risk.analysis.requested` | RiskAnalysisHandler | Processes AI risk analysis |
| `RiskAnalysisCompletedEvent` | `risk.analysis.completed` | ContextStorageHandler | Stores results in RAG |
| `ContextStoreRequestedEvent` | `context.store.requested` | ContextStorageHandler | Manages context storage |

## 🚀 Quick Start

### 1. Environment Setup
```bash
# Create .env file
cp .env.example .env

# Edit with your configuration
nano .env
```

### 2. Start Infrastructure
```bash
# Start PostgreSQL, Redis, and NATS
docker-compose up -d postgres redis nats

# Verify services
docker-compose ps
```

### 3. Run the Application
```bash
# Install dependencies
go mod tidy

# Run the API server
go run ./cmd/api

# In another terminal, test the event flow
go run test_event_flow.go
```

## 📋 Event Flow Testing

### Automated Testing
```bash
# Ensure NATS is running
docker-compose up -d nats

# Run comprehensive event flow test
go run test_event_flow.go
```

### Manual API Testing
```bash
# Submit startup for onboarding
curl -X POST http://localhost:8080/api/v1/startups/onboard \
-H "Content-Type: application/json" \
-d '{
  "name": "TestStartup Inc",
  "description": "AI-powered fintech solution",
  "industry": "Technology",
  "sector": "FinTech",
  "website": "https://teststartup.com",
  "funding_stage": "Seed",
  "location": "San Francisco, CA",
  "founded_date": "2024-01-15T00:00:00Z",
  "team_size": 5,
  "business_model": "B2B SaaS",
  "revenue_streams": ["Subscription", "Transaction fees"],
  "target_market": "SMB Banking",
  "competitor_analysis": "Competing with traditional banking solutions...",
  "implementation_plan": "Phase 1: MVP development...",
  "technology_stack": ["Go", "React", "PostgreSQL"],
  "development_timeline": "12 months to market",
  "go_to_market_strategy": "Direct sales and partnerships...",
  "initial_investment": 500000,
  "monthly_burn_rate": 50000,
  "projected_revenue": 1000000,
  "funding_requirement": 2000000,
  "founder_details": [{
    "name": "John Smith",
    "email": "john@teststartup.com",
    "role": "CEO & Founder",
    "linkedin_url": "https://linkedin.com/in/johnsmith",
    "education": ["MBA Stanford", "BS Computer Science MIT"],
    "experience": [{
      "company": "Google",
      "position": "Senior Software Engineer",
      "start_date": "2020-01-01T00:00:00Z",
      "end_date": "2023-12-31T00:00:00Z",
      "description": "Led payment processing systems development",
      "industry": "Technology"
    }],
    "skills": ["Go", "Python", "Financial modeling", "Team leadership"],
    "achievements": ["Led 10-person engineering team", "300% performance improvement"]
  }]
}'
```

## 🔍 Monitoring Event Flow

### 1. Application Logs
```bash
# Watch real-time logs
tail -f application.log | grep -E "(event|Event)"

# Filter by event type
tail -f application.log | grep "risk.analysis"
```

### 2. NATS Monitoring
```bash
# Check NATS server status
curl http://localhost:8222/varz | jq

# Monitor connections
curl http://localhost:8222/connz | jq

# View subscriptions
curl http://localhost:8222/subsz | jq
```

### 3. Database Monitoring
```bash
# Connect to PostgreSQL
docker exec -it risq_postgres psql -U risq_user -d risq_db

# Check recent startups
SELECT id, name, industry, created_at FROM startups ORDER BY created_at DESC LIMIT 5;

# Check risk assessments
SELECT startup_id, risk_score, risk_level, created_at FROM risk_assessments ORDER BY created_at DESC LIMIT 5;
```

## 🐛 Troubleshooting

### Common Issues

1. **NATS Connection Failed**
   ```bash
   # Check NATS container
   docker-compose ps nats
   
   # Restart NATS
   docker-compose restart nats
   
   # Check logs
   docker-compose logs nats
   ```

2. **Database Connection Issues**
   ```bash
   # Test database connection
   docker exec -it risq_postgres psql -U risq_user -d risq_db -c "SELECT 1;"
   
   # Check database logs
   docker-compose logs postgres
   ```

3. **Events Not Processing**
   ```bash
   # Check if event handlers are registered
   grep -r "Subscribe" internal/handlers/
   
   # Verify NATS subjects in logs
   tail -f application.log | grep "Subject"
   
   # Test event publishing manually
   go run test_event_flow.go
   ```

4. **OpenAI API Issues**
   ```bash
   # Test API key
   curl https://api.openai.com/v1/models \
     -H "Authorization: Bearer $OPENAI_API_KEY"
   
   # Check environment variable
   echo $OPENAI_API_KEY
   ```

### Performance Optimization

1. **NATS Performance**
   ```bash
   # Check message rates
   curl http://localhost:8222/varz | jq '.in_msgs, .out_msgs'
   
   # Monitor memory usage
   curl http://localhost:8222/varz | jq '.mem'
   ```

2. **Database Performance**
   ```sql
   -- Check slow queries
   SELECT query, mean_time, calls 
   FROM pg_stat_statements 
   ORDER BY mean_time DESC LIMIT 10;
   
   -- Monitor connections
   SELECT count(*) FROM pg_stat_activity;
   ```

## 🎯 Next Steps

### Development
1. Add comprehensive integration tests
2. Implement event replay mechanisms
3. Add distributed tracing
4. Create performance benchmarks

### Production Readiness
1. Add monitoring dashboards
2. Implement circuit breakers
3. Add rate limiting
4. Create backup strategies

### Feature Extensions
1. Add decision tracking events
2. Implement user notification events
3. Create audit trail events
4. Add webhook support

## 📊 Event Schema Reference

### StartupOnboardedEvent
```json
{
  "id": "startup-onboarded-uuid",
  "type": "startup.onboarded",
  "source": "startup-controller",
  "subject": "startup.onboarded",
  "timestamp": "2024-01-15T10:30:00Z",
  "startup_id": "uuid",
  "user_id": "uuid",
  "startup_data": { /* comprehensive startup data */ },
  "founder_cv": { /* founder background */ },
  "business_plan": { /* business plan details */ }
}
```

### RiskAnalysisCompletedEvent
```json
{
  "id": "risk-completed-uuid",
  "type": "risk.analysis.completed",
  "source": "risk-analysis-service",
  "subject": "risk.analysis.completed",
  "timestamp": "2024-01-15T10:35:00Z",
  "startup_id": "uuid",
  "risk_score": 75.5,
  "risk_level": "Medium",
  "strengths": ["Strong founder", "Large market"],
  "weaknesses": ["High competition", "Limited track record"],
  "recommendations": ["Build partnerships", "Focus on compliance"],
  "detailed_analysis": { /* comprehensive analysis data */ }
}
```

## 🤝 Contributing

1. Follow the event-driven architecture patterns
2. Add comprehensive tests for new events
3. Update documentation for new event types
4. Ensure proper error handling in event handlers

---

For questions or issues, please check the main README.md or create an issue in the repository.

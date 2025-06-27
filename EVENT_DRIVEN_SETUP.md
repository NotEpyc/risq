# Risk Assessment Backend - Event-Driven Architecture Setup

## 🎯 Quick Start - Event-Driven Startup Onboarding

This system implements a complete event-driven architecture for startup onboarding with automated market validation, risk analysis, and context-aware RAG storage.

### 🚀 Prerequisites

- Go 1.21+
- Docker & Docker Compose
- Postman or curl for testing

### 📦 Setup Instructions

#### 1. Start Infrastructure Services

```bash
# Option 1: Start specific services (Recommended)
docker-compose up -d postgres redis nats

# Option 2: Use the quick setup file
docker-compose -f docker-compose.quick.yml up -d
```

#### 2. Verify Services Are Running

```bash
# Check containers
docker ps

# Check NATS
curl http://localhost:8222/varz

# Check PostgreSQL
docker exec -it risq_postgres psql -U risq_user -d risq_db -c "SELECT 1;"

# Check Redis
docker exec -it risq_redis redis-cli ping
```

#### 3. Build and Run the API

```bash
# Build the application
go build -o api ./cmd/api

# Run the application
./api
```

The API will start on `http://localhost:8080` with full event-driven capabilities.

## 🔄 Event-Driven Flow Overview

The system implements the following event chain:

```
1. Startup Onboarding (API Call)
   ↓
2. StartupOnboardedEvent → Market Validation Handler
   ↓
3. MarketValidationRequestedEvent → External Market Data Service
   ↓
4. MarketValidatedEvent → Risk Analysis Handler
   ↓
5. RiskAnalysisRequestedEvent → OpenAI Risk Analysis
   ↓
6. RiskAnalysisCompletedEvent → Context Storage Handler
   ↓
7. Context Stored in RAG (Vector Database)
```

### 📋 Event Subjects (NATS Topics)

- `startup.onboarded` - Startup has completed onboarding
- `market.validation.requested` - Market validation request
- `market.validated` - Market validation completed
- `risk.analysis.requested` - Risk analysis request
- `risk.analysis.completed` - Risk analysis completed
- `context.store.requested` - Context storage request
- `context.stored` - Context successfully stored

## 🧪 Testing the Complete Flow

### 1. User Registration

```bash
curl -X POST http://localhost:8080/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

### 2. User Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john.doe@example.com",
    "password": "password123"
  }'
```

**Copy the JWT token from the response for use in subsequent requests.**

### 3. Startup Onboarding (Triggers Event Chain)

```bash
curl -X POST http://localhost:8080/api/v1/startup/onboard \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "TechFlow AI",
    "description": "AI-powered workflow automation platform for enterprises",
    "industry": "technology",
    "funding_stage": "seed",
    "location": "San Francisco",
    "founded_date": "2024-01-15T00:00:00Z",
    "team_size": 5
  }'
```

This single API call will trigger the entire event-driven workflow:

1. ✅ **Startup created** - Basic validation and database storage
2. 🔄 **Event published** - `startup.onboarded` event sent to NATS
3. 🏪 **Market validation** - External APIs called for market data
4. 📊 **Risk analysis** - OpenAI analyzes comprehensive data
5. 🧠 **Context storage** - Results stored in RAG for future queries

### 4. Check Risk Assessment Results

```bash
curl -X GET http://localhost:8080/api/v1/risk/current \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### 5. View Risk Evolution History

```bash
curl -X GET http://localhost:8080/api/v1/risk/history \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## 📊 Example Response Flow

### Initial Startup Response
```json
{
  "status": "success",
  "message": "Startup onboarded successfully and event-driven analysis initiated",
  "data": {
    "startup": {
      "id": "123e4567-e89b-12d3-a456-426614174000",
      "name": "TechFlow AI",
      "industry": "technology",
      "funding_stage": "seed"
    },
    "risk_profile": {
      "score": 45.2,
      "level": "medium",
      "factors": ["Market competition", "Team experience"],
      "suggestions": ["Conduct market research", "Build team expertise"]
    }
  }
}
```

### Enhanced Risk Analysis (After Events)
The event-driven analysis will enhance this with:
- **Market data integration** - Real competitor analysis
- **Sector health assessment** - Industry trends and outlook  
- **AI-powered recommendations** - Detailed improvement strategies
- **Context-aware insights** - Stored for future decision making

## 🔧 Architecture Components

### Event Handlers

1. **MarketValidationHandler** (`internal/handlers/market_validation.go`)
   - Processes startup onboarding events
   - Calls external market data APIs
   - Publishes market validation results

2. **RiskAnalysisHandler** (`internal/handlers/risk_analysis.go`) 
   - Processes market validation completion
   - Calls OpenAI for comprehensive risk analysis
   - Publishes risk analysis results

3. **ContextStorageHandler** (`internal/handlers/context_storage.go`)
   - Processes risk analysis completion
   - Stores comprehensive context in RAG
   - Enables future context-aware queries

### External Services

1. **Market Data Service** (`pkg/external/market_data.go`)
   - Industry trends analysis
   - Sector health assessment
   - Competitor landscape evaluation
   - News sentiment analysis

2. **Event Service** (`pkg/events/events.go`)
   - NATS event publishing/subscribing
   - Event type definitions
   - Message routing and handling

### Key Features

- ✅ **Input Validation** - Comprehensive startup data validation
- 🔄 **Event-Driven Architecture** - Asynchronous processing pipeline
- 🌐 **External API Integration** - Real market data (simulated for demo)
- 🤖 **OpenAI Integration** - AI-powered risk analysis
- 🧠 **RAG Context Storage** - Vector-based context storage
- 📊 **Production-Ready** - Clean architecture, error handling, logging

## 🐛 Troubleshooting

### NATS Connection Issues
```bash
# Check NATS server status
curl http://localhost:8222/varz

# View NATS logs
docker logs risq_nats
```

### Database Connection Issues
```bash
# Check PostgreSQL status
docker exec -it risq_postgres pg_isready -U risq_user

# View database logs
docker logs risq_postgres
```

### Event Processing Issues
Check application logs for event processing:
```bash
# The application logs will show event flow
./api 2>&1 | grep -E "(Event|NATS|Handler)"
```

## 🎉 Next Steps

1. **Decision Speculation** - Test business decisions with AI impact analysis
2. **Context Queries** - Query the RAG system for business insights  
3. **Scaling** - Add more event handlers for additional business logic
4. **Integration** - Connect to real external APIs (replace simulated data)
5. **Dashboard** - Build a frontend to visualize the event-driven insights

The system is now ready for production-grade startup onboarding with complete event-driven market validation, risk analysis, and context-aware intelligence!

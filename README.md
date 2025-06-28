# 🎯 RISQ Backend - AI-Powered Startup Risk Assessment API

[![Go Version](https://img.shields.io/badge/Go-1.23+-blue.svg)](https://golang.org)
[![Railway Deploy](https://img.shields.io/badge/Deploy-Railway-blueviolet.svg)](https://railway.app)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![API Version](https://img.shields.io/badge/API-v1.0.0-orange.svg)](https://resqbackend-production.up.railway.app/health)

> **An intelligent, event-driven platform that leverages AI to provide comprehensive risk assessment and decision-making support for startups.**

## 🌟 **Overview**

RISQ Backend is a sophisticated microservice that combines machine learning, real-time market data, and context-aware AI to deliver actionable insights for startup risk assessment. Built with enterprise-grade architecture, it features event-driven workflows, vector embeddings for context memory, and seamless integration with external data sources.

## 🚀 **Key Features**

### **🤖 AI-Powered Intelligence**
- **GPT-4 Integration**: Advanced natural language processing for decision reasoning
- **Vector Embeddings**: Semantic context storage and retrieval using OpenAI embeddings  
- **Context Memory**: RAG (Retrieval-Augmented Generation) system for learning from startup history

### **📊 Risk Assessment Engine**
- **Multi-dimensional Analysis**: Market, Technical, Financial, Regulatory, and Operational risk scoring
- **Real-time Scoring**: Dynamic risk calculation (0-100 scale) with confidence intervals
- **Historical Tracking**: Risk evolution monitoring and trend analysis

### **🎯 Decision Support System**
- **AI Speculation**: Intelligent decision recommendations with reasoning
- **Confidence Scoring**: Quantified certainty levels for all suggestions
- **Decision History**: Complete audit trail of all decisions and outcomes

### **� Real-time Market Intelligence**
- **News Sentiment Analysis**: Automated market trend detection via NewsAPI
- **Industry Data**: Real-time financial market data integration via MarketStack
- **Sector-specific Insights**: Tailored analysis based on startup industry

### **⚡ Event-Driven Architecture**
- **Asynchronous Processing**: Non-blocking workflows for heavy AI computations
- **NATS Messaging**: Reliable event streaming and service communication
- **Workflow Orchestration**: Automated startup analysis pipeline

## 📦 Tech Stack

- **Backend**: Go 1.24+ with Fiber framework
- **Database**: PostgreSQL with custom SQL queries and migrations
- **Cache**: Redis for caching
- **AI**: OpenAI GPT-4 for risk assessment
- **Containerization**: Docker & Docker Compose

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- OpenAI API Key (optional - will use default values if not provided)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd risq_backend
   ```

2. **Set up environment variables (optional)**
   ```bash
   # Create a .env file in the root directory with your configuration
   # Example:
   # OPENAI_API_KEY=your_api_key_here
   # DB_PASSWORD=your_db_password
   ```

3. **Start development services**
   ```bash
   docker-compose -f docker-compose.dev.yml up -d
   ```

4. **Run the Go application**
   ```bash
   go run main.go
   ```

The API will be available at `http://localhost:8080`

### Available Services

When you start the development services, the following will be available:

- **API**: http://localhost:8080 (run locally with `go run main.go`)
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379

## 📁 Project Structure

```
risq_backend/
├── docker-compose.yml           # Production Docker services
├── docker-compose.dev.yml       # Development Docker services
├── go.mod                       # Go module dependencies
├── main.go                      # Application entrypoint (simplified)
├── README.md                    # This file

├── api/                         # API layer (HTTP)
│   ├── controller/             # Route handlers
│   │   ├── user_controller.go
│   │   ├── startup_controller.go
│   │   ├── decision_controller.go
│   │   └── risk_controller.go
│   └── routes.go               # Route definitions

├── cmd/api/main.go             # Application entrypoint (full)

├── config/                     # Configuration loader
│   └── config.go

├── internal/                   # Domain logic
│   ├── user/                   # User management
│   ├── startup/                # Startup profiles
│   ├── decision/               # Decision speculation/confirmation
│   ├── risk/                   # Risk assessment
│   ├── llm/                    # LLM integration
│   └── contextmem/             # Context memory for decisions

├── pkg/                        # Infrastructure
│   ├── app/                    # Application setup & DI  
│   ├── database/               # PostgreSQL connection
│   ├── cache/                  # Redis connection
│   ├── logger/                 # Structured logging
│   ├── llm/                    # OpenAI client
│   ├── middlewares/            # HTTP middlewares
│   └── response/               # Standardized responses

└── types/                      # Shared type definitions
    ├── decision_types.go
    ├── risk_types.go
    └── jwt.go
```

## 🔧 API Endpoints

### Health Check
- `GET /health` - Service health status

### User Management
- `POST /api/v1/public/users` - Create new user
- `GET /api/v1/public/users/email?email=user@example.com` - Get user by email
- `GET /api/v1/users/:id` - Get user by ID (authenticated)

### Startup Management
- `POST /api/v1/startups/submit` - Submit startup for onboarding
- `GET /api/v1/startups/me` - Get current user's startup
- `GET /api/v1/startups/:id` - Get startup by ID

### Decision Engine
- `POST /api/v1/decisions/speculate` - Speculate on a decision
- `POST /api/v1/decisions/confirm` - Confirm a speculated decision
- `GET /api/v1/decisions/` - Get all decisions for startup
- `GET /api/v1/decisions/:id` - Get specific decision

### Risk Assessment
- `GET /api/v1/risks/current` - Get current risk profile
- `GET /api/v1/risks/evolution` - Get risk evolution timeline

## 💡 Core Workflows

### 1. Startup Onboarding
```json
POST /api/v1/startups/submit
{
  "name": "TechCorp",
  "description": "AI-powered analytics platform",
  "industry": "technology",
  "funding_stage": "seed",
  "location": "San Francisco",
  "founded_date": "2024-01-01",
  "team_size": 5
}
```

### 2. Decision Speculation
```json
POST /api/v1/decisions/speculate
{
  "startup_id": "uuid",
  "description": "Hire 10 additional engineers",
  "category": "hiring",
  "context": "Expanding product team for new features",
  "timeline": "3 months",
  "budget": 150000
}
```

### 3. Decision Confirmation
```json
POST /api/v1/decisions/confirm
{
  "decision_id": "uuid",
  "notes": "Proceeding with hiring plan"
}
```

## 🧠 AI Integration

### Risk Analysis
- Uses OpenAI GPT-4 for intelligent risk assessment
- Combines multiple data sources: startup info, market data, historical decisions
- Provides confidence scores and detailed reasoning

### Context Memory
- Stores decision history as vector embeddings
- Enables RAG (Retrieval-Augmented Generation) for context-aware analysis
- Improves accuracy over time as more decisions are made

### Mitigation Suggestions
- AI-generated actionable risk mitigation strategies
- Categorized by risk type (market, financial, operational, etc.)
- Includes implementation steps, resources, and timelines

## 🔄 Event-Driven Architecture

The system uses Redis Streams for event-driven communication:

- `startup.profile.created` - Triggers initial risk assessment
- `decision.confirmed` - Updates risk scores and context memory
- `risk.score.updated` - Notifies downstream systems
- `suggestion.generated` - Creates mitigation recommendations

## 🚀 Deployment

### Docker Production
```bash
# Build and deploy
docker-compose -f docker-compose.prod.yml up -d

# Scale services
docker-compose -f docker-compose.prod.yml up -d --scale app=3
```

### Environment Variables
```bash
# Required
DATABASE_URL=postgres://user:pass@host:port/db
REDIS_URL=redis://host:port
OPENAI_API_KEY=your_openai_key
JWT_SECRET=your_jwt_secret

# Optional
PINECONE_API_KEY=your_pinecone_key
ALPHA_VANTAGE_API_KEY=your_market_data_key
LOG_LEVEL=info
PORT=8080
```

## 🔧 Environment Setup

### 1. Environment Variables

Before running the application, you need to set up your environment variables:

1. **Copy the example environment file:**
   ```bash
   cp .env.example .env
   ```

2. **Edit `.env` and add your actual values:**
   ```bash
   # OpenAI Configuration (Required)
   OPENAI_API_KEY=your_actual_openai_api_key_here
   
   # JWT Secret (Required for production)
   JWT_SECRET=your_secure_jwt_secret_here
   
   # Database credentials (if different from defaults)
   DB_PASSWORD=your_database_password
   ```

3. **Get your OpenAI API key:**
   - Visit [OpenAI API Keys](https://platform.openai.com/api-keys)
   - Create a new API key
   - Add it to your `.env` file

⚠️ **Security Note**: Never commit your `.env` file to version control. It contains sensitive API keys and secrets.

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/decision/
```

## 🔮 Future Enhancements

- [ ] Real-time market data integration
- [ ] Advanced vector search with Pinecone
- [ ] Multi-tenant support
- [ ] GraphQL API
- [ ] Real-time notifications
- [ ] Advanced analytics dashboard
- [ ] Integration with startup ecosystems
- [ ] Mobile app support

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 📞 Support

For support and questions:
- Create an issue in the repository
- Email: support@risq-backend.com
- Documentation: [docs.risq-backend.com](https://docs.risq-backend.com)

---

**Built with ❤️ for startup founders by startup founders**

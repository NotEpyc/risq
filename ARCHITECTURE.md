# Risk Assessment Backend - Clean Architecture

## Event Flow Overview

The application follows a clean, linear event flow designed for first-time entrepreneurs:

```
1. SIGNUP → 2. LOGIN → 3. STARTUP ONBOARDING → 4. RISK ASSESSMENT → 5. DECISION FLOW → 6. CONTEXT STORAGE
```

## 1. User Authentication Flow

### Signup (`POST /api/v1/auth/signup`)
- User creates account with email, name, and password
- System validates input (email format, password strength)
- Password is hashed using bcrypt
- User gets JWT token automatically after signup
- Default role: "founder"

### Login (`POST /api/v1/auth/login`)
- User authenticates with email/password
- System validates credentials
- Returns JWT token on success
- Token expires in 24 hours

## 2. Startup Onboarding Flow

### Submit Startup Profile (`POST /api/v1/startup/onboard`)
- **Requires**: JWT authentication
- User submits startup information:
  - Name, description, industry
  - Funding stage, location
  - Team size, founded date
- System validates required fields
- Prevents duplicate startup creation per user
- Links startup to user account via `startup_id`
- Triggers initial risk assessment
- Stores context for future decisions

## 3. Risk Assessment Flow

### Initial Risk Profile Creation
- Automatically triggered after startup onboarding
- AI analyzes startup information using OpenAI GPT-4
- Generates initial risk score (0-100 scale)
- Identifies key risk factors
- Provides mitigation suggestions
- Stores baseline for future comparisons

### Get Current Risk (`GET /api/v1/risk/current`)
- **Requires**: JWT + Startup onboarding
- Returns latest risk assessment
- Shows current risk score and factors

### Risk History (`GET /api/v1/risk/history`)
- **Requires**: JWT + Startup onboarding
- Shows risk score evolution over time
- Tracks impact of confirmed decisions

## 4. Decision Speculation & Confirmation Flow

### Decision Speculation (`POST /api/v1/decisions/speculate`)
- **Requires**: JWT + Startup onboarding
- User describes a potential decision
- System analyzes impact using AI
- Returns projected risk score change
- Provides confidence level and suggestions
- Stores as "speculative" decision
- Does NOT affect actual risk score

### Decision Confirmation (`POST /api/v1/decisions/confirm`)
- **Requires**: JWT + Startup ownership validation
- Confirms a previously speculated decision
- Updates actual risk score
- Changes status from "speculative" to "confirmed"
- Triggers context memory storage

### Decision Management
- `GET /api/v1/decisions/` - List all decisions for startup
- `GET /api/v1/decisions/:id` - Get specific decision
- All endpoints validate startup ownership

## 5. Context Memory & Learning

### Automatic Context Storage
- Startup onboarding data stored as context
- Confirmed decisions stored with embeddings
- Vector search enables contextual AI responses
- Historical decisions inform future speculation

### Context Retrieval
- System fetches relevant historical context
- Uses vector similarity search
- Improves AI accuracy over time
- Provides context-aware risk assessment

## Security & Validation

### Authentication Flow
1. JWT middleware validates all protected routes
2. StartupContext middleware ensures startup ownership
3. Controllers validate user owns requested resources

### Input Validation
- Email format validation
- Password strength requirements
- Required field validation
- UUID format validation
- Startup ownership verification

### Data Privacy
- Passwords never returned in API responses
- JWT tokens contain minimal claims
- Context data scoped by startup ID

## Database Architecture

### User-Startup Relationship
- Users table has optional `startup_id` foreign key
- One-to-one relationship (one user, one startup)
- Startup creation updates user record

### Decision-Risk Relationship
- Decisions linked to startups via `startup_id`
- Risk assessments track decision impacts
- Historical risk data preserved

## Error Handling

### Validation Errors
- 400 Bad Request for invalid input
- Clear error messages for missing fields
- Input format validation

### Authentication Errors
- 401 Unauthorized for missing/invalid tokens
- 403 Forbidden for insufficient permissions

### Business Logic Errors
- Prevents duplicate startup creation
- Validates decision ownership
- Ensures proper event flow sequence

## API Endpoints Summary

### Authentication
- `POST /api/v1/auth/signup` - User registration
- `POST /api/v1/auth/login` - User authentication

### Startup Management
- `POST /api/v1/startup/onboard` - Startup profile creation
- `GET /api/v1/startup/profile` - Get user's startup

### Risk Assessment
- `GET /api/v1/risk/current` - Current risk profile
- `GET /api/v1/risk/history` - Risk evolution history

### Decision Management
- `POST /api/v1/decisions/speculate` - Test decision impact
- `POST /api/v1/decisions/confirm` - Confirm speculated decision
- `GET /api/v1/decisions/` - List startup decisions
- `GET /api/v1/decisions/:id` - Get specific decision

## Production Best Practices

### Clean Architecture
- Separation of concerns (Controller → Service → Repository)
- Domain-driven design with clear bounded contexts
- Dependency injection for testability

### Error Handling
- Structured error responses
- Comprehensive logging
- Graceful degradation for AI failures

### Performance
- Redis caching for context data
- Connection pooling for database
- Efficient vector storage and retrieval

### Scalability
- Stateless API design
- Horizontal scaling ready
- Microservices preparation

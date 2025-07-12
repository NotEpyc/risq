# 🏗️ RISQ Backend - Complete Flutter Developer Guide & API Documentation

## 🚀 Project Overview

**RISQ** is an AI-powered risk assessment platform for early-stage startups and entrepreneurs. Think of it as a smart advisor that helps entrepreneurs make better business decisions by analyzing risk and providing AI-powered suggestions.

### What Does This App Do?
1. **User Registration & Login** - Entrepreneurs create accounts
2. **Startup Onboarding** - Users provide detailed information about their startup idea/business
3. **AI Risk Analysis** - The system analyzes the startup and provides risk scores and suggestions
4. **Decision Testing** - Users can "test" business decisions before making them to see how they might affect risk
5. **Risk Tracking** - Monitor how risks change over time

### Key Features for Your Flutter App
- **User Authentication** with JWT tokens (login stays active across app restarts)
- **Startup Profile Creation** with forms and data collection
- **AI-Powered Risk Dashboard** showing scores, graphs, and suggestions
- **Decision Simulator** - "What if I hire 5 people?" type scenarios
- **Risk History Charts** - Show how risk changes over time
- **Real-time Updates** - The backend processes data in the background

### Technical Architecture
```
Flutter App ↔️ REST API ↔️ PostgreSQL Database
                     ↕️
                OpenAI API (for AI analysis)
                     ↕️
            Event System (for real-time processing)
```

---

## 🌐 API Configuration

### Base URL
```
Development: http://localhost:8080
Production: https://your-domain.com (when deployed)
```

### Authentication
The API uses JWT (JSON Web Tokens). After login, you get a token that you must include in all protected requests:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

---

## 📊 Standard API Response Format

**Every API endpoint returns this same structure - this makes Flutter development much easier!**

### Success Response
```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": {
    // The actual data you need is always here
  },
  "error": null
}
```

### Error Response
```json
{
  "success": false,
  "message": "User-friendly error message",
  "data": null,
  "error": "Technical error details"
}
```

### Common HTTP Status Codes
- `200` - Success
- `400` - Bad Request (validation errors, missing fields)
- `401` - Unauthorized (invalid/missing token)
- `404` - Not Found (resource doesn't exist)
- `500` - Server Error (something went wrong on our end)

---

## 🔐 Authentication Endpoints

> **Flutter Implementation Tip:** Create an `AuthService` class to handle these endpoints and store the JWT token using `shared_preferences` package.

### 1. User Signup (Create Account)

**Endpoint:** `POST /api/v1/auth/signup`

**Description:** Creates a new user account. Perfect for your "Sign Up" screen in Flutter.

**Request Body:** (What you send from Flutter)
```json
{
  "email": "john.doe@example.com",
  "name": "John Doe", 
  "password": "SecurePassword123!"
}
```

**Validation Rules:**
- `email`: Must be valid email format (user@example.com)
- `name`: Required, 2-100 characters
- `password`: Required, minimum 8 characters

**Success Response:** (What you get back - Status Code 200)
```json
{
  "success": true,
  "message": "User created successfully",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "name": "John Doe",
      "role": "founder",
      "startup_id": null,
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- **400 Bad Request:** Email already exists, password too short, invalid email format
- **500 Internal Error:** Server problem

**Flutter Usage:**
```dart
// Store the token for future API calls
String token = response.data['token'];
await SharedPreferences.getInstance().setString('jwt_token', token);

// Navigate to startup onboarding screen
Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => OnboardingScreen()));
```

---

### 2. User Login

**Endpoint:** `POST /api/v1/auth/login`

**Description:** Authenticate existing user. Use this for your "Login" screen.

**Request Body:**
```json
{
  "email": "john.doe@example.com",
  "password": "SecurePassword123!"
}
```

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "user": {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "email": "john.doe@example.com",
      "name": "John Doe",
      "role": "founder",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Wrong email or password
- **400 Bad Request:** Missing email or password

**Flutter Usage:**
```dart
// Check if user has completed startup onboarding
if (response.data['user']['startup_id'] != null) {
  // User has startup - go to main app
  Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => HomeScreen()));
} else {
  // User needs to complete onboarding
  Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => OnboardingScreen()));
}
```

---

## 🏢 Startup Management Endpoints

> **Flutter Implementation Tip:** These endpoints help you build the onboarding flow and profile management screens.

### 3. Startup Onboarding (Complete Profile)

**Endpoint:** `POST /api/v1/startup/onboard`

**Authentication:** Required (Include JWT token in headers)

**Description:** This is the big one! Users fill out detailed startup information. The backend then automatically:
- Analyzes the startup with AI
- Creates a risk profile
- Generates suggestions
- Sets up market validation
- Prepares everything for the dashboard

**Request Body:** (This is a large form - consider breaking into multiple screens in Flutter)
```json
{
  "name": "EduTech Solutions",
  "description": "AI-powered personalized learning platform for K-12 students",
  "industry": "Education Technology",
  "sector": "edutech",
  "funding_stage": "seed",
  "location": "San Francisco, CA",
  "founded_date": "2024-01-15",
  "team_size": 5,
  "website": "https://edutech-solutions.com",
  
  "business_model": "SaaS subscription model with tiered pricing",
  "revenue_streams": [
    "Monthly subscription fees",
    "Premium content licensing", 
    "Corporate training packages"
  ],
  "target_market": "K-12 schools and homeschooling families",
  "competitor_analysis": "Competing with Khan Academy and Coursera K-12",
  
  "implementation_plan": "Phase 1: MVP development (3 months), Phase 2: Beta testing (2 months)",
  "technology_stack": ["React Native", "Node.js", "PostgreSQL", "OpenAI API"],
  "development_timeline": "6 months to full launch",
  "go_to_market_strategy": "Direct sales to schools, partnerships with educational consultants",
  
  "initial_investment": 250000,
  "monthly_burn_rate": 45000,
  "projected_revenue": 500000,
  "funding_requirement": 1000000,
  
  "founder_details": [
    {
      "name": "John Doe",
      "email": "john.doe@example.com",
      "role": "CEO & Co-founder",
      "linkedin_url": "https://linkedin.com/in/johndoe",
      "education": ["MBA from Stanford", "BS Computer Science from MIT"],
      "experience": [
        {
          "company": "Google",
          "position": "Senior Software Engineer",
          "start_date": "2020-01-01",
          "end_date": "2023-12-31",
          "description": "Led AI/ML projects for Google Education",
          "industry": "Technology"
        }
      ],
      "skills": ["Machine Learning", "Product Management", "Team Leadership"],
      "achievements": ["Led team that increased user engagement by 300%"],
      "previous_startups": ["EdTech Analytics (successful exit)"]
    }
  ]
}
```

**Field Descriptions for Your Flutter Forms:**

**Basic Information:**
- `name`: Startup name (required)
- `description`: What does your startup do? (required)
- `industry`: Broad category like "Technology", "Healthcare" (required)
- `sector`: Specific area like "edutech", "fintech", "healthtech" (required)
- `funding_stage`: "idea", "pre_seed", "seed", "series_a", "series_b", "series_c", "ipo"
- `location`: City, State/Country
- `founded_date`: "YYYY-MM-DD" format
- `team_size`: Number of people
- `website`: Optional URL

**Business Model:**
- `business_model`: How you make money (required)
- `revenue_streams`: Array of income sources (required)
- `target_market`: Who are your customers (required)
- `competitor_analysis`: Who do you compete with

**Implementation:**
- `implementation_plan`: Your roadmap (required)
- `technology_stack`: Array of technologies you use
- `development_timeline`: How long to build
- `go_to_market_strategy`: How you'll get customers

**Financial Information:**
- `initial_investment`: Money invested so far (USD)
- `monthly_burn_rate`: Money spent per month (USD)
- `projected_revenue`: Expected annual revenue (USD)
- `funding_requirement`: Money needed to raise (USD)

**Founder Information:**
- `founder_details`: Array of founder profiles (required)
  - Each founder needs: name, email, role, experience, skills

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Startup onboarding completed successfully",
  "data": {
    "startup": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "EduTech Solutions",
      "description": "AI-powered personalized learning platform for K-12 students",
      "industry": "Education Technology",
      "funding_stage": "seed",
      "location": "San Francisco, CA", 
      "founded_date": "2024-01-15T00:00:00Z",
      "team_size": 5,
      "website": "https://edutech-solutions.com",
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- **400 Bad Request:** Missing required fields, invalid dates, invalid email formats
- **401 Unauthorized:** Invalid or missing JWT token
- **500 Internal Error:** Server problem

**Flutter Usage:**
```dart
// After successful onboarding, the user can now access all other features
// Navigate to main dashboard
Navigator.pushReplacement(context, MaterialPageRoute(builder: (context) => DashboardScreen()));

// The backend automatically starts AI analysis in the background
// You can now call the risk assessment endpoints
```

---

### 4. Get Startup Profile

**Endpoint:** `GET /api/v1/startup/profile`

**Authentication:** Required (JWT token)

**Description:** Get the current user's startup information. Use this to populate profile screens or check if onboarding is complete.

**Request:** No body needed, just the JWT token in headers

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Startup profile retrieved successfully",
  "data": {
    "startup": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "EduTech Solutions",
      "description": "AI-powered personalized learning platform for K-12 students",
      "industry": "Education Technology",
      "funding_stage": "seed",
      "location": "San Francisco, CA",
      "founded_date": "2024-01-15T00:00:00Z",
      "team_size": 5,
      "website": "https://edutech-solutions.com",
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    }
  }
}
```

**Error Responses:**
- **400 Bad Request:** User hasn't completed onboarding yet
- **401 Unauthorized:** Invalid or missing JWT token

**Flutter Usage:**
```dart
// Use this to check if user needs onboarding
try {
  var response = await apiClient.get('/api/v1/startup/profile');
  // User has startup - show main app
  return DashboardScreen();
} catch (e) {
  // User needs onboarding
  return OnboardingScreen();
}
```
- All basic fields (name, description, industry, etc.) are required
- `funding_stage`: Must be one of: "idea", "pre_seed", "seed", "series_a", "series_b", "series_c", "ipo"
- `founder_details`: At least one founder required
- Financial values must be non-negative numbers

**Success Response (200):**
```json
{
  "success": true,
  "message": "Startup onboarded successfully",
  "data": {
    "startup": {
      "id": "660e8400-e29b-41d4-a716-446655440000",
      "name": "EduTech Solutions",
      "description": "AI-powered personalized learning platform for K-12 students",
      "industry": "Education Technology",
      "funding_stage": "seed",
      "location": "San Francisco, CA",
      "founded_date": "2024-01-15T00:00:00Z",
      "team_size": 5,
      "website": "https://edutech-solutions.com",
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    },
    "risk_profile": {
      "id": "770e8400-e29b-41d4-a716-446655440000",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "score": 6.8,
      "level": "medium",
      "confidence": 0.85,
      "factors": [
        "Strong technical team with relevant experience",
        "Competitive market with established players",
        "AI technology provides differentiation",
        "Funding stage appropriate for development timeline"
      ],
      "suggestions": [
        "Focus on unique AI differentiation in marketing",
        "Develop partnerships with educational institutions early",
        "Consider pilot programs to validate product-market fit"
      ],
      "reasoning": "The startup shows strong technical foundation with experienced team from Google. The education technology market is competitive but growing rapidly. AI personalization provides clear differentiation opportunity.",
      "created_at": "2025-06-28T10:30:00Z",
      "updated_at": "2025-06-28T10:30:00Z"
    },
    "message": "Complete startup profile created successfully with AI-powered risk analysis initiated"
  }
}
```

**Error Responses:**
- `400` - Missing required fields, invalid funding stage, user already has startup
- `401` - Authentication required
- `500` - Internal server error

**Event-Driven Processing:**
After successful onboarding, the system automatically triggers:
1. **Market Validation** - AI analysis of market conditions
2. **Risk Analysis** - Comprehensive risk assessment
3. **Context Storage** - Startup information stored for future decisions

---

### 4. Get Startup Profile

**Endpoint:** `GET /api/v1/startup/profile`

**Authentication:** Required (JWT Token)

**Description:** Retrieve the current user's startup profile.

**Success Response (200):**
```json
{
  "success": true,
  "message": "Startup profile retrieved successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "name": "EduTech Solutions",
    "description": "AI-powered personalized learning platform for K-12 students",
    "industry": "Education Technology",
    "funding_stage": "seed",
    "location": "San Francisco, CA",
    "founded_date": "2024-01-15T00:00:00Z",
    "team_size": 5,
    "website": "https://edutech-solutions.com",
    "created_at": "2025-06-28T10:30:00Z",
    "updated_at": "2025-06-28T10:30:00Z"
  }
}
```

**Error Responses:**
- `401` - Authentication required
- `404` - User has no startup profile (must complete onboarding first)
- `500` - Internal server error

---

---

## 🎯 Decision Management Endpoints

> **Flutter Implementation Tip:** These endpoints power the "Decision Simulator" feature - like a crystal ball for business decisions! Users can ask "What if I hire 5 people?" and get AI-powered risk analysis.

*All decision endpoints require both authentication and startup context (user must have completed onboarding).*

### 5. Speculate Decision (Test a Decision)

**Endpoint:** `POST /api/v1/decisions/speculate`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** This is the core "What-If" feature! Users describe a potential business decision and the AI analyzes how it might affect their startup's risk. Think of it as a business decision simulator.

**Use Cases in Flutter:**
- "Decision Simulator" screen with form inputs
- "Risk Impact Calculator" 
- "AI Business Advisor" feature

**Request Body:**
```json
{
  "description": "Hire 3 senior engineers to accelerate product development",
  "category": "hiring",
  "context": "Current team is 5 people, need to launch MVP in 3 months. Competition is heating up.",
  "timeline": "2 months",
  "budget": 180000
}
```

**Field Descriptions for Flutter Forms:**
- `description`: What decision are you considering? (required) - Text input
- `category`: Type of decision (required) - Dropdown with these options:
  - `hiring` - Team expansion, recruitment
  - `funding` - Investment rounds, grants  
  - `product` - Feature development, pivots
  - `marketing` - Campaigns, partnerships
  - `operations` - Process changes, tools
  - `strategy` - Business model changes
  - `legal` - Compliance, contracts
  - `other` - Miscellaneous decisions
- `context`: Why are you considering this? (optional) - Text area
- `timeline`: When would you implement this? (optional) - Text input
- `budget`: How much would this cost? (optional) - Number input (USD)

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Decision speculation completed",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "startup_id": "660e8400-e29b-41d4-a716-446655440000",
    "description": "Hire 3 senior engineers to accelerate product development",
    "category": "hiring",
    "context": "Current team is 5 people, need to launch MVP in 3 months. Competition is heating up.",
    "timeline": "2 months",
    "budget": 180000,
    "status": "speculative",
    "previous_risk_score": 6.8,
    "projected_risk_score": 7.2,
    "risk_delta": 0.4,
    "confidence": 0.82,
    "suggestions": [
      "Ensure new hires have relevant EdTech experience",
      "Consider contracting before full-time hiring",
      "Budget for additional infrastructure costs",
      "Plan for longer onboarding time given tight timeline"
    ],
    "reasoning": "Hiring senior engineers will accelerate development but increases burn rate and management complexity. The timeline pressure justifies the risk increase. Strong technical leadership will be crucial for success.",
    "created_at": "2025-06-28T11:00:00Z",
    "updated_at": "2025-06-28T11:00:00Z",
    "confirmed_at": null
  }
}
```

**What the Response Means for Flutter:**
- `previous_risk_score`: Current risk score (6.8/10)
- `projected_risk_score`: What risk would be after this decision (7.2/10)  
- `risk_delta`: Change in risk (+0.4 means risk increased)
- `confidence`: How confident AI is in this analysis (0.82 = 82%)
- `suggestions`: Array of actionable advice from AI
- `reasoning`: AI's explanation of the analysis
- `status`: "speculative" means it's just a simulation, not confirmed

**Flutter UI Ideas:**
```dart
// Show risk change with colored indicators
Color getRiskDeltaColor(double delta) {
  if (delta > 0) return Colors.red; // Risk increased
  if (delta < 0) return Colors.green; // Risk decreased
  return Colors.grey; // No change
}

// Display suggestions as cards or list tiles
ListView.builder(
  itemCount: suggestions.length,
  itemBuilder: (context, index) => Card(
    child: ListTile(
      leading: Icon(Icons.lightbulb),
      title: Text(suggestions[index]),
    ),
  ),
);
```

**Error Responses:**
- **400 Bad Request:** Invalid input, missing required fields
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** User has no startup (must complete onboarding first)
- **500 Internal Error:** Server problem

---

### 6. Confirm Decision (Make It Real)

**Endpoint:** `POST /api/v1/decisions/confirm`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** After speculating on a decision, users can confirm it to actually update their startup's risk profile. This moves the decision from "speculative" to "confirmed" and affects the real risk calculations.

**Request Body:**
```json
{
  "decision_id": "880e8400-e29b-41d4-a716-446655440000",
  "notes": "Decided to proceed with hiring. Focusing on candidates with EdTech background."
}
```

**Field Descriptions:**
- `decision_id`: ID from the speculation response (required)
- `notes`: Any additional notes about the decision (optional)

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Decision confirmed successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "startup_id": "660e8400-e29b-41d4-a716-446655440000",
    "description": "Hire 3 senior engineers to accelerate product development",
    "category": "hiring",
    "context": "Current team is 5 people, need to launch MVP in 3 months. Competition is heating up.",
    "timeline": "2 months",
    "budget": 180000,
    "status": "confirmed",
    "previous_risk_score": 6.8,
    "projected_risk_score": 7.2,
    "risk_delta": 0.4,
    "confidence": 0.82,
    "suggestions": [
      "Ensure new hires have relevant EdTech experience",
      "Consider contracting before full-time hiring",
      "Budget for additional infrastructure costs",
      "Plan for longer onboarding time given tight timeline"
    ],
    "reasoning": "Hiring senior engineers will accelerate development but increases burn rate and management complexity. The timeline pressure justifies the risk increase. Strong technical leadership will be crucial for success.",
    "created_at": "2025-06-28T11:00:00Z",
    "updated_at": "2025-06-28T11:05:00Z",
    "confirmed_at": "2025-06-28T11:05:00Z"
  }
}
```

**What Changes After Confirmation:**
- `status` changes from "speculative" to "confirmed"
- `confirmed_at` gets a timestamp
- The startup's actual risk profile is updated
- This decision now affects future risk calculations

**Flutter Usage:**
```dart
// Show confirmation dialog before confirming
showDialog(
  context: context,
  builder: (context) => AlertDialog(
    title: Text('Confirm Decision'),
    content: Text('This will update your risk profile. Are you sure?'),
    actions: [
      TextButton(onPressed: () => Navigator.pop(context), child: Text('Cancel')),
      ElevatedButton(onPressed: () => confirmDecision(), child: Text('Confirm')),
    ],
  ),
);
```

**Error Responses:**
- **400 Bad Request:** Invalid decision ID
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** Decision doesn't belong to user's startup
- **404 Not Found:** Decision not found
- **500 Internal Error:** Server problem

---

### 7. Get All Decisions (Decision History)

**Endpoint:** `GET /api/v1/decisions/`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** Get all decisions (both speculative and confirmed) for the user's startup. Perfect for building a "Decision History" screen.

**Request:** No body needed, just JWT token in headers

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Decisions retrieved successfully",
  "data": [
    {
      "id": "880e8400-e29b-41d4-a716-446655440000",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "description": "Hire 3 senior engineers to accelerate product development",
      "category": "hiring",
      "status": "confirmed",
      "previous_risk_score": 6.8,
      "projected_risk_score": 7.2,
      "risk_delta": 0.4,
      "confidence": 0.82,
      "created_at": "2025-06-28T11:00:00Z",
      "confirmed_at": "2025-06-28T11:05:00Z"
    },
    {
      "id": "990e8400-e29b-41d4-a716-446655440000",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "description": "Launch paid marketing campaign on Google Ads",
      "category": "marketing",
      "status": "speculative",
      "previous_risk_score": 7.2,
      "projected_risk_score": 6.9,
      "risk_delta": -0.3,
      "confidence": 0.75,
      "created_at": "2025-06-28T12:00:00Z",
      "confirmed_at": null
    }
  ]
}
```

**Flutter Usage Ideas:**
```dart
// Group decisions by status
List<Decision> confirmedDecisions = decisions.where((d) => d.status == 'confirmed').toList();
List<Decision> speculativeDecisions = decisions.where((d) => d.status == 'speculative').toList();

// Show with different visual indicators
Card(
  color: decision.status == 'confirmed' ? Colors.green[50] : Colors.orange[50],
  child: ListTile(
    leading: Icon(
      decision.status == 'confirmed' ? Icons.check_circle : Icons.help_outline,
      color: decision.status == 'confirmed' ? Colors.green : Colors.orange,
    ),
    title: Text(decision.description),
    subtitle: Text('Risk ${decision.riskDelta > 0 ? 'increased' : 'decreased'} by ${decision.riskDelta.abs()}'),
    trailing: Text('${(decision.confidence * 100).toInt()}% confident'),
  ),
);
```

**Error Responses:**
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** User has no startup
- **500 Internal Error:** Server problem

---

### 8. Get Specific Decision (Decision Details)

**Endpoint:** `GET /api/v1/decisions/{decision_id}`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** Get detailed information about a specific decision. Use this for a "Decision Details" screen.

**URL Parameters:**
- `decision_id` (string): UUID of the decision

**Example:** `GET /api/v1/decisions/880e8400-e29b-41d4-a716-446655440000`

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Decision retrieved successfully",
  "data": {
    "id": "880e8400-e29b-41d4-a716-446655440000",
    "startup_id": "660e8400-e29b-41d4-a716-446655440000",
    "description": "Hire 3 senior engineers to accelerate product development",
    "category": "hiring",
    "context": "Current team is 5 people, need to launch MVP in 3 months. Competition is heating up.",
    "timeline": "2 months",
    "budget": 180000,
    "status": "confirmed",
    "previous_risk_score": 6.8,
    "projected_risk_score": 7.2,
    "risk_delta": 0.4,
    "confidence": 0.82,
    "suggestions": [
      "Ensure new hires have relevant EdTech experience",
      "Consider contracting before full-time hiring",
      "Budget for additional infrastructure costs",
      "Plan for longer onboarding time given tight timeline"
    ],
    "reasoning": "Hiring senior engineers will accelerate development but increases burn rate and management complexity. The timeline pressure justifies the risk increase. Strong technical leadership will be crucial for success.",
    "created_at": "2025-06-28T11:00:00Z",
    "updated_at": "2025-06-28T11:05:00Z",
    "confirmed_at": "2025-06-28T11:05:00Z"
  }
}
```

**Flutter Usage:**
Use this for a detailed view when user taps on a decision from the list. Show all the AI analysis, suggestions, and reasoning.

**Error Responses:**
- **400 Bad Request:** Invalid decision ID format
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** Decision doesn't belong to user's startup
- **404 Not Found:** Decision not found
- **500 Internal Error:** Server problem

---

## 📊 Risk Assessment Endpoints

> **Flutter Implementation Tip:** These endpoints power your main dashboard! Show beautiful risk scores, charts, and trend analysis. Perfect for data visualization widgets.

*All risk endpoints require both authentication and startup context.*

### 9. Get Current Risk Profile (Main Dashboard Data)

**Endpoint:** `GET /api/v1/risk/current`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** Get the current AI-powered risk assessment for the user's startup. This is perfect for your main dashboard screen - show the current risk score, level, and AI suggestions.

**Request:** No body needed, just JWT token in headers

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Risk profile retrieved successfully",
  "data": {
    "id": "770e8400-e29b-41d4-a716-446655440000",
    "startup_id": "660e8400-e29b-41d4-a716-446655440000",
    "score": 7.2,
    "level": "medium",
    "confidence": 0.85,
    "factors": [
      "Increased team size adds management complexity",
      "Strong technical team with relevant experience",
      "Competitive market with established players",
      "Higher burn rate due to recent hiring decisions",
      "AI technology provides clear differentiation"
    ],
    "suggestions": [
      "Implement agile management practices for larger team",
      "Focus on unique AI differentiation in marketing",
      "Monitor burn rate closely with new hires",
      "Develop partnerships with educational institutions early"
    ],
    "reasoning": "Recent hiring decision has increased operational complexity and burn rate, leading to slight risk increase. However, stronger technical team should accelerate product development and time-to-market.",
    "created_at": "2025-06-28T11:05:00Z",
    "updated_at": "2025-06-28T11:05:00Z"
  }
}
```

**What Each Field Means for Flutter:**

**Risk Score & Level:**
- `score`: Risk score from 0.0 to 10.0 (lower is better)
- `level`: Text representation of risk level:
  - `low` (0.0 - 3.0): Green - Minimal risk, strong position
  - `medium` (3.1 - 7.0): Yellow - Moderate risk, manageable with attention  
  - `high` (7.1 - 9.0): Orange - Significant risk, requires immediate action
  - `critical` (9.1 - 10.0): Red - Extreme risk, major changes needed

**AI Analysis:**
- `confidence`: How confident the AI is (0.85 = 85% confident)
- `factors`: Array of things affecting the risk (both positive and negative)
- `suggestions`: Array of actionable advice from AI
- `reasoning`: AI's explanation of the overall assessment

**Flutter UI Ideas:**
```dart
// Risk Score Gauge/Progress Indicator
class RiskScoreWidget extends StatelessWidget {
  final double score;
  final String level;
  
  Color getRiskColor() {
    switch (level) {
      case 'low': return Colors.green;
      case 'medium': return Colors.yellow[700]!;
      case 'high': return Colors.orange;
      case 'critical': return Colors.red;
      default: return Colors.grey;
    }
  }
  
  Widget build(BuildContext context) {
    return Card(
      child: Column(
        children: [
          CircularProgressIndicator(
            value: score / 10.0,
            backgroundColor: Colors.grey[300],
            valueColor: AlwaysStoppedAnimation<Color>(getRiskColor()),
            strokeWidth: 8,
          ),
          Text('${score.toStringAsFixed(1)}/10'),
          Text(level.toUpperCase(), style: TextStyle(color: getRiskColor())),
        ],
      ),
    );
  }
}

// Suggestions List
class SuggestionsList extends StatelessWidget {
  final List<String> suggestions;
  
  Widget build(BuildContext context) {
    return Column(
      children: suggestions.map((suggestion) => Card(
        child: ListTile(
          leading: Icon(Icons.lightbulb, color: Colors.amber),
          title: Text(suggestion),
        ),
      )).toList(),
    );
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** User has no startup
- **404 Not Found:** Risk profile not found (system error)
- **500 Internal Error:** Server problem

---

### 10. Get Risk Evolution History (Risk Chart Data)

**Endpoint:** `GET /api/v1/risk/history`

**Authentication:** Required (JWT Token + Startup Context)

**Description:** Get the historical risk scores to show how risk has changed over time. Perfect for charts and graphs showing risk trends.

**Query Parameters:**
- `limit` (optional, integer): Number of records to return (default: 10, max: 100)

**Example:** `GET /api/v1/risk/history?limit=20`

**Success Response:** (Status Code 200)
```json
{
  "success": true,
  "message": "Risk evolution retrieved successfully",
  "data": [
    {
      "id": "aa0e8400-e29b-41d4-a716-446655440000",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "score": 7.2,
      "level": "medium",
      "trigger": "Decision confirmed: Hire 3 senior engineers to accelerate product development",
      "created_at": "2025-06-28T11:05:00Z"
    },
    {
      "id": "bb0e8400-e29b-41d4-a716-446655440000",
      "startup_id": "660e8400-e29b-41d4-a716-446655440000",
      "score": 6.8,
      "level": "medium",
      "trigger": "Initial risk assessment after startup onboarding",
      "created_at": "2025-06-28T10:30:00Z"
    }
  ]
}
```

**What Each Field Means:**
- `score`: Risk score at that point in time
- `level`: Risk level at that point in time
- `trigger`: What caused this risk assessment (decision, onboarding, etc.)
- `created_at`: When this risk assessment was created

**Flutter Chart Ideas:**
```dart
// Using fl_chart package for line chart
class RiskEvolutionChart extends StatelessWidget {
  final List<RiskEvolution> riskHistory;
  
  Widget build(BuildContext context) {
    return LineChart(
      LineChartData(
        gridData: FlGridData(show: true),
        titlesData: FlTitlesData(show: true),
        borderData: FlBorderData(show: true),
        minX: 0,
        maxX: riskHistory.length.toDouble() - 1,
        minY: 0,
        maxY: 10,
        lineBarsData: [
          LineChartBarData(
            spots: riskHistory.asMap().entries.map((entry) {
              return FlSpot(entry.key.toDouble(), entry.value.score);
            }).toList(),
            isCurved: true,
            colors: [Colors.blue],
            barWidth: 3,
            dotData: FlDotData(show: true),
          ),
        ],
      ),
    );
  }
}

// Simple list view showing risk changes
class RiskHistoryList extends StatelessWidget {
  final List<RiskEvolution> riskHistory;
  
  Widget build(BuildContext context) {
    return ListView.builder(
      itemCount: riskHistory.length,
      itemBuilder: (context, index) {
        final risk = riskHistory[index];
        return Card(
          child: ListTile(
            leading: CircleAvatar(
              backgroundColor: getRiskColor(risk.level),
              child: Text(risk.score.toStringAsFixed(1)),
            ),
            title: Text(risk.trigger),
            subtitle: Text('Risk Level: ${risk.level}'),
            trailing: Text(formatDate(risk.createdAt)),
          ),
        );
      },
    );
  }
}
```

**Error Responses:**
- **401 Unauthorized:** Authentication required
- **403 Forbidden:** User has no startup
- **500 Internal Error:** Server problem

---

## 🔧 Utility Endpoints

### 11. Health Check (Test API Connection)

**Endpoint:** `GET /health`

**Authentication:** Not required

**Description:** Simple endpoint to check if the API is running. Use this to test connectivity during app startup or debug network issues.

**Success Response:** (Status Code 200)
```json
{
  "status": "ok",
  "message": "Risk Assessment API is running"
}
```

**Flutter Usage:**
```dart
// Check API connectivity during app startup
Future<bool> checkApiHealth() async {
  try {
    final response = await http.get(Uri.parse('$baseUrl/health'));
    return response.statusCode == 200;
  } catch (e) {
    return false;
  }
}

// Show status in debug screen
Text(isApiHealthy ? '✅ API Connected' : '❌ API Disconnected');
```

---

## 🎯 User Flow Recommendations for Flutter Developers

### 1. App Startup Flow (What happens when app opens)
```dart
// 1. Check for stored JWT token
final prefs = await SharedPreferences.getInstance();
final token = prefs.getString('jwt_token');

if (token != null) {
  // 2. Validate token by checking user profile
  try {
    final response = await apiClient.get('/api/v1/startup/profile');
    // User is logged in and has completed onboarding
    Navigator.pushReplacement(context, MaterialPageRoute(
      builder: (context) => DashboardScreen()
    ));
  } catch (e) {
    if (e.statusCode == 400) {
      // User is logged in but needs onboarding
      Navigator.pushReplacement(context, MaterialPageRoute(
        builder: (context) => OnboardingScreen()
      ));
    } else {
      // Token is invalid, go to login
      Navigator.pushReplacement(context, MaterialPageRoute(
        builder: (context) => LoginScreen()
      ));
    }
  }
} else {
  // No token, show login/signup
  Navigator.pushReplacement(context, MaterialPageRoute(
    builder: (context) => LoginScreen()
  ));
}
```

### 2. Recommended Screen Structure

**Login/Signup Screens:**
- Simple forms using `/api/v1/auth/signup` and `/api/v1/auth/login`
- Store JWT token immediately after successful authentication
- Navigate based on whether user has `startup_id` in response

**Onboarding Screen/Flow:**
- Break the large startup onboarding form into multiple screens
- Save progress locally in case user closes app
- Use `/api/v1/startup/onboard` when all data is collected
- Show loading indicator while AI processes the data

**Main Dashboard:**
- Use `/api/v1/risk/current` to show current risk score and suggestions
- Use `/api/v1/risk/history` for risk trend chart
- Quick "Simulate Decision" button leading to decision speculation

**Decision Simulator:**
- Form to collect decision details
- Use `/api/v1/decisions/speculate` to get AI analysis
- Show risk impact clearly (before/after scores)
- Allow user to save as speculative or confirm immediately

**Decision History:**
- Use `/api/v1/decisions/` to list all decisions
- Separate tabs for "Confirmed" and "Speculative" decisions
- Tap on decision for details using `/api/v1/decisions/{id}`

**Profile/Settings:**
- Use `/api/v1/startup/profile` to show current startup info
- Allow editing (though editing endpoints aren't implemented yet)
- Logout functionality (clear stored token)

### 3. Decision Flow (Core Feature)
```dart
// 1. User fills out decision form
final decisionData = {
  'description': 'Hire 5 developers',
  'category': 'hiring',
  'context': 'Need to scale product team',
  'timeline': '3 months',
  'budget': 300000,
};

// 2. Get AI analysis
final speculationResponse = await decisionService.speculateDecision(decisionData);
final decision = Decision.fromJson(speculationResponse);

// 3. Show results to user
showDialog(
  context: context,
  builder: (context) => DecisionResultsDialog(
    decision: decision,
    onConfirm: () => confirmDecision(decision.id),
    onSave: () => saveAsSpeculative(decision.id),
  ),
);

// 4. If user confirms, update their risk profile
if (userConfirms) {
  await decisionService.confirmDecision(decision.id);
  // Refresh dashboard to show new risk score
  Navigator.pushReplacement(context, MaterialPageRoute(
    builder: (context) => DashboardScreen()
  ));
}
```

### 4. Error Handling Strategy
```dart
class ApiErrorHandler {
  static void handle(BuildContext context, ApiException error) {
    String message;
    
    switch (error.statusCode) {
      case 400:
        message = 'Please check your input: ${error.message}';
        break;
      case 401:
        // Token expired or invalid
        _handleAuthError(context);
        return;
      case 403:
        message = 'Please complete startup onboarding first';
        Navigator.pushReplacement(context, MaterialPageRoute(
          builder: (context) => OnboardingScreen()
        ));
        return;
      case 404:
        message = 'Information not found';
        break;
      case 500:
        message = 'Server error. Please try again later.';
        break;
      default:
        message = 'Something went wrong. Please try again.';
    }
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text(message), backgroundColor: Colors.red),
    );
  }
  
  static void _handleAuthError(BuildContext context) {
    // Clear stored token
    SharedPreferences.getInstance().then((prefs) => prefs.remove('jwt_token'));
    
    // Navigate to login
    Navigator.pushAndRemoveUntil(
      context,
      MaterialPageRoute(builder: (context) => LoginScreen()),
      (route) => false,
    );
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(content: Text('Session expired. Please login again.')),
    );
  }
}
```

### 5. Data Refresh Strategy
```dart
// Refresh key data when user returns to dashboard
class DashboardScreen extends StatefulWidget {
  @override
  _DashboardScreenState createState() => _DashboardScreenState();
}

class _DashboardScreenState extends State<DashboardScreen> with RouteAware {
  RiskProfile? currentRisk;
  List<Decision> recentDecisions = [];
  bool isLoading = true;
  
  @override
  void initState() {
    super.initState();
    _loadDashboardData();
  }
  
  Future<void> _loadDashboardData() async {
    setState(() => isLoading = true);
    
    try {
      // Load data in parallel
      final futures = await Future.wait([
        riskService.getCurrentRisk(),
        decisionService.getDecisions(),
      ]);
      
      setState(() {
        currentRisk = futures[0] as RiskProfile;
        recentDecisions = (futures[1] as List<Decision>).take(5).toList();
        isLoading = false;
      });
    } catch (e) {
      ApiErrorHandler.handle(context, e);
      setState(() => isLoading = false);
    }
  }
  
  // Refresh when returning from other screens
  @override
  void didPopNext() {
    _loadDashboardData();
  }
}
```

### 6. Recommended Flutter Packages

**Essential packages for this API:**
```yaml
dependencies:
  http: ^0.13.5                 # HTTP requests
  shared_preferences: ^2.0.17   # Store JWT token
  fl_chart: ^0.62.0            # Charts for risk history
  provider: ^6.0.5             # State management
  
dev_dependencies:
  json_annotation: ^4.8.1      # JSON serialization
  build_runner: ^2.3.3         # Code generation
  json_serializable: ^6.6.2    # Generate fromJson/toJson
```

**Nice-to-have packages:**
```yaml
dependencies:
  flutter_secure_storage: ^9.0.0  # More secure token storage
  dio: ^5.1.2                      # More powerful HTTP client
  get_it: ^7.6.0                   # Dependency injection
  rxdart: ^0.27.7                  # Reactive programming
  cached_network_image: ^3.2.3     # Image caching
  shimmer: ^3.0.0                  # Loading animations
```

## 🏗️ Flutter Integration Guide

### 1. HTTP Client Setup

```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

class ApiClient {
  static const String baseUrl = 'http://localhost:8080';
  String? _token;

  void setToken(String token) {
    _token = token;
  }

  Map<String, String> get _headers {
    final headers = {
      'Content-Type': 'application/json',
    };
    
    if (_token != null) {
      headers['Authorization'] = 'Bearer $_token';
    }
    
    return headers;
  }

  Future<Map<String, dynamic>> post(String endpoint, Map<String, dynamic> data) async {
    final response = await http.post(
      Uri.parse('$baseUrl$endpoint'),
      headers: _headers,
      body: json.encode(data),
    );
    
    return _handleResponse(response);
  }

  Future<Map<String, dynamic>> get(String endpoint) async {
    final response = await http.get(
      Uri.parse('$baseUrl$endpoint'),
      headers: _headers,
    );
    
    return _handleResponse(response);
  }

  Map<String, dynamic> _handleResponse(http.Response response) {
    final Map<String, dynamic> data = json.decode(response.body);
    
    if (!data['success']) {
      throw ApiException(data['message'] ?? 'Unknown error', response.statusCode);
    }
    
    return data;
  }
}

class ApiException implements Exception {
  final String message;
  final int statusCode;
  
  ApiException(this.message, this.statusCode);
  
  @override
  String toString() => 'ApiException: $message (Status: $statusCode)';
}
```

### 2. Authentication Service

```dart
import 'package:shared_preferences/shared_preferences.dart';

class AuthService {
  final ApiClient _apiClient;
  static const String _tokenKey = 'jwt_token';
  
  AuthService(this._apiClient);

  Future<User> signup(String email, String name, String password) async {
    final response = await _apiClient.post('/api/v1/auth/signup', {
      'email': email,
      'name': name,
      'password': password,
    });
    
    final userData = response['data'];
    final token = userData['token'];
    
    await _saveToken(token);
    _apiClient.setToken(token);
    
    return User.fromJson(userData['user']);
  }

  Future<User> login(String email, String password) async {
    final response = await _apiClient.post('/api/v1/auth/login', {
      'email': email,
      'password': password,
    });
    
    final userData = response['data'];
    final token = userData['token'];
    
    await _saveToken(token);
    _apiClient.setToken(token);
    
    return User.fromJson(userData['user']);
  }

  Future<void> _saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
  }

  Future<bool> loadToken() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString(_tokenKey);
    
    if (token != null) {
      _apiClient.setToken(token);
      return true;
    }
    
    return false;
  }

  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    _apiClient.setToken('');
  }
}
```

### 3. Data Models

```dart
class User {
  final String id;
  final String email;
  final String name;
  final String role;
  final String? startupId;
  final DateTime createdAt;
  final DateTime updatedAt;

  User({
    required this.id,
    required this.email,
    required this.name,
    required this.role,
    this.startupId,
    required this.createdAt,
    required this.updatedAt,
  });

  factory User.fromJson(Map<String, dynamic> json) {
    return User(
      id: json['id'],
      email: json['email'],
      name: json['name'],
      role: json['role'],
      startupId: json['startup_id'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class Startup {
  final String id;
  final String name;
  final String description;
  final String industry;
  final String fundingStage;
  final String location;
  final DateTime foundedDate;
  final int teamSize;
  final String? website;
  final DateTime createdAt;
  final DateTime updatedAt;

  Startup({
    required this.id,
    required this.name,
    required this.description,
    required this.industry,
    required this.fundingStage,
    required this.location,
    required this.foundedDate,
    required this.teamSize,
    this.website,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Startup.fromJson(Map<String, dynamic> json) {
    return Startup(
      id: json['id'],
      name: json['name'],
      description: json['description'],
      industry: json['industry'],
      fundingStage: json['funding_stage'],
      location: json['location'],
      foundedDate: DateTime.parse(json['founded_date']),
      teamSize: json['team_size'],
      website: json['website'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class RiskProfile {
  final String id;
  final String startupId;
  final double score;
  final String level;
  final double confidence;
  final List<String> factors;
  final List<String> suggestions;
  final String reasoning;
  final DateTime createdAt;
  final DateTime updatedAt;

  RiskProfile({
    required this.id,
    required this.startupId,
    required this.score,
    required this.level,
    required this.confidence,
    required this.factors,
    required this.suggestions,
    required this.reasoning,
    required this.createdAt,
    required this.updatedAt,
  });

  factory RiskProfile.fromJson(Map<String, dynamic> json) {
    return RiskProfile(
      id: json['id'],
      startupId: json['startup_id'],
      score: json['score'].toDouble(),
      level: json['level'],
      confidence: json['confidence'].toDouble(),
      factors: List<String>.from(json['factors']),
      suggestions: List<String>.from(json['suggestions']),
      reasoning: json['reasoning'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
    );
  }
}

class Decision {
  final String id;
  final String startupId;
  final String description;
  final String category;
  final String? context;
  final String? timeline;
  final double? budget;
  final String status;
  final double previousRiskScore;
  final double projectedRiskScore;
  final double riskDelta;
  final double confidence;
  final List<String> suggestions;
  final String reasoning;
  final DateTime createdAt;
  final DateTime updatedAt;
  final DateTime? confirmedAt;

  Decision({
    required this.id,
    required this.startupId,
    required this.description,
    required this.category,
    this.context,
    this.timeline,
    this.budget,
    required this.status,
    required this.previousRiskScore,
    required this.projectedRiskScore,
    required this.riskDelta,
    required this.confidence,
    required this.suggestions,
    required this.reasoning,
    required this.createdAt,
    required this.updatedAt,
    this.confirmedAt,
  });

  factory Decision.fromJson(Map<String, dynamic> json) {
    return Decision(
      id: json['id'],
      startupId: json['startup_id'],
      description: json['description'],
      category: json['category'],
      context: json['context'],
      timeline: json['timeline'],
      budget: json['budget']?.toDouble(),
      status: json['status'],
      previousRiskScore: json['previous_risk_score'].toDouble(),
      projectedRiskScore: json['projected_risk_score'].toDouble(),
      riskDelta: json['risk_delta'].toDouble(),
      confidence: json['confidence'].toDouble(),
      suggestions: List<String>.from(json['suggestions']),
      reasoning: json['reasoning'],
      createdAt: DateTime.parse(json['created_at']),
      updatedAt: DateTime.parse(json['updated_at']),
      confirmedAt: json['confirmed_at'] != null ? DateTime.parse(json['confirmed_at']) : null,
    );
  }
}
```

### 4. Service Classes

```dart
class StartupService {
  final ApiClient _apiClient;
  
  StartupService(this._apiClient);

  Future<Map<String, dynamic>> onboardStartup(Map<String, dynamic> startupData) async {
    final response = await _apiClient.post('/api/v1/startup/onboard', startupData);
    return response['data'];
  }

  Future<Startup> getProfile() async {
    final response = await _apiClient.get('/api/v1/startup/profile');
    return Startup.fromJson(response['data']);
  }
}

class RiskService {
  final ApiClient _apiClient;
  
  RiskService(this._apiClient);

  Future<RiskProfile> getCurrentRisk() async {
    final response = await _apiClient.get('/api/v1/risk/current');
    return RiskProfile.fromJson(response['data']);
  }

  Future<List<RiskEvolution>> getRiskHistory({int limit = 10}) async {
    final response = await _apiClient.get('/api/v1/risk/history?limit=$limit');
    final List<dynamic> data = response['data'];
    return data.map((item) => RiskEvolution.fromJson(item)).toList();
  }
}

class DecisionService {
  final ApiClient _apiClient;
  
  DecisionService(this._apiClient);

  Future<Decision> speculateDecision(Map<String, dynamic> decisionData) async {
    final response = await _apiClient.post('/api/v1/decisions/speculate', decisionData);
    return Decision.fromJson(response['data']);
  }

  Future<Decision> confirmDecision(String decisionId, {String? notes}) async {
    final data = {'decision_id': decisionId};
    if (notes != null) data['notes'] = notes;
    
    final response = await _apiClient.post('/api/v1/decisions/confirm', data);
    return Decision.fromJson(response['data']);
  }

  Future<List<Decision>> getDecisions() async {
    final response = await _apiClient.get('/api/v1/decisions/');
    final List<dynamic> data = response['data'];
    return data.map((item) => Decision.fromJson(item)).toList();
  }

  Future<Decision> getDecision(String decisionId) async {
    final response = await _apiClient.get('/api/v1/decisions/$decisionId');
    return Decision.fromJson(response['data']);
  }
}
```

### 5. Error Handling

```dart
// Add this to your widgets
class ErrorHandler {
  static void handleApiException(BuildContext context, ApiException e) {
    String message;
    
    switch (e.statusCode) {
      case 400:
        message = 'Invalid input: ${e.message}';
        break;
      case 401:
        message = 'Authentication required. Please login again.';
        // Navigate to login screen
        Navigator.of(context).pushReplacementNamed('/login');
        return;
      case 403:
        message = 'Access denied: ${e.message}';
        break;
      case 404:
        message = 'Resource not found: ${e.message}';
        break;
      case 500:
        message = 'Server error. Please try again later.';
        break;
      default:
        message = 'An unexpected error occurred: ${e.message}';
    }
    
    ScaffoldMessenger.of(context).showSnackBar(
      SnackBar(
        content: Text(message),
        backgroundColor: Colors.red,
      ),
    );
  }
}
```

---

## 🎯 User Flow Recommendations

### 1. App Startup Flow
1. Check for stored JWT token
2. If token exists, validate with API call
3. Navigate to appropriate screen (onboarding vs dashboard)

### 2. Onboarding Flow
1. User signup/login
2. Startup profile creation (comprehensive form)
3. Wait for AI analysis completion
4. Show risk assessment results
5. Navigate to main dashboard

### 3. Decision Flow
1. User describes potential decision
2. AI provides risk analysis and suggestions
3. User can save as speculative or confirm immediately
4. Confirmed decisions update risk profile
5. Show updated risk metrics

### 4. Dashboard Elements
- Current risk score with visual indicator
- Recent decisions (speculative and confirmed)
- Risk evolution chart
- Quick decision speculation button
- Suggestions from AI analysis

---

## � Quick Start Guide for Absolute Beginners

### Step 1: Set Up Your Flutter Project

1. **Create a new Flutter project:**
```bash
flutter create risq_app
cd risq_app
```

2. **Add required dependencies to `pubspec.yaml`:**
```yaml
dependencies:
  flutter:
    sdk: flutter
  http: ^0.13.5
  shared_preferences: ^2.0.17
  fl_chart: ^0.62.0
  provider: ^6.0.5

dev_dependencies:
  flutter_test:
    sdk: flutter
  flutter_lints: ^2.0.0
```

3. **Run `flutter pub get` to install packages**

### Step 2: Create Your API Client (Copy and Paste This!)

Create `lib/services/api_client.dart`:
```dart
import 'package:http/http.dart' as http;
import 'dart:convert';

class ApiClient {
  static const String baseUrl = 'http://localhost:8080'; // Change for production
  String? _token;

  void setToken(String token) {
    _token = token;
  }

  Map<String, String> get _headers {
    final headers = {
      'Content-Type': 'application/json',
    };
    
    if (_token != null) {
      headers['Authorization'] = 'Bearer $_token';
    }
    
    return headers;
  }

  Future<Map<String, dynamic>> post(String endpoint, Map<String, dynamic> data) async {
    final response = await http.post(
      Uri.parse('$baseUrl$endpoint'),
      headers: _headers,
      body: json.encode(data),
    );
    
    return _handleResponse(response);
  }

  Future<Map<String, dynamic>> get(String endpoint) async {
    final response = await http.get(
      Uri.parse('$baseUrl$endpoint'),
      headers: _headers,
    );
    
    return _handleResponse(response);
  }

  Map<String, dynamic> _handleResponse(http.Response response) {
    final Map<String, dynamic> data = json.decode(response.body);
    
    if (!data['success']) {
      throw ApiException(data['message'] ?? 'Unknown error', response.statusCode);
    }
    
    return data;
  }
}

class ApiException implements Exception {
  final String message;
  final int statusCode;
  
  ApiException(this.message, this.statusCode);
  
  @override
  String toString() => 'ApiException: $message (Status: $statusCode)';
}
```

### Step 3: Create Your Authentication Service

Create `lib/services/auth_service.dart`:
```dart
import 'package:shared_preferences/shared_preferences.dart';
import 'api_client.dart';

class AuthService {
  final ApiClient _apiClient;
  static const String _tokenKey = 'jwt_token';
  
  AuthService(this._apiClient);

  Future<Map<String, dynamic>> signup(String email, String name, String password) async {
    final response = await _apiClient.post('/api/v1/auth/signup', {
      'email': email,
      'name': name,
      'password': password,
    });
    
    final userData = response['data'];
    final token = userData['token'];
    
    await _saveToken(token);
    _apiClient.setToken(token);
    
    return userData;
  }

  Future<Map<String, dynamic>> login(String email, String password) async {
    final response = await _apiClient.post('/api/v1/auth/login', {
      'email': email,
      'password': password,
    });
    
    final userData = response['data'];
    final token = userData['token'];
    
    await _saveToken(token);
    _apiClient.setToken(token);
    
    return userData;
  }

  Future<void> _saveToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_tokenKey, token);
  }

  Future<bool> loadToken() async {
    final prefs = await SharedPreferences.getInstance();
    final token = prefs.getString(_tokenKey);
    
    if (token != null) {
      _apiClient.setToken(token);
      return true;
    }
    
    return false;
  }

  Future<void> logout() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove(_tokenKey);
    _apiClient.setToken('');
  }
}
```

### Step 4: Create a Simple Login Screen

Create `lib/screens/login_screen.dart`:
```dart
import 'package:flutter/material.dart';
import '../services/api_client.dart';
import '../services/auth_service.dart';

class LoginScreen extends StatefulWidget {
  @override
  _LoginScreenState createState() => _LoginScreenState();
}

class _LoginScreenState extends State<LoginScreen> {
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _apiClient = ApiClient();
  late final AuthService _authService;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _authService = AuthService(_apiClient);
  }

  Future<void> _login() async {
    setState(() => _isLoading = true);
    
    try {
      final userData = await _authService.login(
        _emailController.text,
        _passwordController.text,
      );
      
      // Check if user has completed onboarding
      if (userData['user']['startup_id'] != null) {
        // Navigate to dashboard
        Navigator.pushReplacementNamed(context, '/dashboard');
      } else {
        // Navigate to onboarding
        Navigator.pushReplacementNamed(context, '/onboarding');
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Login failed: $e')),
      );
    } finally {
      setState(() => _isLoading = false);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('RISQ Login')),
      body: Padding(
        padding: EdgeInsets.all(16.0),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            TextField(
              controller: _emailController,
              decoration: InputDecoration(labelText: 'Email'),
              keyboardType: TextInputType.emailAddress,
            ),
            SizedBox(height: 16),
            TextField(
              controller: _passwordController,
              decoration: InputDecoration(labelText: 'Password'),
              obscureText: true,
            ),
            SizedBox(height: 24),
            _isLoading
                ? CircularProgressIndicator()
                : ElevatedButton(
                    onPressed: _login,
                    child: Text('Login'),
                  ),
          ],
        ),
      ),
    );
  }
}
```

### Step 5: Test Your API Connection

Create `lib/screens/test_screen.dart`:
```dart
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

class TestScreen extends StatefulWidget {
  @override
  _TestScreenState createState() => _TestScreenState();
}

class _TestScreenState extends State<TestScreen> {
  String _status = 'Tap to test API connection';

  Future<void> _testConnection() async {
    try {
      final response = await http.get(Uri.parse('http://localhost:8080/health'));
      if (response.statusCode == 200) {
        setState(() => _status = '✅ API is working! Response: ${response.body}');
      } else {
        setState(() => _status = '❌ API error: ${response.statusCode}');
      }
    } catch (e) {
      setState(() => _status = '❌ Connection failed: $e');
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(title: Text('API Test')),
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Text(_status, textAlign: TextAlign.center),
            SizedBox(height: 20),
            ElevatedButton(
              onPressed: _testConnection,
              child: Text('Test API Connection'),
            ),
          ],
        ),
      ),
    );
  }
}
```

### Step 6: Update Your main.dart

Replace your `lib/main.dart`:
```dart
import 'package:flutter/material.dart';
import 'screens/login_screen.dart';
import 'screens/test_screen.dart';

void main() {
  runApp(MyApp());
}

class MyApp extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'RISQ App',
      theme: ThemeData(
        primarySwatch: Colors.blue,
      ),
      home: TestScreen(), // Start with test screen, then change to LoginScreen()
      routes: {
        '/login': (context) => LoginScreen(),
        '/test': (context) => TestScreen(),
      },
    );
  }
}
```

### Step 7: Test Everything

1. **Start the backend server:**
   - Make sure the Docker containers are running
   - API should be available at `http://localhost:8080`

2. **Run your Flutter app:**
```bash
flutter run
```

3. **Test the connection:**
   - Tap "Test API Connection" - you should see "✅ API is working!"
   - If it fails, check that the backend is running on port 8080

4. **Test user signup:**
   - Change `home: TestScreen()` to `home: LoginScreen()` in main.dart
   - Try creating a new account using the signup endpoints

### Next Steps:

1. **Add more screens:** Dashboard, Onboarding, Decision Simulator
2. **Add state management:** Use Provider or Bloc for managing app state
3. **Add error handling:** Better error messages and retry logic
4. **Add UI polish:** Beautiful charts, animations, loading states
5. **Add offline support:** Cache important data locally

### Common Issues & Solutions:

**❌ "Connection refused" error:**
- Make sure backend is running: `docker-compose up`
- Check if using correct IP (localhost vs your machine's IP)

**❌ "CORS error" (web only):**
- Backend already handles CORS, but use mobile/desktop for testing

**❌ "Invalid token" errors:**
- Token might be expired, implement auto-refresh or re-login

**❌ JSON parsing errors:**
- Check API response format matches your model classes

This should get you started! The backend is already fully functional, so focus on building a great Flutter UI that calls these endpoints. Start simple with basic screens, then add more features gradually.

---

## �🔒 Security Considerations

1. **JWT Token Management**
   - Store tokens securely using flutter_secure_storage
   - Implement token refresh mechanism
   - Handle token expiration gracefully

2. **Input Validation**
   - Validate all user inputs on client side
   - Handle API validation errors properly
   - Sanitize sensitive data before storage

3. **Error Handling**
   - Never expose sensitive error information to users
   - Log errors securely for debugging
   - Implement proper retry mechanisms

4. **Network Security**
   - Use HTTPS in production
   - Implement certificate pinning
   - Handle network timeouts and failures

---

## 🚀 Production Deployment Notes

- Change base URL to production domain
- Update CORS settings for production domain
- Use environment variables for API endpoints
- Implement proper logging and analytics
- Add offline support for critical features
- Implement proper state management (Provider/Bloc)

This comprehensive documentation should provide your Flutter developer with everything needed to integrate with the RISQ backend API successfully.

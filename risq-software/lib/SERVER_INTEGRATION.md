# Server Integration Configuration

## Current Status
The onboarding form is currently configured to work **without a server**. The data is logged to the console and a success message is shown to the user.

## How to Enable Server Integration

When your Railway server is deployed and ready, follow these simple steps:

### 1. Update Configuration in `onboarding_page.dart`

```dart
// Change these constants at the top of _OnboardingPageState class:
static const bool _useServerAPI = true;  // Change from false to true
static const String _serverEndpoint = 'https://your-actual-railway-app.railway.app/api/startups';  // Update with your real endpoint
```

### 2. Server Endpoint Requirements

Your server should accept a POST request to the `/api/startups` endpoint with the following JSON structure:

```json
{
  "name": "Company Name",
  "description": "Company description",
  "industry": "Technology",
  "sector": "FinTech",
  "funding_stage": "Seed",
  "location": "San Francisco, CA",
  "founded_date": "2025-06-28T10:30:00.000Z",
  "team_size": 5,
  "website": "https://company.com",
  "business_model": "B2B SaaS",
  "revenue_streams": ["Monthly Subscription", "Enterprise Licenses"],
  "target_market": "Small to medium businesses",
  "competitive_advantage": "AI-powered automation",
  "implementation_plan": "Phase 1: MVP, Phase 2: Scale",
  "technology_stack": ["React", "Node.js", "PostgreSQL"],
  "development_timeline": "12 months to market",
  "go_to_market_strategy": "Direct sales and partnerships",
  "initial_investment": 500000.0,
  "monthly_burn_rate": 45000.0,
  "projected_revenue": 1200000.0,
  "funding_requirement": 2000000.0,
  "founder_details": [
    {
      "name": "John Doe",
      "email": "john@company.com",
      "role": "CEO & Founder",
      "linkedin_url": "https://linkedin.com/in/johndoe",
      "education": [
        {
          "degree": "Bachelor of Computer Science",
          "institution": "MIT",
          "graduation_year": "2020"
        }
      ],
      "experience": [
        {
          "company": "Google",
          "position": "Senior Software Engineer",
          "start_date": "2020-01-15",
          "end_date": "2023-12-31",
          "description": "Led development of AI features"
        }
      ],
      "skills": ["AI/ML", "Product Strategy", "Team Leadership"],
      "achievements": ["Built AI platform serving 1M+ users", "Led team of 10 engineers"]
    }
  ]
}
```

### 3. Expected Server Response

The server should respond with:
- **Status Code**: 200 or 201 for success
- **Response Body**: Can be any JSON (the app doesn't parse it currently)

For errors, return appropriate HTTP status codes (400, 500, etc.) and the app will show a generic error message.

### 4. Headers Required

The app sends these headers:
- `Content-Type: application/json`

Make sure your server accepts these headers and handles CORS if needed.

## Current Behavior (Server Disabled)

With `_useServerAPI = false`:
1. Form validates all inputs
2. Generates complete JSON payload
3. Logs the payload to console (check Flutter logs/debug console)
4. Navigates to a data display page showing all submitted information after 1.5 second delay
5. User can review their submitted data in a formatted, readable layout

## Navigation

The app now navigates to a dedicated data display page (`DataDisplayPage`) after successful submission. This page shows:
- All company information in organized sections
- Business model and market data
- Technical implementation details
- Financial information
- Complete founder details including education and experience
- Action buttons to return home or continue to the main app

When the server is enabled, the same navigation will occur after successful server response.

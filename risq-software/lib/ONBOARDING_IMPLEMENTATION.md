# Multi-Step Onboarding Implementation

## Overview
The onboarding page has been completely refactored into a comprehensive 4-step startup registration form that collects all the data matching your JSON structure and sends it to your railway server.

## Changes Made

### 1. Multi-Step Form Structure
- **Step 1: Company Information**
  - Company Name, Description, Industry, Sector, Location, Website

- **Step 2: Business Model & Market**
  - Business Model, Target Market, Competitive Advantage, Go-to-Market Strategy
  - Funding Stage (dropdown), Team Size (slider), Revenue Streams

- **Step 3: Technical & Implementation**
  - Implementation Plan, Development Timeline, Technology Stack

- **Step 4: Financial & Founder Information**
  - Financial: Initial Investment, Monthly Burn Rate, Projected Revenue, Funding Requirement
  - Founder: Name, Email, Role, LinkedIn, Education, Experience, Skills, Achievements

### 2. Features Implemented
- **Progress Indicator**: Shows current step and overall progress
- **Form Validation**: Each step has comprehensive validation
- **Step Navigation**: Next/Previous buttons with form validation
- **Responsive Design**: Adapts to different screen sizes
- **Video Background**: Maintained from original design
- **Data Collection**: All form data is structured into the JSON format you specified

### 3. API Integration
The form submits data to your railway server with the exact JSON structure:

```json
{
  "name": "Company Name",
  "description": "Business description...",
  "industry": "Technology", 
  "sector": "FinTech",
  "funding_stage": "Seed",
  "location": "San Francisco, CA",
  "founded_date": "2024-01-15T00:00:00Z",
  "team_size": 5,
  "website": "https://company.com",
  "business_model": "B2B SaaS",
  "revenue_streams": ["Monthly Subscription", "Enterprise Licenses"],
  "target_market": "SMBs in retail and e-commerce",
  "competitive_advantage": "Proprietary AI algorithms...",
  "implementation_plan": "Phase 1: MVP development...",
  "technology_stack": ["Go", "React", "PostgreSQL"],
  "development_timeline": "12 months to full market release",
  "go_to_market_strategy": "Direct sales to SMBs...",
  "initial_investment": 500000,
  "monthly_burn_rate": 45000,
  "projected_revenue": 1200000,
  "funding_requirement": 2000000,
  "founder_details": [
    {
      "name": "John Smith",
      "email": "john@company.com",
      "role": "CEO & Founder",
      "linkedin_url": "https://linkedin.com/in/johnsmith",
      "education": ["MBA from Stanford", "BS from MIT"],
      "experience": [{"description": "Led development at Google..."}],
      "skills": ["AI/ML", "Product Strategy", "Team Leadership"],
      "achievements": ["Improved targeting by 40%", "Published 5 papers"]
    }
  ]
}
```

## To Configure for Your Railway Server

### 1. Update the API Endpoint
In `_submitOnboarding()` method, replace the placeholder URL:

```dart
final response = await http.post(
  Uri.parse('https://your-railway-app.railway.app/api/startups'), // <- Update this
  headers: {
    'Content-Type': 'application/json',
  },
  body: json.encode(payload),
);
```

### 2. Add Authentication (if needed)
If your railway server requires authentication, add headers:

```dart
headers: {
  'Content-Type': 'application/json',
  'Authorization': 'Bearer $token',
  // Add other headers as needed
},
```

### 3. Error Handling
The current implementation shows generic error messages. You can customize based on your API responses:

```dart
if (response.statusCode == 200 || response.statusCode == 201) {
  // Success handling
} else {
  // Parse error response and show specific messages
  final errorData = json.decode(response.body);
  setState(() {
    _errorMessage = errorData['message'] ?? 'Failed to submit startup data';
  });
}
```

## Form Validation
Each step includes comprehensive validation:
- Required field validation
- Email format validation
- URL format validation
- Numeric field validation
- Minimum/maximum length validation

## User Experience
- **Progress Tracking**: Visual progress bar showing completion percentage
- **Step Labels**: Clear indication of current step (e.g., "Step 2 of 4")
- **Navigation**: Previous/Next buttons with form validation
- **Loading States**: Loading indicator during submission
- **Error Handling**: User-friendly error messages
- **Responsive Design**: Works on all screen sizes

## Testing
To test the form:
1. Fill out all required fields in each step
2. Use the Next button to progress through steps
3. Submit the final form to see the JSON payload in your server logs

## Next Steps
1. Deploy your railway server endpoint
2. Update the API URL in the `_submitOnboarding()` method
3. Test the complete flow
4. Add any additional validation or fields as needed
5. Implement navigation to the main app after successful submission

The form is now ready for production use and will send properly structured data to your railway server according to your specified JSON format.

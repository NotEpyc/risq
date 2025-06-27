package main

import (
	"encoding/json"
	"fmt"
	"risq_backend/internal/startup"
)

func main() {
	// Test JSON string
	jsonStr := `{
		"name": "Test Startup",
		"description": "A test startup for validation",
		"industry": "Technology",
		"sector": "fintech",
		"funding_stage": "seed",
		"location": "San Francisco, CA",
		"founded_date": "2024-01-15",
		"team_size": 3,
		"business_model": "SaaS subscription",
		"revenue_streams": ["Monthly subscriptions"],
		"target_market": "Small businesses",
		"implementation_plan": "Build MVP in 6 months",
		"initial_investment": 100000,
		"monthly_burn_rate": 15000,
		"projected_revenue": 200000,
		"funding_requirement": 500000,
		"founder_details": [{
			"name": "Test Founder",
			"email": "founder@test.com",
			"role": "CEO",
			"experience": [{
				"company": "Tech Corp",
				"position": "Developer",
				"duration": "2 years",
				"description": "Built web applications"
			}],
			"skills": ["Programming", "Leadership"],
			"achievements": ["Built successful side project"],
			"previous_startups": []
		}]
	}`

	var input startup.StartupOnboardingInput
	err := json.Unmarshal([]byte(jsonStr), &input)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		fmt.Printf("Error type: %T\n", err)
		return
	}

	fmt.Printf("Success! Parsed startup: %+v\n", input)
}

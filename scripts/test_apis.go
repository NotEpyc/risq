package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// APITestResult represents the result of an API test
type APITestResult struct {
	Service string `json:"service"`
	Status  string `json:"status"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

func main() {
	fmt.Println("🔧 Risk Assessment Backend - API Configuration Tester")
	fmt.Println("====================================================")
	fmt.Println()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		fmt.Println("⚠️  Warning: Could not load .env file, using system environment variables")
	}

	var results []APITestResult

	// Test NewsAPI
	fmt.Println("📰 Testing NewsAPI...")
	newsResult := testNewsAPI()
	results = append(results, newsResult)
	printResult(newsResult)
	fmt.Println()

	// Test MarketStack API
	fmt.Println("📊 Testing MarketStack API...")
	marketResult := testMarketStackAPI()
	results = append(results, marketResult)
	printResult(marketResult)
	fmt.Println()

	// Test OpenAI API
	fmt.Println("🤖 Testing OpenAI API...")
	openaiResult := testOpenAIAPI()
	results = append(results, openaiResult)
	printResult(openaiResult)
	fmt.Println()

	// Print summary
	printSummary(results)
}

func testNewsAPI() APITestResult {
	apiKey := os.Getenv("NEWS_API_KEY")
	if apiKey == "" || apiKey == "your_news_api_key_here" {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "❌ FAILED",
			Message: "API key not configured",
			Details: "Set NEWS_API_KEY in your .env file",
		}
	}

	// Test API call
	url := fmt.Sprintf("https://newsapi.org/v2/everything?q=technology&pageSize=1&apiKey=%s", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "❌ FAILED",
			Message: "Network error",
			Details: err.Error(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "❌ FAILED",
			Message: "Failed to read response",
			Details: err.Error(),
		}
	}

	var newsResponse struct {
		Status       string `json:"status"`
		TotalResults int    `json:"totalResults"`
		Code         string `json:"code"`
		Message      string `json:"message"`
	}

	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "❌ FAILED",
			Message: "Invalid response format",
			Details: err.Error(),
		}
	}

	if newsResponse.Status == "ok" {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "✅ SUCCESS",
			Message: fmt.Sprintf("API working! Found %d articles", newsResponse.TotalResults),
		}
	} else {
		return APITestResult{
			Service: "NewsAPI",
			Status:  "❌ FAILED",
			Message: "API error: " + newsResponse.Message,
			Details: "Code: " + newsResponse.Code,
		}
	}
}

func testMarketStackAPI() APITestResult {
	apiKey := os.Getenv("MARKET_DATA_API_KEY")
	if apiKey == "" || apiKey == "your_market_data_api_key_here" {
		return APITestResult{
			Service: "MarketStack",
			Status:  "❌ FAILED",
			Message: "API key not configured",
			Details: "Set MARKET_DATA_API_KEY in your .env file",
		}
	}

	// Test API call - using a simple endpoint that should work with free tier
	url := fmt.Sprintf("https://api.marketstack.com/v1/exchanges?access_key=%s&limit=1", apiKey)

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return APITestResult{
			Service: "MarketStack",
			Status:  "❌ FAILED",
			Message: "Network error",
			Details: err.Error(),
		}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return APITestResult{
			Service: "MarketStack",
			Status:  "❌ FAILED",
			Message: "Failed to read response",
			Details: err.Error(),
		}
	}

	var marketResponse struct {
		Data  []interface{} `json:"data"`
		Error struct {
			Code int    `json:"code"`
			Type string `json:"type"`
			Info string `json:"info"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &marketResponse); err != nil {
		return APITestResult{
			Service: "MarketStack",
			Status:  "❌ FAILED",
			Message: "Invalid response format",
			Details: err.Error(),
		}
	}

	if marketResponse.Error.Code != 0 {
		return APITestResult{
			Service: "MarketStack",
			Status:  "❌ FAILED",
			Message: "API error: " + marketResponse.Error.Info,
			Details: fmt.Sprintf("Code: %d, Type: %s", marketResponse.Error.Code, marketResponse.Error.Type),
		}
	}

	return APITestResult{
		Service: "MarketStack",
		Status:  "✅ SUCCESS",
		Message: fmt.Sprintf("API working! Retrieved %d data points", len(marketResponse.Data)),
	}
}

func testOpenAIAPI() APITestResult {
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" || apiKey == "sk-your_openai_api_key_here" {
		return APITestResult{
			Service: "OpenAI",
			Status:  "❌ FAILED",
			Message: "API key not configured",
			Details: "Set OPENAI_API_KEY in your .env file",
		}
	}

	// Test API call - using models endpoint which is lightweight
	url := "https://api.openai.com/v1/models"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return APITestResult{
			Service: "OpenAI",
			Status:  "❌ FAILED",
			Message: "Failed to create request",
			Details: err.Error(),
		}
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	resp, err := client.Do(req)
	if err != nil {
		return APITestResult{
			Service: "OpenAI",
			Status:  "❌ FAILED",
			Message: "Network error",
			Details: err.Error(),
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode == 401 {
		return APITestResult{
			Service: "OpenAI",
			Status:  "❌ FAILED",
			Message: "Invalid API key",
			Details: "Check your OpenAI API key",
		}
	}

	if resp.StatusCode != 200 {
		return APITestResult{
			Service: "OpenAI",
			Status:  "❌ FAILED",
			Message: fmt.Sprintf("HTTP %d error", resp.StatusCode),
		}
	}

	return APITestResult{
		Service: "OpenAI",
		Status:  "✅ SUCCESS",
		Message: "API working! Ready for AI analysis",
	}
}

func printResult(result APITestResult) {
	fmt.Printf("   %s %s: %s\n", result.Status, result.Service, result.Message)
	if result.Details != "" {
		fmt.Printf("   Details: %s\n", result.Details)
	}
}

func printSummary(results []APITestResult) {
	fmt.Println("📋 Summary")
	fmt.Println("==========")

	successCount := 0
	for _, result := range results {
		if result.Status == "✅ SUCCESS" {
			successCount++
		}
	}

	fmt.Printf("APIs configured: %d/%d\n\n", successCount, len(results))

	if successCount == len(results) {
		fmt.Println("🎉 All APIs are properly configured!")
		fmt.Println("   Your backend is ready to use real external data.")
		fmt.Println()
		fmt.Println("Next steps:")
		fmt.Println("   1. Start your backend: ./api")
		fmt.Println("   2. Test startup onboarding")
		fmt.Println("   3. Monitor logs for real API calls")
	} else {
		fmt.Println("⚠️  Some APIs need configuration:")
		for _, result := range results {
			if result.Status != "✅ SUCCESS" {
				fmt.Printf("   - %s: %s\n", result.Service, result.Message)
			}
		}
		fmt.Println()
		fmt.Println("📖 See EXTERNAL_API_SETUP.md for detailed setup instructions")
	}
}

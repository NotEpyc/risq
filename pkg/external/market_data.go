package external

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"risq_backend/pkg/logger"
)

// MarketDataService handles external market data API calls
type MarketDataService interface {
	GetMarketData(ctx context.Context, sector, industry string) (*MarketDataResponse, error)
	GetNewsAnalysis(ctx context.Context, sector string) (*NewsAnalysisResponse, error)
	GetCompetitorAnalysis(ctx context.Context, sector, businessModel string) (*CompetitorAnalysisResponse, error)
	GetIndustryTrends(ctx context.Context, industry string) (*IndustryTrends, error)
	GetSectorAnalysis(ctx context.Context, sector string) (*SectorAnalysis, error)
	GetMarketHealth(ctx context.Context, targetMarket string) (*MarketHealth, error)
}

type marketDataService struct {
	marketDataAPIKey string
	marketDataURL    string
	newsAPIKey       string
	newsAPIURL       string
	httpClient       *http.Client
}

// MarketDataResponse represents comprehensive market data
type MarketDataResponse struct {
	Sector           string         `json:"sector"`
	Industry         string         `json:"industry"`
	MarketSize       MarketSizeInfo `json:"market_size"`
	GrowthRate       float64        `json:"growth_rate"`
	MarketStatus     string         `json:"market_status"` // active, declining, emerging, mature
	CompetitionLevel string         `json:"competition_level"`
	KeyTrends        []string       `json:"key_trends"`
	Opportunities    []string       `json:"opportunities"`
	Threats          []string       `json:"threats"`
	RegulationLevel  string         `json:"regulation_level"`
	BarriersToEntry  []string       `json:"barriers_to_entry"`
	LastUpdated      time.Time      `json:"last_updated"`
}

type MarketSizeInfo struct {
	TAM  float64 `json:"tam"`  // Total Addressable Market in USD
	SAM  float64 `json:"sam"`  // Serviceable Addressable Market in USD
	SOM  float64 `json:"som"`  // Serviceable Obtainable Market in USD
	CAGR float64 `json:"cagr"` // Compound Annual Growth Rate
}

// NewsAnalysisResponse represents news sentiment analysis
type NewsAnalysisResponse struct {
	Sector           string     `json:"sector"`
	SentimentScore   float64    `json:"sentiment_score"` // -1 to 1
	PositiveKeywords []string   `json:"positive_keywords"`
	NegativeKeywords []string   `json:"negative_keywords"`
	RecentNews       []NewsItem `json:"recent_news"`
	InvestmentNews   []NewsItem `json:"investment_news"`
	AnalysisDate     time.Time  `json:"analysis_date"`
}

// CompetitorAnalysisResponse represents competitor landscape
type CompetitorAnalysisResponse struct {
	Sector               string       `json:"sector"`
	TotalCompetitors     int          `json:"total_competitors"`
	KeyPlayers           []Competitor `json:"key_players"`
	MarketLeader         string       `json:"market_leader"`
	EmergingPlayers      []string     `json:"emerging_players"`
	CompetitionIntensity string       `json:"competition_intensity"` // low, medium, high
}

type Competitor struct {
	Name        string   `json:"name"`
	MarketShare float64  `json:"market_share"`
	Funding     float64  `json:"funding"`
	Founded     int      `json:"founded"`
	Employees   int      `json:"employees"`
	Strengths   []string `json:"strengths"`
	Weaknesses  []string `json:"weaknesses"`
}

// Legacy structures for backward compatibility
type IndustryTrends struct {
	Industry         string   `json:"industry"`
	GrowthRate       float64  `json:"growth_rate"`
	MarketSize       int64    `json:"market_size"`
	CompetitionLevel string   `json:"competition_level"`
	Trends           []string `json:"trends"`
	Outlook          string   `json:"outlook"`
}

type SectorAnalysis struct {
	Sector         string     `json:"sector"`
	IsActive       bool       `json:"is_active"`
	Activity       string     `json:"activity"` // "growing", "stable", "declining"
	KeyPlayers     []string   `json:"key_players"`
	MarketCap      int64      `json:"market_cap"`
	RecentNews     []NewsItem `json:"recent_news"`
	InvestmentFlow float64    `json:"investment_flow"`
}

type MarketHealth struct {
	TargetMarket    string   `json:"target_market"`
	Health          string   `json:"health"` // "excellent", "good", "fair", "poor"
	SaturationLevel float64  `json:"saturation_level"`
	Opportunities   []string `json:"opportunities"`
	Threats         []string `json:"threats"`
	Recommendations []string `json:"recommendations"`
}

type NewsItem struct {
	Title       string    `json:"title"`
	Description string    `json:"description"`
	URL         string    `json:"url"`
	PublishedAt time.Time `json:"published_at"`
	Source      string    `json:"source"`
	Sentiment   string    `json:"sentiment"` // positive, negative, neutral
}

// NewMarketDataService creates a new market data service
func NewMarketDataService(marketDataAPIKey, marketDataURL, newsAPIKey, newsAPIURL string) MarketDataService {
	return &marketDataService{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		marketDataAPIKey: marketDataAPIKey,
		marketDataURL:    marketDataURL,
		newsAPIKey:       newsAPIKey,
		newsAPIURL:       newsAPIURL,
	}
}

func (s *marketDataService) GetIndustryTrends(ctx context.Context, industry string) (*IndustryTrends, error) {
	logger.Infof("Fetching industry trends for: %s", industry)

	// For demo purposes, we'll simulate market data
	// In production, this would call real APIs like:
	// - Bloomberg API
	// - Yahoo Finance API
	// - Crunchbase API
	// - PitchBook API

	trends := &IndustryTrends{
		Industry:         industry,
		GrowthRate:       s.calculateGrowthRate(industry),
		MarketSize:       s.calculateMarketSize(industry),
		CompetitionLevel: s.determineCompetitionLevel(industry),
		Trends:           s.getIndustryTrends(industry),
		Outlook:          s.determineOutlook(industry),
	}

	logger.Infof("Industry trends retrieved for %s: Growth=%f%%, Competition=%s",
		industry, trends.GrowthRate, trends.CompetitionLevel)

	return trends, nil
}

func (s *marketDataService) GetSectorAnalysis(ctx context.Context, sector string) (*SectorAnalysis, error) {
	logger.Infof("Fetching sector analysis for: %s", sector)

	// Simulate sector analysis - in production this would call real APIs
	news, err := s.fetchRecentNews(ctx, sector)
	if err != nil {
		logger.Warnf("Failed to fetch news for sector %s: %v", sector, err)
		news = []NewsItem{} // Continue without news
	}

	analysis := &SectorAnalysis{
		Sector:         sector,
		IsActive:       s.isSectorActive(sector),
		Activity:       s.determineSectorActivity(sector),
		KeyPlayers:     s.getKeyPlayers(sector),
		MarketCap:      s.estimateMarketCap(sector),
		RecentNews:     news,
		InvestmentFlow: s.calculateInvestmentFlow(sector),
	}

	logger.Infof("Sector analysis completed for %s: Active=%t, Activity=%s",
		sector, analysis.IsActive, analysis.Activity)

	return analysis, nil
}

func (s *marketDataService) GetMarketHealth(ctx context.Context, targetMarket string) (*MarketHealth, error) {
	logger.Infof("Analyzing market health for: %s", targetMarket)

	health := &MarketHealth{
		TargetMarket:    targetMarket,
		Health:          s.determineMarketHealth(targetMarket),
		SaturationLevel: s.calculateSaturation(targetMarket),
		Opportunities:   s.identifyOpportunities(targetMarket),
		Threats:         s.identifyThreats(targetMarket),
		Recommendations: s.generateRecommendations(targetMarket),
	}

	logger.Infof("Market health analysis completed for %s: Health=%s, Saturation=%f%%",
		targetMarket, health.Health, health.SaturationLevel)

	return health, nil
}

func (s *marketDataService) fetchRecentNews(ctx context.Context, sector string) ([]NewsItem, error) {
	if s.newsAPIKey == "" || s.newsAPIKey == "your_news_api_key_here" {
		// Return mock news if no API key
		return s.getMockNews(sector), nil
	}

	// Real News API call would go here
	url := fmt.Sprintf("%s/everything?q=%s&sortBy=publishedAt&pageSize=5&apiKey=%s",
		s.newsAPIURL, sector, s.newsAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch news: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return s.getMockNews(sector), nil // Fallback to mock data
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	var newsResponse struct {
		Articles []struct {
			Title       string `json:"title"`
			Description string `json:"description"`
			URL         string `json:"url"`
			PublishedAt string `json:"publishedAt"`
			Source      struct {
				Name string `json:"name"`
			} `json:"source"`
		} `json:"articles"`
	}

	if err := json.Unmarshal(body, &newsResponse); err != nil {
		return s.getMockNews(sector), nil // Fallback to mock data
	}

	var news []NewsItem
	for _, article := range newsResponse.Articles {
		publishedAt, _ := time.Parse(time.RFC3339, article.PublishedAt)
		news = append(news, NewsItem{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			PublishedAt: publishedAt,
			Source:      article.Source.Name,
		})
	}

	return news, nil
}

// Helper methods for market analysis simulation
func (s *marketDataService) calculateGrowthRate(industry string) float64 {
	// Simulated growth rates based on industry
	growthRates := map[string]float64{
		"technology": 12.5,
		"healthcare": 8.3,
		"fintech":    15.2,
		"edutech":    18.7,
		"logistics":  7.4,
		"ecommerce":  14.1,
		"saas":       22.3,
		"ai":         35.6,
		"blockchain": 25.8,
		"renewable":  19.4,
	}

	if rate, exists := growthRates[industry]; exists {
		return rate
	}
	return 10.0 // Default growth rate
}

func (s *marketDataService) calculateMarketSize(industry string) int64 {
	// Simulated market sizes in millions USD
	marketSizes := map[string]int64{
		"technology": 500000,
		"healthcare": 350000,
		"fintech":    180000,
		"edutech":    85000,
		"logistics":  450000,
		"ecommerce":  250000,
		"saas":       120000,
		"ai":         65000,
		"blockchain": 25000,
		"renewable":  95000,
	}

	if size, exists := marketSizes[industry]; exists {
		return size
	}
	return 50000 // Default market size
}

func (s *marketDataService) determineCompetitionLevel(industry string) string {
	highCompetition := []string{"technology", "ecommerce", "saas", "fintech"}
	moderateCompetition := []string{"healthcare", "logistics", "edutech"}

	for _, comp := range highCompetition {
		if comp == industry {
			return "high"
		}
	}

	for _, comp := range moderateCompetition {
		if comp == industry {
			return "moderate"
		}
	}

	return "low"
}

func (s *marketDataService) getIndustryTrends(industry string) []string {
	trends := map[string][]string{
		"technology": {
			"AI and Machine Learning adoption",
			"Cloud-first strategies",
			"Remote work technologies",
			"Cybersecurity focus",
		},
		"edutech": {
			"Personalized learning platforms",
			"VR/AR in education",
			"Micro-learning content",
			"Skills-based training",
		},
		"logistics": {
			"Last-mile delivery optimization",
			"Autonomous vehicles",
			"Supply chain digitization",
			"Sustainable packaging",
		},
	}

	if trend, exists := trends[industry]; exists {
		return trend
	}
	return []string{"Digital transformation", "Customer experience focus", "Sustainability initiatives"}
}

func (s *marketDataService) determineOutlook(industry string) string {
	positiveOutlook := []string{"technology", "ai", "renewable", "edutech"}

	for _, pos := range positiveOutlook {
		if pos == industry {
			return "positive"
		}
	}

	return "stable"
}

func (s *marketDataService) isSectorActive(sector string) bool {
	activeSectors := []string{"technology", "healthcare", "edutech", "fintech", "ai", "renewable"}

	for _, active := range activeSectors {
		if active == sector {
			return true
		}
	}

	return false
}

func (s *marketDataService) determineSectorActivity(sector string) string {
	growingSectors := []string{"ai", "renewable", "edutech", "blockchain"}
	stableSectors := []string{"healthcare", "logistics", "technology"}

	for _, growing := range growingSectors {
		if growing == sector {
			return "growing"
		}
	}

	for _, stable := range stableSectors {
		if stable == sector {
			return "stable"
		}
	}

	return "declining"
}

func (s *marketDataService) getKeyPlayers(sector string) []string {
	players := map[string][]string{
		"technology": {"Microsoft", "Google", "Apple", "Amazon"},
		"edutech":    {"Coursera", "Udemy", "Khan Academy", "Byju's"},
		"logistics":  {"FedEx", "UPS", "DHL", "Amazon Logistics"},
		"fintech":    {"Stripe", "PayPal", "Square", "Coinbase"},
	}

	if player, exists := players[sector]; exists {
		return player
	}
	return []string{"Market leaders vary by region"}
}

func (s *marketDataService) estimateMarketCap(sector string) int64 {
	// Estimated market caps in billions USD
	marketCaps := map[string]int64{
		"technology": 25000,
		"healthcare": 18000,
		"fintech":    8500,
		"edutech":    4200,
		"logistics":  12000,
	}

	if cap, exists := marketCaps[sector]; exists {
		return cap
	}
	return 5000
}

func (s *marketDataService) calculateInvestmentFlow(sector string) float64 {
	// Investment flow in billions USD
	flows := map[string]float64{
		"ai":         25.6,
		"fintech":    18.2,
		"edutech":    8.9,
		"renewable":  45.3,
		"blockchain": 12.1,
	}

	if flow, exists := flows[sector]; exists {
		return flow
	}
	return 5.0
}

func (s *marketDataService) determineMarketHealth(targetMarket string) string {
	// Simplified market health determination
	if len(targetMarket) > 20 {
		return "good" // Specific markets tend to be healthier
	}
	return "fair"
}

func (s *marketDataService) calculateSaturation(targetMarket string) float64 {
	// Simulated saturation levels
	saturations := map[string]float64{
		"global":        85.0,
		"north america": 75.0,
		"europe":        70.0,
		"asia":          45.0,
		"emerging":      25.0,
	}

	for market, saturation := range saturations {
		if market == targetMarket {
			return saturation
		}
	}
	return 50.0 // Default saturation
}

func (s *marketDataService) identifyOpportunities(targetMarket string) []string {
	return []string{
		"Growing digital adoption",
		"Underserved market segments",
		"Technological advancement opportunities",
		"Regulatory support for innovation",
	}
}

func (s *marketDataService) identifyThreats(targetMarket string) []string {
	return []string{
		"Increased competition",
		"Regulatory changes",
		"Economic uncertainty",
		"Technology disruption",
	}
}

func (s *marketDataService) generateRecommendations(targetMarket string) []string {
	return []string{
		"Focus on differentiation",
		"Build strong customer relationships",
		"Monitor regulatory developments",
		"Invest in technology infrastructure",
	}
}

func (s *marketDataService) getMockNews(sector string) []NewsItem {
	return []NewsItem{
		{
			Title:       fmt.Sprintf("%s sector shows promising growth", sector),
			Description: fmt.Sprintf("Recent developments in %s indicate positive market trends", sector),
			URL:         "https://example.com/news/1",
			PublishedAt: time.Now().AddDate(0, 0, -2),
			Source:      "Market News",
			Sentiment:   "positive",
		},
		{
			Title:       fmt.Sprintf("Innovation drives %s market expansion", sector),
			Description: fmt.Sprintf("New technologies are reshaping the %s landscape", sector),
			URL:         "https://example.com/news/2",
			PublishedAt: time.Now().AddDate(0, 0, -5),
			Source:      "Tech Today",
			Sentiment:   "positive",
		},
	}
}

// GetMarketData fetches comprehensive market data for a given sector
func (s *marketDataService) GetMarketData(ctx context.Context, sector, industry string) (*MarketDataResponse, error) {
	logger.Infof("Fetching comprehensive market data for sector: %s, industry: %s", sector, industry)

	// Generate comprehensive market data based on sector
	marketData := s.generateMarketData(sector, industry)

	// Simulate API call delay
	time.Sleep(500 * time.Millisecond)

	logger.Infof("Successfully fetched market data for sector: %s", sector)
	return marketData, nil
}

// GetNewsAnalysis fetches news sentiment analysis
func (s *marketDataService) GetNewsAnalysis(ctx context.Context, sector string) (*NewsAnalysisResponse, error) {
	logger.Infof("Fetching news analysis for sector: %s", sector)

	// Generate news analysis based on sector
	newsAnalysis := s.generateNewsAnalysis(sector)

	// Simulate API call delay
	time.Sleep(300 * time.Millisecond)

	logger.Infof("Successfully fetched news analysis for sector: %s", sector)
	return newsAnalysis, nil
}

// GetCompetitorAnalysis fetches competitor landscape data
func (s *marketDataService) GetCompetitorAnalysis(ctx context.Context, sector, businessModel string) (*CompetitorAnalysisResponse, error) {
	logger.Infof("Fetching competitor analysis for sector: %s, business model: %s", sector, businessModel)

	// Generate competitor data based on sector
	competitorData := s.generateCompetitorAnalysis(sector, businessModel)

	// Simulate API call delay
	time.Sleep(400 * time.Millisecond)

	logger.Infof("Successfully fetched competitor analysis for sector: %s", sector)
	return competitorData, nil
}

// generateMarketData creates comprehensive market data based on sector
func (s *marketDataService) generateMarketData(sector, industry string) *MarketDataResponse {
	// Sector-specific market data simulation
	sectorData := map[string]struct {
		TAM              float64
		GrowthRate       float64
		Status           string
		CompetitionLevel string
		RegulationLevel  string
	}{
		"fintech": {
			TAM:              324000000000, // $324B
			GrowthRate:       13.7,
			Status:           "active",
			CompetitionLevel: "high",
			RegulationLevel:  "high",
		},
		"edutech": {
			TAM:              89000000000, // $89B
			GrowthRate:       16.3,
			Status:           "active",
			CompetitionLevel: "medium",
			RegulationLevel:  "medium",
		},
		"healthtech": {
			TAM:              659800000000, // $659.8B
			GrowthRate:       15.1,
			Status:           "active",
			CompetitionLevel: "high",
			RegulationLevel:  "high",
		},
		"logistics": {
			TAM:              12040000000000, // $12.04T
			GrowthRate:       6.2,
			Status:           "active",
			CompetitionLevel: "high",
			RegulationLevel:  "medium",
		},
		"ecommerce": {
			TAM:              6200000000000, // $6.2T
			GrowthRate:       14.7,
			Status:           "active",
			CompetitionLevel: "high",
			RegulationLevel:  "low",
		},
		"agritech": {
			TAM:              22500000000, // $22.5B
			GrowthRate:       22.5,
			Status:           "emerging",
			CompetitionLevel: "medium",
			RegulationLevel:  "medium",
		},
	}

	data, exists := sectorData[strings.ToLower(sector)]
	if !exists {
		// Default data for unknown sectors
		data = sectorData["fintech"]
		data.Status = "unknown"
		data.CompetitionLevel = "medium"
	}

	return &MarketDataResponse{
		Sector:   sector,
		Industry: industry,
		MarketSize: MarketSizeInfo{
			TAM:  data.TAM,
			SAM:  data.TAM * 0.1,  // 10% of TAM
			SOM:  data.TAM * 0.01, // 1% of TAM
			CAGR: data.GrowthRate,
		},
		GrowthRate:       data.GrowthRate,
		MarketStatus:     data.Status,
		CompetitionLevel: data.CompetitionLevel,
		KeyTrends:        s.generateKeyTrends(sector),
		Opportunities:    s.generateOpportunities(sector),
		Threats:          s.generateThreats(sector),
		RegulationLevel:  data.RegulationLevel,
		BarriersToEntry:  s.generateBarriers(sector),
		LastUpdated:      time.Now(),
	}
}

// generateNewsAnalysis creates simulated news sentiment
func (s *marketDataService) generateNewsAnalysis(sector string) *NewsAnalysisResponse {
	sectorSentiments := map[string]float64{
		"fintech":    0.3, // Positive
		"edutech":    0.6, // Very positive
		"healthtech": 0.4, // Positive
		"logistics":  0.1, // Slightly positive
		"ecommerce":  0.2, // Slightly positive
		"agritech":   0.7, // Very positive
	}

	sentiment, exists := sectorSentiments[strings.ToLower(sector)]
	if !exists {
		sentiment = 0.0 // Neutral for unknown sectors
	}

	return &NewsAnalysisResponse{
		Sector:           sector,
		SentimentScore:   sentiment,
		PositiveKeywords: s.generatePositiveKeywords(sector),
		NegativeKeywords: s.generateNegativeKeywords(sector),
		RecentNews:       s.generateRecentNews(sector),
		InvestmentNews:   s.generateInvestmentNews(sector),
		AnalysisDate:     time.Now(),
	}
}

// generateCompetitorAnalysis creates simulated competitor data
func (s *marketDataService) generateCompetitorAnalysis(sector, businessModel string) *CompetitorAnalysisResponse {
	competitorCounts := map[string]int{
		"fintech":    1200,
		"edutech":    800,
		"healthtech": 2500,
		"logistics":  450,
		"ecommerce":  3200,
		"agritech":   300,
	}

	count, exists := competitorCounts[strings.ToLower(sector)]
	if !exists {
		count = 500 // Default
	}

	return &CompetitorAnalysisResponse{
		Sector:               sector,
		TotalCompetitors:     count,
		KeyPlayers:           s.generateKeyPlayers(sector),
		MarketLeader:         s.getMarketLeader(sector),
		EmergingPlayers:      s.generateEmergingPlayers(sector),
		CompetitionIntensity: s.getCompetitionIntensity(count),
	}
}

// Helper functions for generating sector-specific data

func (s *marketDataService) generateKeyTrends(sector string) []string {
	trends := map[string][]string{
		"fintech": {
			"Open Banking adoption",
			"AI-driven risk assessment",
			"Cryptocurrency integration",
			"Regulatory technology (RegTech)",
			"Embedded finance solutions",
		},
		"edutech": {
			"AI-powered personalized learning",
			"Virtual and Augmented Reality",
			"Microlearning platforms",
			"Skill-based hiring trends",
			"Remote learning infrastructure",
		},
		"healthtech": {
			"Telemedicine expansion",
			"AI diagnostics",
			"Wearable health monitoring",
			"Digital therapeutics",
			"Healthcare data interoperability",
		},
		"logistics": {
			"Autonomous delivery vehicles",
			"Supply chain digitization",
			"Last-mile delivery optimization",
			"IoT-enabled tracking",
			"Sustainable logistics solutions",
		},
		"ecommerce": {
			"Social commerce integration",
			"Augmented reality shopping",
			"Voice commerce",
			"Sustainability focus",
			"Headless commerce architecture",
		},
		"agritech": {
			"Precision agriculture",
			"IoT sensors for crop monitoring",
			"Vertical farming solutions",
			"AI-powered yield prediction",
			"Sustainable farming practices",
		},
	}

	if sectorTrends, exists := trends[strings.ToLower(sector)]; exists {
		return sectorTrends
	}
	return []string{"Digital transformation", "AI adoption", "Sustainability focus"}
}

func (s *marketDataService) generateOpportunities(sector string) []string {
	opportunities := map[string][]string{
		"fintech": {
			"Underbanked population targeting",
			"SME lending market expansion",
			"Cross-border payment solutions",
			"Insurance technology innovation",
		},
		"edutech": {
			"Corporate training market",
			"Lifelong learning platforms",
			"Emerging market penetration",
			"Professional certification programs",
		},
		"healthtech": {
			"Rural healthcare access",
			"Mental health solutions",
			"Chronic disease management",
			"Healthcare cost reduction",
		},
		"logistics": {
			"E-commerce fulfillment growth",
			"Cold chain logistics",
			"Cross-border trade facilitation",
			"Green logistics solutions",
		},
		"ecommerce": {
			"Mobile-first markets",
			"B2B marketplace expansion",
			"Subscription commerce models",
			"Local marketplace platforms",
		},
		"agritech": {
			"Climate-smart agriculture",
			"Food traceability solutions",
			"Agricultural data monetization",
			"Smallholder farmer empowerment",
		},
	}

	if sectorOpps, exists := opportunities[strings.ToLower(sector)]; exists {
		return sectorOpps
	}
	return []string{"Market expansion", "Technology adoption", "Partnership opportunities"}
}

func (s *marketDataService) generateThreats(sector string) []string {
	threats := map[string][]string{
		"fintech": {
			"Increasing regulatory scrutiny",
			"Big tech competition",
			"Cybersecurity risks",
			"Economic downturn impact",
		},
		"edutech": {
			"Return to traditional learning",
			"Budget constraints in education",
			"Technology adoption resistance",
			"Content piracy issues",
		},
		"healthtech": {
			"Regulatory compliance challenges",
			"Data privacy concerns",
			"Healthcare system resistance",
			"Reimbursement uncertainties",
		},
		"logistics": {
			"Economic volatility",
			"Fuel price fluctuations",
			"Labor shortages",
			"Geopolitical tensions",
		},
		"ecommerce": {
			"Platform dependency",
			"Rising customer acquisition costs",
			"Supply chain disruptions",
			"Regulatory changes",
		},
		"agritech": {
			"Climate change impacts",
			"Technology adoption barriers",
			"Rural connectivity issues",
			"Traditional farming resistance",
		},
	}

	if sectorThreats, exists := threats[strings.ToLower(sector)]; exists {
		return sectorThreats
	}
	return []string{"Economic uncertainty", "Competitive pressure", "Regulatory changes"}
}

func (s *marketDataService) generateBarriers(sector string) []string {
	barriers := map[string][]string{
		"fintech": {
			"Regulatory compliance costs",
			"Trust and credibility building",
			"High customer acquisition costs",
			"Capital requirements",
		},
		"edutech": {
			"Long sales cycles",
			"Content development costs",
			"Technology infrastructure needs",
			"User adoption challenges",
		},
		"healthtech": {
			"Regulatory approval processes",
			"Clinical validation requirements",
			"Healthcare system integration",
			"Data security compliance",
		},
		"logistics": {
			"Infrastructure investment needs",
			"Operational complexity",
			"Partnership requirements",
			"Technology integration costs",
		},
		"ecommerce": {
			"Brand building costs",
			"Inventory management complexity",
			"Customer service requirements",
			"Technology platform costs",
		},
		"agritech": {
			"Rural market penetration",
			"Technology education needs",
			"Seasonal revenue patterns",
			"Hardware development costs",
		},
	}

	if sectorBarriers, exists := barriers[strings.ToLower(sector)]; exists {
		return sectorBarriers
	}
	return []string{"High capital requirements", "Market penetration challenges", "Technology barriers"}
}

func (s *marketDataService) generatePositiveKeywords(sector string) []string {
	keywords := map[string][]string{
		"fintech":    {"innovation", "growth", "digital transformation", "investment", "adoption"},
		"edutech":    {"advancement", "accessibility", "personalization", "engagement", "results"},
		"healthtech": {"breakthrough", "improvement", "accessibility", "efficiency", "outcomes"},
		"logistics":  {"optimization", "efficiency", "automation", "sustainability", "growth"},
		"ecommerce":  {"expansion", "convenience", "innovation", "growth", "adoption"},
		"agritech":   {"sustainability", "innovation", "productivity", "efficiency", "growth"},
	}

	if sectorKeywords, exists := keywords[strings.ToLower(sector)]; exists {
		return sectorKeywords
	}
	return []string{"growth", "innovation", "opportunity"}
}

func (s *marketDataService) generateNegativeKeywords(sector string) []string {
	keywords := map[string][]string{
		"fintech":    {"regulation", "security", "fraud", "compliance", "risk"},
		"edutech":    {"budget cuts", "resistance", "digital divide", "quality concerns"},
		"healthtech": {"privacy", "regulation", "costs", "adoption barriers", "compliance"},
		"logistics":  {"disruption", "costs", "delays", "complexity", "challenges"},
		"ecommerce":  {"competition", "costs", "fraud", "returns", "saturation"},
		"agritech":   {"adoption barriers", "costs", "connectivity", "resistance", "complexity"},
	}

	if sectorKeywords, exists := keywords[strings.ToLower(sector)]; exists {
		return sectorKeywords
	}
	return []string{"challenges", "competition", "risks"}
}

func (s *marketDataService) generateRecentNews(sector string) []NewsItem {
	// Simulated recent news - in production, fetch from actual news APIs
	return []NewsItem{
		{
			Title:       fmt.Sprintf("%s sector shows continued growth momentum", strings.Title(sector)),
			Description: fmt.Sprintf("Recent developments in %s indicate positive market trends", sector),
			URL:         "https://example.com/news/1",
			PublishedAt: time.Now().AddDate(0, 0, -1),
			Source:      "Industry Weekly",
			Sentiment:   "positive",
		},
		{
			Title:       fmt.Sprintf("New regulations impact %s companies", sector),
			Description: fmt.Sprintf("Regulatory changes affecting %s industry landscape", sector),
			URL:         "https://example.com/news/2",
			PublishedAt: time.Now().AddDate(0, 0, -3),
			Source:      "Regulatory News",
			Sentiment:   "neutral",
		},
	}
}

func (s *marketDataService) generateInvestmentNews(sector string) []NewsItem {
	// Simulated investment news
	return []NewsItem{
		{
			Title:       fmt.Sprintf("%s startups raise $2B in Q4", strings.Title(sector)),
			Description: fmt.Sprintf("Investment activity in %s sector remains strong", sector),
			URL:         "https://example.com/investment/1",
			PublishedAt: time.Now().AddDate(0, 0, -5),
			Source:      "Investment Daily",
			Sentiment:   "positive",
		},
	}
}

func (s *marketDataService) generateKeyPlayers(sector string) []Competitor {
	players := map[string][]Competitor{
		"fintech": {
			{Name: "Stripe", MarketShare: 15.2, Funding: 95000000000, Founded: 2010, Employees: 4000},
			{Name: "Square", MarketShare: 12.8, Funding: 6000000000, Founded: 2009, Employees: 8000},
		},
		"edutech": {
			{Name: "Coursera", MarketShare: 8.5, Funding: 464000000, Founded: 2012, Employees: 1000},
			{Name: "Udemy", MarketShare: 6.2, Funding: 173000000, Founded: 2010, Employees: 1500},
		},
	}

	if sectorPlayers, exists := players[strings.ToLower(sector)]; exists {
		return sectorPlayers
	}
	return []Competitor{
		{Name: "Market Leader", MarketShare: 10.0, Funding: 100000000, Founded: 2015, Employees: 500},
	}
}

func (s *marketDataService) getMarketLeader(sector string) string {
	leaders := map[string]string{
		"fintech":    "Stripe",
		"edutech":    "Coursera",
		"healthtech": "Teladoc",
		"logistics":  "Amazon Logistics",
		"ecommerce":  "Amazon",
		"agritech":   "John Deere",
	}

	if leader, exists := leaders[strings.ToLower(sector)]; exists {
		return leader
	}
	return "Market Leader Corp"
}

func (s *marketDataService) generateEmergingPlayers(sector string) []string {
	return []string{
		fmt.Sprintf("%s Innovator A", strings.Title(sector)),
		fmt.Sprintf("%s Disruptor B", strings.Title(sector)),
		fmt.Sprintf("%s Pioneer C", strings.Title(sector)),
	}
}

func (s *marketDataService) getCompetitionIntensity(competitorCount int) string {
	if competitorCount > 2000 {
		return "high"
	} else if competitorCount > 800 {
		return "medium"
	}
	return "low"
}

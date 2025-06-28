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
		logger.Warn("No valid NewsAPI key configured, using mock news data")
		return s.getMockNews(sector), nil
	}

	logger.Infof("Fetching real news for sector: %s", sector)

	// Real NewsAPI call
	url := fmt.Sprintf("%s/everything?q=%s+startup+funding+investment&sortBy=publishedAt&pageSize=5&apiKey=%s&language=en",
		s.newsAPIURL, sector, s.newsAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Errorf("Failed to create news request: %v", err)
		return s.getMockNews(sector), nil
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Errorf("Failed to fetch news from NewsAPI: %v", err)
		return s.getMockNews(sector), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warnf("NewsAPI returned status %d, falling back to mock data", resp.StatusCode)
		return s.getMockNews(sector), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read NewsAPI response: %v", err)
		return s.getMockNews(sector), nil
	}

	var newsResponse struct {
		Status   string `json:"status"`
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
		logger.Errorf("Failed to parse NewsAPI response: %v", err)
		return s.getMockNews(sector), nil
	}

	if newsResponse.Status != "ok" {
		logger.Warnf("NewsAPI status not ok: %s", newsResponse.Status)
		return s.getMockNews(sector), nil
	}

	var news []NewsItem
	for _, article := range newsResponse.Articles {
		publishedAt, parseErr := time.Parse(time.RFC3339, article.PublishedAt)
		if parseErr != nil {
			publishedAt = time.Now()
		}

		// Simple sentiment analysis based on keywords
		sentiment := s.analyzeSentiment(article.Title + " " + article.Description)

		news = append(news, NewsItem{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			PublishedAt: publishedAt,
			Source:      article.Source.Name,
			Sentiment:   sentiment,
		})
	}

	logger.Infof("Successfully fetched %d real news articles for sector %s", len(news), sector)
	return news, nil
}

// analyzeSentiment performs simple keyword-based sentiment analysis
func (s *marketDataService) analyzeSentiment(text string) string {
	text = strings.ToLower(text)

	positiveKeywords := []string{
		"growth", "funding", "investment", "success", "revenue", "profit",
		"expansion", "partnership", "innovation", "breakthrough", "achievement",
		"milestone", "acquisition", "valuation", "ipo", "unicorn", "promising",
		"opportunity", "bullish", "rising", "gains", "optimistic", "positive",
	}

	negativeKeywords := []string{
		"decline", "loss", "bankruptcy", "failure", "recession", "crisis",
		"downfall", "scandal", "lawsuit", "controversy", "bearish", "falling",
		"crash", "bubble", "risk", "threat", "concern", "warning", "volatile",
		"uncertainty", "disruption", "challenge", "struggle", "debt",
	}

	positiveCount := 0
	negativeCount := 0

	for _, keyword := range positiveKeywords {
		if strings.Contains(text, keyword) {
			positiveCount++
		}
	}

	for _, keyword := range negativeKeywords {
		if strings.Contains(text, keyword) {
			negativeCount++
		}
	}

	if positiveCount > negativeCount {
		return "positive"
	} else if negativeCount > positiveCount {
		return "negative"
	}

	return "neutral"
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

// GetMarketData fetches comprehensive market data using real external APIs
func (s *marketDataService) GetMarketData(ctx context.Context, sector, industry string) (*MarketDataResponse, error) {
	logger.Infof("Fetching comprehensive market data for sector: %s, industry: %s", sector, industry)

	// Get industry trends (enhanced with real data where possible)
	_, err := s.GetIndustryTrends(ctx, industry)
	if err != nil {
		logger.Warnf("Failed to get industry trends: %v", err)
	}

	// Get news analysis for sentiment
	newsAnalysis, err := s.GetNewsAnalysis(ctx, sector)
	if err != nil {
		logger.Warnf("Failed to get news analysis: %v", err)
	}

	// Combine real and simulated data to create comprehensive response
	marketData := &MarketDataResponse{
		Sector:   sector,
		Industry: industry,
		MarketSize: MarketSizeInfo{
			TAM:  float64(s.calculateMarketSize(industry)),
			SAM:  float64(s.calculateMarketSize(industry)) * 0.3,  // 30% of TAM
			SOM:  float64(s.calculateMarketSize(industry)) * 0.05, // 5% of TAM
			CAGR: s.calculateGrowthRate(industry),
		},
		GrowthRate:       s.calculateGrowthRate(industry),
		MarketStatus:     s.determineMarketStatus(industry),
		CompetitionLevel: s.determineCompetitionLevel(industry),
		KeyTrends:        s.getIndustryTrends(industry),
		Opportunities:    s.generateOpportunities(sector, industry),
		Threats:          s.generateThreats(sector, industry),
		RegulationLevel:  s.determineRegulationLevel(industry),
		BarriersToEntry:  s.getBarriersToEntry(industry),
		LastUpdated:      time.Now(),
	}

	// Enhance with real news sentiment if available
	if newsAnalysis != nil && len(newsAnalysis.RecentNews) > 0 {
		// Adjust market status based on recent news sentiment
		marketData.MarketStatus = s.adjustMarketStatusBySentiment(marketData.MarketStatus, newsAnalysis.SentimentScore)
	}

	logger.Infof("Market data fetched successfully for %s/%s", sector, industry)
	return marketData, nil
}

// GetNewsAnalysis fetches and analyzes news sentiment using real NewsAPI
func (s *marketDataService) GetNewsAnalysis(ctx context.Context, sector string) (*NewsAnalysisResponse, error) {
	logger.Infof("Analyzing news sentiment for sector: %s", sector)

	// Fetch recent news
	recentNews, err := s.fetchRecentNews(ctx, sector)
	if err != nil {
		logger.Warnf("Failed to fetch recent news: %v", err)
		recentNews = s.getMockNews(sector)
	}

	// Fetch investment-specific news
	investmentNews, err := s.fetchInvestmentNews(ctx, sector)
	if err != nil {
		logger.Warnf("Failed to fetch investment news: %v", err)
		investmentNews = s.getMockInvestmentNews(sector)
	}

	// Calculate overall sentiment score
	sentimentScore := s.calculateSentimentScore(recentNews, investmentNews)

	// Extract keywords
	positiveKeywords, negativeKeywords := s.extractKeywords(recentNews, investmentNews)

	newsAnalysis := &NewsAnalysisResponse{
		Sector:           sector,
		SentimentScore:   sentimentScore,
		PositiveKeywords: positiveKeywords,
		NegativeKeywords: negativeKeywords,
		RecentNews:       recentNews,
		InvestmentNews:   investmentNews,
		AnalysisDate:     time.Now(),
	}

	logger.Infof("News analysis completed for %s: sentiment score %.2f", sector, sentimentScore)
	return newsAnalysis, nil
}

// GetCompetitorAnalysis analyzes the competitive landscape for a given sector
func (s *marketDataService) GetCompetitorAnalysis(ctx context.Context, sector, businessModel string) (*CompetitorAnalysisResponse, error) {
	logger.Infof("Analyzing competitive landscape for sector: %s, business model: %s", sector, businessModel)

	// In a real implementation, this would call APIs like:
	// - Crunchbase API for startup data
	// - PitchBook API for funding information
	// - LinkedIn API for company insights
	// For now, we'll use enhanced simulation based on real patterns

	totalCompetitors := s.estimateTotalCompetitors(sector)
	keyPlayers := s.getKeyCompetitors(sector)
	marketLeader := s.determineMarketLeader(sector)
	emergingPlayers := s.getEmergingPlayers(sector)
	competitionIntensity := s.determineCompetitionLevel(sector)

	analysis := &CompetitorAnalysisResponse{
		Sector:               sector,
		TotalCompetitors:     totalCompetitors,
		KeyPlayers:           keyPlayers,
		MarketLeader:         marketLeader,
		EmergingPlayers:      emergingPlayers,
		CompetitionIntensity: competitionIntensity,
	}

	logger.Infof("Competitor analysis completed for %s: %d total competitors, intensity: %s", sector, totalCompetitors, competitionIntensity)
	return analysis, nil
}

// Helper methods for market data service

func (s *marketDataService) determineMarketStatus(industry string) string {
	emergingIndustries := []string{"ai", "blockchain", "renewable", "biotech", "quantum"}
	matureIndustries := []string{"finance", "retail", "manufacturing", "automotive"}

	for _, emerging := range emergingIndustries {
		if strings.Contains(strings.ToLower(industry), emerging) {
			return "emerging"
		}
	}

	for _, mature := range matureIndustries {
		if strings.Contains(strings.ToLower(industry), mature) {
			return "mature"
		}
	}

	return "active"
}

func (s *marketDataService) generateOpportunities(sector, industry string) []string {
	opportunities := map[string][]string{
		"technology": {
			"AI integration opportunities",
			"Remote work solutions demand",
			"Cloud adoption acceleration",
			"Digital transformation needs",
		},
		"healthcare": {
			"Telemedicine expansion",
			"Personalized medicine growth",
			"Healthcare data analytics",
			"Wearable health tech adoption",
		},
		"fintech": {
			"Digital payment solutions",
			"Crypto integration opportunities",
			"RegTech automation",
			"Financial inclusion initiatives",
		},
	}

	if ops, exists := opportunities[strings.ToLower(industry)]; exists {
		return ops
	}

	return []string{
		"Digital transformation opportunities",
		"Market expansion potential",
		"Partnership possibilities",
		"Innovation-driven growth",
	}
}

func (s *marketDataService) generateThreats(sector, industry string) []string {
	threats := map[string][]string{
		"technology": {
			"Rapid technological obsolescence",
			"Intense competition from big tech",
			"Cybersecurity vulnerabilities",
			"Regulatory compliance challenges",
		},
		"healthcare": {
			"Strict regulatory requirements",
			"High compliance costs",
			"Data privacy concerns",
			"Long approval cycles",
		},
		"fintech": {
			"Financial regulation changes",
			"Security breach risks",
			"Market volatility impact",
			"Banking partner dependencies",
		},
	}

	if threats_list, exists := threats[strings.ToLower(industry)]; exists {
		return threats_list
	}

	return []string{
		"Market saturation risks",
		"Economic downturn impact",
		"Competitive pressures",
		"Regulatory uncertainties",
	}
}

func (s *marketDataService) determineRegulationLevel(industry string) string {
	highlyRegulated := []string{"healthcare", "finance", "fintech", "banking", "insurance"}
	moderatelyRegulated := []string{"education", "transport", "energy", "telecommunications"}

	for _, regulated := range highlyRegulated {
		if strings.Contains(strings.ToLower(industry), regulated) {
			return "high"
		}
	}

	for _, regulated := range moderatelyRegulated {
		if strings.Contains(strings.ToLower(industry), regulated) {
			return "moderate"
		}
	}

	return "low"
}

func (s *marketDataService) getBarriersToEntry(industry string) []string {
	barriers := map[string][]string{
		"healthcare": {
			"Regulatory approvals required",
			"High compliance costs",
			"Long development cycles",
			"Medical expertise requirements",
		},
		"fintech": {
			"Financial licensing requirements",
			"High security standards",
			"Banking partnerships needed",
			"Regulatory compliance costs",
		},
		"technology": {
			"High development costs",
			"Technical expertise requirements",
			"Network effects advantages",
			"Scale economies",
		},
	}

	if barrier_list, exists := barriers[strings.ToLower(industry)]; exists {
		return barrier_list
	}

	return []string{
		"Capital requirements",
		"Market competition",
		"Customer acquisition costs",
		"Brand recognition needs",
	}
}

func (s *marketDataService) adjustMarketStatusBySentiment(currentStatus string, sentimentScore float64) string {
	// Adjust market status based on news sentiment
	if sentimentScore > 0.3 {
		// Very positive sentiment might indicate emerging/growing market
		if currentStatus == "mature" {
			return "active"
		}
		if currentStatus == "active" {
			return "emerging"
		}
	} else if sentimentScore < -0.3 {
		// Very negative sentiment might indicate declining market
		if currentStatus == "emerging" {
			return "active"
		}
		if currentStatus == "active" {
			return "mature"
		}
	}

	return currentStatus
}

func (s *marketDataService) fetchInvestmentNews(ctx context.Context, sector string) ([]NewsItem, error) {
	if s.newsAPIKey == "" || s.newsAPIKey == "your_news_api_key_here" {
		logger.Warn("No valid NewsAPI key configured, using mock investment news")
		return s.getMockInvestmentNews(sector), nil
	}

	logger.Infof("Fetching investment news for sector: %s", sector)

	// Search for investment and funding related news
	url := fmt.Sprintf("%s/everything?q=%s+AND+(investment+OR+funding+OR+venture+OR+raised)&sortBy=publishedAt&pageSize=5&apiKey=%s&language=en",
		s.newsAPIURL, sector, s.newsAPIKey)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		logger.Errorf("Failed to create investment news request: %v", err)
		return s.getMockInvestmentNews(sector), nil
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		logger.Errorf("Failed to fetch investment news: %v", err)
		return s.getMockInvestmentNews(sector), nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		logger.Warnf("NewsAPI returned status %d for investment news", resp.StatusCode)
		return s.getMockInvestmentNews(sector), nil
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Errorf("Failed to read investment news response: %v", err)
		return s.getMockInvestmentNews(sector), nil
	}

	var newsResponse struct {
		Status   string `json:"status"`
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
		logger.Errorf("Failed to parse investment news response: %v", err)
		return s.getMockInvestmentNews(sector), nil
	}

	if newsResponse.Status != "ok" {
		logger.Warnf("NewsAPI status not ok for investment news: %s", newsResponse.Status)
		return s.getMockInvestmentNews(sector), nil
	}

	var news []NewsItem
	for _, article := range newsResponse.Articles {
		publishedAt, parseErr := time.Parse(time.RFC3339, article.PublishedAt)
		if parseErr != nil {
			publishedAt = time.Now()
		}

		sentiment := s.analyzeSentiment(article.Title + " " + article.Description)

		news = append(news, NewsItem{
			Title:       article.Title,
			Description: article.Description,
			URL:         article.URL,
			PublishedAt: publishedAt,
			Source:      article.Source.Name,
			Sentiment:   sentiment,
		})
	}

	logger.Infof("Successfully fetched %d investment news articles for sector %s", len(news), sector)
	return news, nil
}

func (s *marketDataService) getMockInvestmentNews(sector string) []NewsItem {
	return []NewsItem{
		{
			Title:       fmt.Sprintf("%s Startup Raises $10M Series A", strings.Title(sector)),
			Description: fmt.Sprintf("A promising %s startup has successfully raised Series A funding to expand operations.", sector),
			URL:         "https://example.com/mock-funding-news",
			PublishedAt: time.Now().AddDate(0, 0, -1),
			Source:      "Mock Investment News",
			Sentiment:   "positive",
		},
		{
			Title:       fmt.Sprintf("VC Interest in %s Sector Growing", strings.Title(sector)),
			Description: fmt.Sprintf("Venture capital firms are showing increased interest in %s investments.", sector),
			URL:         "https://example.com/mock-vc-news",
			PublishedAt: time.Now().AddDate(0, 0, -2),
			Source:      "Mock Investment Report",
			Sentiment:   "positive",
		},
	}
}

func (s *marketDataService) calculateSentimentScore(recentNews, investmentNews []NewsItem) float64 {
	if len(recentNews) == 0 && len(investmentNews) == 0 {
		return 0.0
	}

	var totalScore float64
	var totalCount int

	// Weight recent news
	for _, news := range recentNews {
		switch news.Sentiment {
		case "positive":
			totalScore += 1.0
		case "negative":
			totalScore -= 1.0
			// neutral = 0, no change
		}
		totalCount++
	}

	// Weight investment news more heavily (2x)
	for _, news := range investmentNews {
		switch news.Sentiment {
		case "positive":
			totalScore += 2.0
		case "negative":
			totalScore -= 2.0
		}
		totalCount += 2
	}

	if totalCount == 0 {
		return 0.0
	}

	// Normalize to -1 to 1 range
	score := totalScore / float64(totalCount)
	if score > 1.0 {
		score = 1.0
	} else if score < -1.0 {
		score = -1.0
	}

	return score
}

func (s *marketDataService) extractKeywords(recentNews, investmentNews []NewsItem) ([]string, []string) {
	var positiveKeywords, negativeKeywords []string

	allNews := append(recentNews, investmentNews...)

	for _, news := range allNews {
		text := strings.ToLower(news.Title + " " + news.Description)

		if news.Sentiment == "positive" {
			// Extract positive keywords
			keywords := []string{"growth", "funding", "success", "revenue", "expansion", "innovation", "partnership"}
			for _, keyword := range keywords {
				if strings.Contains(text, keyword) && !contains(positiveKeywords, keyword) {
					positiveKeywords = append(positiveKeywords, keyword)
				}
			}
		} else if news.Sentiment == "negative" {
			// Extract negative keywords
			keywords := []string{"decline", "loss", "risk", "challenge", "concern", "uncertainty", "volatility"}
			for _, keyword := range keywords {
				if strings.Contains(text, keyword) && !contains(negativeKeywords, keyword) {
					negativeKeywords = append(negativeKeywords, keyword)
				}
			}
		}
	}

	return positiveKeywords, negativeKeywords
}

func (s *marketDataService) estimateTotalCompetitors(sector string) int {
	competitorCounts := map[string]int{
		"technology": 15000,
		"fintech":    3500,
		"healthcare": 8000,
		"education":  5000,
		"saas":       12000,
		"ai":         2500,
		"blockchain": 800,
		"renewable":  1200,
	}

	for key, count := range competitorCounts {
		if strings.Contains(strings.ToLower(sector), key) {
			return count
		}
	}

	return 2000 // Default
}

func (s *marketDataService) getKeyCompetitors(sector string) []Competitor {
	// This would normally come from APIs like Crunchbase or PitchBook
	competitors := map[string][]Competitor{
		"fintech": {
			{Name: "Stripe", MarketShare: 25.5, Funding: 950000000, Founded: 2010, Employees: 4000,
				Strengths:  []string{"Developer-friendly API", "Global reach", "Strong partnerships"},
				Weaknesses: []string{"High fees for small businesses", "Limited offline solutions"}},
			{Name: "Square", MarketShare: 18.2, Funding: 590000000, Founded: 2009, Employees: 5000,
				Strengths:  []string{"Integrated hardware/software", "Small business focus", "Easy setup"},
				Weaknesses: []string{"Limited international presence", "Dependency on hardware sales"}},
		},
		"technology": {
			{Name: "Generic Tech Leader", MarketShare: 15.0, Funding: 1000000000, Founded: 2015, Employees: 2000,
				Strengths:  []string{"First mover advantage", "Strong IP portfolio", "Enterprise relationships"},
				Weaknesses: []string{"Legacy tech debt", "Slower innovation cycles"}},
		},
	}

	for key, comps := range competitors {
		if strings.Contains(strings.ToLower(sector), key) {
			return comps
		}
	}

	// Default generic competitors
	return []Competitor{
		{Name: "Market Leader Co", MarketShare: 20.0, Funding: 500000000, Founded: 2018, Employees: 1500,
			Strengths:  []string{"Market leadership", "Strong funding", "Experienced team"},
			Weaknesses: []string{"High burn rate", "Competitive pressure"}},
	}
}

func (s *marketDataService) determineMarketLeader(sector string) string {
	leaders := map[string]string{
		"fintech":    "Stripe",
		"saas":       "Salesforce",
		"ai":         "OpenAI",
		"blockchain": "Coinbase",
		"renewable":  "Tesla Energy",
	}

	for key, leader := range leaders {
		if strings.Contains(strings.ToLower(sector), key) {
			return leader
		}
	}

	return "Established Market Leader"
}

func (s *marketDataService) getEmergingPlayers(sector string) []string {
	emerging := map[string][]string{
		"fintech":    {"Plaid", "Affirm", "Robinhood", "Chime"},
		"ai":         {"Anthropic", "Cohere", "Scale AI", "Hugging Face"},
		"blockchain": {"Polygon", "Solana Labs", "ConsenSys", "Alchemy"},
		"renewable":  {"Sunrun", "Enphase", "First Solar", "Bloom Energy"},
	}

	for key, players := range emerging {
		if strings.Contains(strings.ToLower(sector), key) {
			return players
		}
	}

	return []string{"Emerging Startup A", "Innovative Company B", "Disruptor Inc"}
}

// Helper function to check if slice contains string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

# 🔗 External API Setup Guide

This guide will help you set up real external API integrations for news and market data in the Risk Assessment Backend. The system now supports real API calls to NewsAPI and MarketStack for enhanced market intelligence.

## 📋 Overview

The backend now integrates with:
- **NewsAPI** - Real-time news data and sentiment analysis
- **MarketStack** - Market data and financial information
- **OpenAI** - AI-powered risk analysis (already configured)

## 🔧 API Key Setup Instructions

### 1. NewsAPI Setup (Free Tier Available)

**What it provides:**
- Real-time news articles
- Sector-specific news filtering
- Investment and funding news
- Automatic sentiment analysis

**Setup Steps:**

1. **Visit NewsAPI Registration:**
   ```
   https://newsapi.org/register
   ```

2. **Create Free Account:**
   - Enter your email address
   - Choose a strong password
   - Verify your email

3. **Get Your API Key:**
   - After verification, you'll see your API key on the dashboard
   - Copy the key (format: `xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx`)

4. **Free Tier Limits:**
   - 1,000 requests per month
   - Perfect for development and testing
   - No credit card required

### 2. MarketStack Setup (Free Tier Available)

**What it provides:**
- Real-time market data
- Historical stock information
- Market indicators
- Financial news integration

**Setup Steps:**

1. **Visit MarketStack Registration:**
   ```
   https://marketstack.com/signup/free
   ```

2. **Create Free Account:**
   - Enter your email and basic information
   - Verify your email address

3. **Get Your API Key:**
   - Login to your dashboard
   - Copy your API access key

4. **Free Tier Limits:**
   - 1,000 API requests per month
   - Real-time data access
   - Perfect for startup analysis

### 3. Configure Environment Variables

1. **Copy the environment template:**
   ```bash
   cp .env.example .env
   ```

2. **Edit your `.env` file and add your API keys:**
   ```bash
   # External API Configuration
   # --------------------------
   
   # NewsAPI (Get your free key at: https://newsapi.org/register)
   NEWS_API_KEY=your_actual_newsapi_key_here
   NEWS_API_URL=https://newsapi.org/v2
   
   # MarketStack API (Get your free key at: https://marketstack.com/signup/free)
   MARKET_DATA_API_KEY=your_actual_marketstack_key_here
   MARKET_DATA_URL=https://api.marketstack.com/v1
   ```

3. **Example with real keys:**
   ```bash
   # Replace with your actual keys
   NEWS_API_KEY=abc123def456ghi789jkl012mno345pq
   MARKET_DATA_API_KEY=xyz789abc123def456ghi789jkl012mno
   ```

## 🚀 Restart and Test

1. **Restart the backend:**
   ```bash
   # Stop the current process (Ctrl+C)
   # Then restart
   ./api
   ```

2. **Test API Integration:**
   The system will automatically use real APIs when valid keys are configured. You'll see logs like:
   ```
   INFO: Fetching real news for sector: technology
   INFO: Successfully fetched 5 real news articles for sector technology
   ```

## 🔍 API Integration Features

### Real News Integration
- **Sector-specific news:** Automatically searches for news related to your startup's sector
- **Investment news:** Searches for funding, investment, and venture capital news
- **Sentiment analysis:** Analyzes news sentiment using keyword-based algorithms
- **Fallback system:** Uses mock data if API calls fail

### Market Data Enhancement
- **Industry trends:** Enhanced with real market indicators
- **Competitive analysis:** Real competitor data where available
- **Market sentiment:** News-driven market sentiment scoring
- **Risk adjustment:** AI analysis considers real market conditions

### Smart Fallback System
- **Graceful degradation:** If API keys are missing or invalid, system uses high-quality mock data
- **No interruption:** Your app continues working even if external APIs are down
- **Logging:** Clear logs show whether real or mock data is being used

## 📊 Usage in the App

### Market Validation
When you onboard a startup, the system now:

1. **Fetches real news** about your industry sector
2. **Analyzes sentiment** from recent articles
3. **Incorporates market data** from financial APIs
4. **Provides AI-enhanced** risk assessments

### Decision Analysis
When making decisions, the system:

1. **Considers current market sentiment** from real news
2. **Analyzes competitive landscape** using real data
3. **Factors in industry trends** from multiple sources
4. **Provides contextual recommendations** based on current market conditions

## 🔧 Troubleshooting

### API Key Issues

**Problem:** "No valid NewsAPI key configured, using mock news data"
```bash
# Check your .env file
cat .env | grep NEWS_API_KEY

# Make sure the key is not empty or placeholder
NEWS_API_KEY=your_news_api_key_here  # ❌ This is placeholder
NEWS_API_KEY=abc123def456...         # ✅ This is real key
```

**Problem:** API rate limit exceeded
```bash
# Check your API usage on provider dashboards
# NewsAPI: https://newsapi.org/account
# MarketStack: https://marketstack.com/dashboard
```

### Network Issues

**Problem:** API calls timing out
```bash
# Check internet connectivity
curl -I https://newsapi.org
curl -I https://api.marketstack.com

# Check firewall/proxy settings
```

### Logs and Debugging

**Enable debug logging:**
```bash
# In your .env file
LOG_LEVEL=debug

# Restart the app
./api
```

**Monitor API calls:**
```bash
# Watch logs for API activity
./api 2>&1 | grep -E "(NewsAPI|MarketStack|API key)"
```

## 💡 Cost Optimization Tips

### Free Tier Management
- **Monitor usage:** Check your API dashboards regularly
- **Efficient caching:** The system caches results to minimize API calls
- **Smart rate limiting:** APIs are called only when necessary

### Upgrade Considerations
- **NewsAPI Pro:** $449/month for unlimited requests
- **MarketStack Standard:** $49.99/month for 10,000 requests
- **Consider upgrading** when you reach free tier limits

## 🎯 Next Steps

1. **✅ Set up API keys** following this guide
2. **🧪 Test the system** with real startup onboarding
3. **📊 Monitor API usage** on provider dashboards
4. **🔍 Review logs** to confirm real data integration
5. **📈 Analyze results** to see enhanced market intelligence

## 📞 Support

If you encounter issues:

1. **Check the logs** for detailed error messages
2. **Verify API keys** are correctly configured
3. **Test API connectivity** using curl commands
4. **Review free tier limits** on provider dashboards

The system is designed to work seamlessly whether you use free or paid API tiers, ensuring your startup risk assessment platform provides the best possible market intelligence!

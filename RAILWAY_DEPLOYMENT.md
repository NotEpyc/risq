# 🚀 Railway Deployment Guide

## Prerequisites
1. Railway account (sign up at https://railway.app)
2. Git repository with your code
3. API keys for OpenAI, NewsAPI, and MarketStack

## Step-by-Step Deployment

### 1. Prepare Your Repository
```bash
# Make sure all files are committed
git add .
git commit -m "Prepare for Railway deployment"
git push origin main
```

### 2. Create Railway Project
1. Go to https://railway.app
2. Click "New Project"
3. Choose "Deploy from GitHub repo"
4. Select your repository
5. Railway will automatically detect it's a Go project

### 3. Add Required Services
Your app needs these services:

#### A. Add PostgreSQL Database
1. In your Railway project dashboard
2. Click "New Service" → "Database" → "PostgreSQL"
3. Railway will automatically create the database and set environment variables

#### B. Add Redis
1. Click "New Service" → "Database" → "Redis"
2. Railway will automatically create Redis and set environment variables

#### C. Add NATS (Optional - for event streaming)
1. Click "New Service" → "Template" → Search for "NATS"
2. Or use Railway's community templates

### 4. Configure Environment Variables
In your Railway project:

1. Go to your app service (not database services)
2. Click on "Variables" tab
3. Add these environment variables:

#### Required Variables:
```env
# Application
APP_NAME=Smart Risk Assessment API
APP_HOST=0.0.0.0
APP_ENV=production

# Database (Railway auto-populates these)
DB_SSL_MODE=require

# JWT Secret (IMPORTANT: Use a strong 64-character secret)
JWT_SECRET=your_super_secure_jwt_secret_key_change_in_production_use_64_chars

# OpenAI API
OPENAI_API_KEY=sk-proj-your-actual-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini
OPENAI_TEMPERATURE=0.3
OPENAI_MAX_TOKENS=1000

# External APIs
MARKET_DATA_API_KEY=your_marketstack_api_key_here
MARKET_DATA_URL=https://api.marketstack.com/v1
NEWS_API_KEY=your_newsapi_key_here
NEWS_API_URL=https://newsapi.org/v2

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# NATS Configuration
NATS_URL=nats://nats-service:4222
NATS_CLUSTER_ID=risq-cluster
NATS_CLIENT_ID=risq-api
NATS_MAX_RECONNECTS=10
NATS_RECONNECT_WAIT=2s
NATS_CONNECTION_TIMEOUT=5s
```

### 5. Configure Build Settings
1. In your service settings, go to "Settings" tab
2. Set the following:

#### Build Configuration:
- **Build Command**: (Leave empty - Railway will auto-detect)
- **Start Command**: (Leave empty - uses CMD from Dockerfile)
- **Dockerfile Path**: `railway.Dockerfile`

#### Port Configuration:
- Railway automatically sets the `PORT` environment variable
- Your app reads this via `$PORT` in the environment

### 6. Deploy
1. Railway will automatically start building when you push to GitHub
2. Monitor the build logs in the "Deployments" tab
3. Once deployed, Railway will provide you with a public URL

### 7. Verify Deployment
Test your endpoints:
```bash
# Health check
curl https://your-app-name.railway.app/health

# API test (replace with your Railway URL)
curl -X POST https://your-app-name.railway.app/api/v1/auth/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"password123"}'
```

## Important Notes

### Security
- ✅ **Never commit API keys** to your repository
- ✅ **Use strong JWT secrets** (64+ characters)
- ✅ **Enable SSL mode** for database connections in production
- ✅ **Rotate API keys** regularly

### Performance
- Railway provides **512MB RAM** and **1 vCPU** on the free tier
- Upgrade to **Pro Plan** for production workloads
- Monitor **database connections** - Railway PostgreSQL has connection limits

### Monitoring
- Use Railway's built-in **logs and metrics**
- Set up **alerts** for critical errors
- Monitor **API usage** for OpenAI/NewsAPI to avoid rate limits

## Troubleshooting

### Build Failures
1. Check the build logs in Railway dashboard
2. Ensure `go.mod` and `go.sum` are committed
3. Verify Dockerfile syntax

### Database Connection Issues
1. Verify PostgreSQL service is running
2. Check environment variables are set correctly
3. Ensure SSL mode is set to `require`

### API Key Issues
1. Verify all API keys are valid and active
2. Check API quotas and billing status
3. Test API keys locally first

### Redis Connection Issues
1. Verify Redis service is running
2. Check Redis connection string format
3. Ensure Redis password is set if required

## Production Checklist
- [ ] Strong JWT secret (64+ characters)
- [ ] Valid API keys for all external services
- [ ] Database SSL mode enabled
- [ ] Appropriate log level (info/warn)
- [ ] Health check endpoint working
- [ ] All environment variables set
- [ ] CORS configuration for your frontend domain
- [ ] API rate limiting configured
- [ ] Monitoring and alerts set up

## Cost Optimization
- **Free Tier**: Good for development and testing
- **Pro Tier ($20/month)**: Required for production
- **Database**: Consider Railway's database pricing
- **API Costs**: Monitor OpenAI/NewsAPI usage

## Next Steps After Deployment
1. Update your Flutter app to use the Railway URL
2. Set up custom domain (optional)
3. Configure CDN for static assets (if any)
4. Set up monitoring and alerts
5. Plan backup and disaster recovery

Your backend is now ready for production! 🎉

# Railway Deployment Guide for Risk Assessment Backend

## Prerequisites

1. **Railway Account**: Sign up at [railway.app](https://railway.app)
2. **GitHub Repository**: Push your code to GitHub
3. **API Keys**: Prepare your OpenAI, NewsAPI, and MarketStack API keys

## Step-by-Step Deployment

### 1. Create New Railway Project

1. Go to [railway.app](https://railway.app) and click "Start a New Project"
2. Choose "Deploy from GitHub repo"
3. Connect your GitHub account and select your repository
4. Railway will automatically detect the Go application

### 2. Configure Build Settings

1. In your Railway project dashboard, go to Settings
2. Under "Build", set:
   - **Build Command**: `go build -o main ./cmd/api/main.go`
   - **Start Command**: `./main`
   - **Dockerfile Path**: `railway.Dockerfile` (if using custom Dockerfile)

### 3. Add Required Services

#### Add PostgreSQL Database
1. Click "New" → "Database" → "Add PostgreSQL"
2. Railway will automatically provide `DATABASE_URL`
3. Note: You can also set individual DB variables if needed

#### Add Redis (Optional but Recommended)
1. Click "New" → "Database" → "Add Redis"
2. Railway will automatically provide `REDIS_URL`

### 4. Configure Environment Variables

Go to your project → Settings → Variables and add:

```bash
# Application
APP_NAME=Smart Risk Assessment API
APP_HOST=0.0.0.0
APP_PORT=8080
APP_ENV=production

# Database (if not using DATABASE_URL)
DB_HOST=${{Postgres.PGHOST}}
DB_PORT=${{Postgres.PGPORT}}
DB_USER=${{Postgres.PGUSER}}
DB_PASSWORD=${{Postgres.PGPASSWORD}}
DB_NAME=${{Postgres.PGDATABASE}}
DB_SSL_MODE=require

# Redis (if not using REDIS_URL)
REDIS_HOST=${{Redis.REDIS_HOST}}
REDIS_PORT=${{Redis.REDIS_PORT}}
REDIS_PASSWORD=${{Redis.REDIS_PASSWORD}}
REDIS_DB=0

# OpenAI (REQUIRED)
OPENAI_API_KEY=your-openai-api-key-here
OPENAI_MODEL=gpt-4o-mini
OPENAI_TEMPERATURE=0.3
OPENAI_MAX_TOKENS=1000

# JWT (REQUIRED)
JWT_SECRET=your-super-secure-jwt-secret-key-for-production
JWT_EXPIRES_IN=24h

# Logging
LOG_LEVEL=info
LOG_FORMAT=json

# External APIs (REQUIRED)
MARKET_DATA_API_KEY=your-marketstack-api-key
MARKET_DATA_URL=https://api.marketstack.com/v1
NEWS_API_KEY=your-newsapi-key
NEWS_API_URL=https://newsapi.org/v2

# NATS (External service needed)
NATS_URL=nats://your-external-nats-server:4222
NATS_CLUSTER_ID=risq-cluster
NATS_CLIENT_ID=risq-api
```

### 5. Deploy

1. Click "Deploy" in Railway
2. Monitor the build logs for any issues
3. Once deployed, Railway will provide a public URL

### 6. Test Deployment

Test your API endpoints:
```bash
curl https://your-railway-app.railway.app/health
```

## Important Notes

### NATS Server
Railway doesn't provide NATS as a managed service. You have options:
1. **Use Railway's NATS template** from the template gallery
2. **Use external NATS service** (NATS Cloud, DigitalOcean, etc.)
3. **Deploy NATS separately** on Railway using Docker

### Database Migrations
The app automatically runs migrations on startup, so no manual migration needed.

### SSL/HTTPS
Railway automatically provides HTTPS for your deployed application.

### Monitoring
Use Railway's built-in monitoring and logs to track your application.

## Troubleshooting

### Build Fails
- Check go.mod and go.sum are committed
- Ensure all dependencies are properly listed
- Check Railway build logs for specific errors

### Database Connection Issues
- Verify DATABASE_URL or individual DB variables
- Ensure SSL mode is set correctly (require for production)

### External API Issues
- Verify all API keys are set correctly
- Check API key permissions and billing status

### NATS Connection Issues
- Ensure NATS server is accessible from Railway
- Check NATS URL format and credentials

## Cost Optimization

1. **Hobby Plan**: Free with limitations
2. **Pro Plan**: $5/month per user
3. **Database Costs**: PostgreSQL (~$5/month), Redis (~$3/month)

## Security Considerations

1. Use strong JWT secrets (32+ characters)
2. Enable SSL for database connections
3. Use environment variables for all secrets
4. Monitor API usage and set billing limits
5. Regularly rotate API keys

## Post-Deployment

1. Test all API endpoints
2. Verify event-driven flows work
3. Check logs for any errors
4. Set up monitoring alerts
5. Configure automatic backups for database

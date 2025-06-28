#!/bin/bash

# Railway Deployment Setup Script
echo "🚀 Setting up Railway deployment for Risk Assessment Backend..."

# Check if git is initialized
if [ ! -d ".git" ]; then
    echo "❌ Error: This is not a git repository. Please run 'git init' first."
    exit 1
fi

# Check if files are committed
if [ -n "$(git status --porcelain)" ]; then
    echo "⚠️  Warning: You have uncommitted changes. Please commit them first:"
    echo "   git add ."
    echo "   git commit -m 'Prepare for Railway deployment'"
    echo "   git push origin main"
    exit 1
fi

echo "✅ Git repository is clean and ready for deployment"

# Instructions for Railway deployment
echo ""
echo "📋 Railway Deployment Instructions:"
echo ""
echo "1. 🌐 Go to https://railway.app and sign in"
echo "2. 📝 Click 'New Project' → 'Deploy from GitHub repo'"
echo "3. 🔗 Select this repository"
echo "4. 🗄️  Add services:"
echo "   - PostgreSQL database"
echo "   - Redis database" 
echo "   - NATS (optional, for events)"
echo ""
echo "5. ⚙️  Set environment variables (see railway.env.template):"
echo "   📌 REQUIRED:"
echo "   - OPENAI_API_KEY=your_openai_key"
echo "   - JWT_SECRET=your_64_char_secret"
echo "   - MARKET_DATA_API_KEY=your_marketstack_key"
echo "   - NEWS_API_KEY=your_newsapi_key"
echo ""
echo "6. 🔧 In service settings:"
echo "   - Dockerfile Path: railway.Dockerfile"
echo "   - Port: Railway auto-detects from \$PORT"
echo ""
echo "7. 🚀 Deploy!"
echo ""
echo "📖 For detailed instructions, see: RAILWAY_DEPLOYMENT.md"
echo ""
echo "🎉 Your backend will be live at: https://your-app-name.railway.app"

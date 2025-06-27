#!/bin/bash

# Risk Assessment Backend - Quick Setup Script
echo "🚀 Setting up Risk Assessment Backend API..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "📝 Creating .env file from template..."
    cp .env.example .env
    echo "✅ .env file created. Please edit it with your actual values:"
    echo "   - Set your OpenAI API key"
    echo "   - Update JWT secret"
    echo "   - Configure database credentials if needed"
else
    echo "✅ .env file already exists"
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.19+ from https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed: $(go version)"

# Install dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

# Check if Docker is installed
if command -v docker &> /dev/null; then
    echo "✅ Docker is available"
    echo "🐳 You can use 'docker-compose up -d' to start PostgreSQL and Redis"
else
    echo "⚠️  Docker not found. You'll need to set up PostgreSQL and Redis manually."
fi

# Build the application
echo "🔨 Building the application..."
go build -o main cmd/api/main.go

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo ""
    echo "🎉 Setup complete! Next steps:"
    echo "   1. Edit .env file with your configuration"
    echo "   2. Start PostgreSQL and Redis (docker-compose up -d)"
    echo "   3. Run the application: ./main"
    echo ""
    echo "📚 API will be available at: http://localhost:8080"
    echo "📄 Check ARCHITECTURE.md for detailed API documentation"
else
    echo "❌ Build failed. Please check for errors above."
    exit 1
fi

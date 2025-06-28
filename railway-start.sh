#!/bin/sh

# Railway startup script for the Risk Assessment Backend
echo "Starting Risk Assessment Backend on Railway..."
echo "Environment: $RAILWAY_ENVIRONMENT"
echo "Port: $PORT"

# Run database migrations if needed
echo "Starting application..."

# Start the application
exec ./main

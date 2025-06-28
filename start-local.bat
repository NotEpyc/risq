@echo off
echo Starting RISQ Backend Local Development Environment...
echo.

REM Check if Docker is running
docker version >nul 2>&1
if %errorlevel% neq 0 (
    echo Error: Docker is not running. Please start Docker Desktop first.
    pause
    exit /b 1
)

echo Step 1: Starting supporting services (PostgreSQL, Redis, NATS)...
docker-compose -f docker-compose.local.yml up -d

echo.
echo Waiting for services to be ready...
timeout /t 10 /nobreak >nul

REM Check service health
echo Step 2: Checking service health...
docker-compose -f docker-compose.local.yml ps

echo.
echo Step 3: Building and starting Go application...
echo Loading environment from .env.local...

REM Set environment variables from .env.local
if exist .env.local (
    for /f "tokens=1,2 delims==" %%A in (.env.local) do (
        if not "%%A"=="" if not "%%A:~0,1%"=="#" set %%A=%%B
    )
)

echo Building Go application...
go build -o risq-backend.exe ./cmd/api

if %errorlevel% neq 0 (
    echo Error: Failed to build Go application
    pause
    exit /b 1
)

echo.
echo ================================================
echo RISQ Backend is starting...
echo ================================================
echo API will be available at: http://localhost:8080
echo.
echo Your computer's IP addresses:
ipconfig | findstr "IPv4"
echo.
echo For Flutter app, use one of the above IP addresses with port 8080
echo Example: http://192.168.1.100:8080
echo ================================================
echo.
echo Press Ctrl+C to stop the server
echo.

REM Start the application
risq-backend.exe

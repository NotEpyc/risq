@echo off
echo Stopping RISQ Backend Local Development Environment...
echo.

echo Stopping supporting services...
docker-compose -f docker-compose.local.yml down

echo.
echo Removing built executable...
if exist risq-backend.exe del risq-backend.exe

echo.
echo Local development environment stopped.
pause

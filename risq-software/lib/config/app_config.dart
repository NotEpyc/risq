class AppConfig {
  // Backend configuration
  static const String _defaultBackendUrl = 'https://resqbackend-production.up.railway.app';
  
  // Get backend URL from environment or use default
  static String get backendUrl {
    // In a real app, you would use flutter_dotenv or similar to read .env files
    // For now, we'll use the default URL
    // const String.fromEnvironment('BACKEND_URL', defaultValue: _defaultBackendUrl);
    return _defaultBackendUrl;
  }
  
  // API timeout configuration
  static const int apiTimeout = 30000; // 30 seconds
  
  // Debug configuration
  static const bool debugMode = bool.fromEnvironment('DEBUG_MODE', defaultValue: false);
  
  // Other app configurations
  static const String appName = 'RISQ';
  static const String appVersion = '1.0.0';
}

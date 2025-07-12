# JWT Token Management Implementation

This document explains how JWT tokens are handled throughout the RISQ Flutter application.

## Overview

The JWT token management system is implemented using a layered approach:

1. **AuthService** - Core authentication and token storage
2. **TokenManager** - Centralized token management utility
3. **API Services** - Services that use tokens for authenticated requests
4. **UI Components** - Widgets that handle authentication state

## Core Components

### 1. AuthService (`lib/services/auth_service.dart`)

The foundation layer that handles:
- User login and registration
- Automatic token storage upon successful authentication
- Token persistence using SharedPreferences
- Basic token and user data operations

Key methods:
```dart
// Authentication
AuthService.login(email, password)
AuthService.register(name, email, password)
AuthService.logout()

// Token management
AuthService.storeAuthToken(token)
AuthService.getAuthToken()
AuthService.clearAuthToken()

// User data management
AuthService.storeUserData(userData)
AuthService.getUserData()
AuthService.clearUserData()
```

### 2. TokenManager (`lib/utils/token_manager.dart`)

A utility class that provides enhanced token management:
- Caching for better performance
- Convenient access methods
- Authentication status checks
- User information extraction

Key methods:
```dart
// Token operations
TokenManager.storeToken(token, userData: userData)
TokenManager.getToken()
TokenManager.isAuthenticated()

// User information
TokenManager.getCurrentUserEmail()
TokenManager.getCurrentUserName()
TokenManager.getCurrentUserId()
TokenManager.getUserData()

// Utilities
TokenManager.getAuthHeaders()
TokenManager.clearAll()
TokenManager.refreshCache()
```

### 3. API Services

Services that automatically handle authentication:

#### StartupService (`lib/services/startup_service.dart`)
```dart
// Automatically uses stored JWT token
StartupService.onboardStartup(startupData: data)
StartupService.getStartupProfile()
```

#### ApiService (`lib/services/api_service.dart`)
```dart
// Generic authenticated API calls
ApiService.authenticatedGet('/api/v1/endpoint')
ApiService.authenticatedPost('/api/v1/endpoint', data)
ApiService.getCurrentUserProfile()
```

## Authentication Flow

### 1. Login/Registration
```dart
// In signin_page.dart or signup_page.dart
final result = await AuthService.login(email: email, password: password);
if (result['success'] == true) {
  // Token is automatically stored by AuthService
  // Navigate to main app
  Navigator.pushReplacement(context, MaterialPageRoute(...));
}
```

### 2. Making Authenticated API Calls
```dart
// Option 1: Using StartupService (recommended for startup-related APIs)
final result = await StartupService.onboardStartup(startupData: data);

// Option 2: Using ApiService for generic APIs
final result = await ApiService.authenticatedGet('/api/v1/user/profile');

// Option 3: Manual token handling
final token = await TokenManager.getToken();
final headers = await TokenManager.getAuthHeaders();
```

### 3. Checking Authentication Status
```dart
// Quick check
final isAuth = await TokenManager.isAuthenticated();

// Detailed status with user info
final authStatus = await ApiService.getAuthStatus();

// Get user information
final userName = await TokenManager.getCurrentUserName();
final userEmail = await TokenManager.getCurrentUserEmail();
```

### 4. Logout
```dart
// Clear all authentication data
await TokenManager.clearAll();
// Navigate to login screen
Navigator.pushNamedAndRemoveUntil(context, '/login', (route) => false);
```

## Usage Examples

### In UI Components

```dart
class MyWidget extends StatefulWidget {
  @override
  _MyWidgetState createState() => _MyWidgetState();
}

class _MyWidgetState extends State<MyWidget> {
  String _userGreeting = 'Loading...';
  
  @override
  void initState() {
    super.initState();
    _loadUserInfo();
  }
  
  Future<void> _loadUserInfo() async {
    if (await TokenManager.isAuthenticated()) {
      final greeting = await ApiService.getUserGreeting();
      setState(() => _userGreeting = greeting);
    } else {
      // Redirect to login
      Navigator.pushNamed(context, '/login');
    }
  }
  
  @override
  Widget build(BuildContext context) {
    return Text(_userGreeting);
  }
}
```

### In Services

```dart
class MyApiService {
  static Future<Map<String, dynamic>> getMyData() async {
    final token = await TokenManager.getToken();
    
    if (token == null) {
      return {'success': false, 'message': 'Authentication required'};
    }
    
    final response = await http.get(
      Uri.parse('${baseUrl}/api/v1/my-endpoint'),
      headers: {
        'Authorization': 'Bearer $token',
        'Content-Type': 'application/json',
      },
    );
    
    // Handle response...
  }
}
```

### Authentication Guard

```dart
class ProtectedPage extends StatelessWidget {
  @override
  Widget build(BuildContext context) {
    return AuthGuard(
      child: Scaffold(
        appBar: AppBar(title: Text('Protected Content')),
        body: Column(
          children: [
            UserInfoWidget(),
            // Other protected content...
          ],
        ),
      ),
    );
  }
}
```

## Troubleshooting

### Authentication Error on Onboarding

If you're getting "Authentication error. Please try again" when clicking "Complete Setup" in the onboarding page:

1. **Check Token Storage**: Verify that the token is being stored after login/signup
2. **Debug Token Flow**: Add the AuthDebugWidget temporarily to check token status
3. **Verify API Response**: Check if the authentication server is returning a valid JWT token
4. **Check Network Issues**: Ensure the app can communicate with the backend server

#### Debug Steps:

1. **Add Debug Widget** (temporarily):
```dart
// In your onboarding page or any page
import 'package:risq/widgets/auth_debug_widget.dart';

// Add to the widget tree
Column(
  children: [
    AuthDebugWidget(), // Remove in production
    // ... other widgets
  ],
)
```

2. **Check Console Logs**: Look for authentication debugging messages:
   - "Registration successful, navigating to onboarding..."
   - "Token stored after registration: Yes/No"
   - "StartupService: Token available: Yes/No"

3. **Common Solutions**:
   - Ensure the backend returns a valid JWT token
   - Check that the token storage is not being cleared
   - Verify the API endpoint is correct
   - Ensure network connectivity

#### Fallback Mechanism:

The app now includes a fallback mechanism that:
- Passes the token directly from login/signup to onboarding
- Automatically retries authentication if the first attempt fails
- Provides detailed debugging information

## Security Considerations

1. **Token Storage**: Tokens are stored using SharedPreferences, which is acceptable for development but consider using flutter_secure_storage for production.

2. **Token Expiration**: The current implementation includes basic token validation. For production, implement proper JWT decoding and expiration checking.

3. **Token Refresh**: Consider implementing automatic token refresh when tokens expire.

4. **Network Security**: Always use HTTPS in production and validate SSL certificates.

## File Locations

- `lib/services/auth_service.dart` - Core authentication service
- `lib/utils/token_manager.dart` - Token management utility
- `lib/services/startup_service.dart` - Startup-specific authenticated APIs
- `lib/services/api_service.dart` - Generic authenticated API service
- `lib/widgets/auth_widgets.dart` - Authentication UI components
- `lib/screens/login/signin_page.dart` - Login page with token storage
- `lib/screens/login/signup_page.dart` - Registration page with token storage
- `lib/screens/pages/onboarding_page.dart` - Onboarding with authentication verification

## Best Practices

1. **Always check authentication** before making API calls
2. **Use TokenManager** for consistent token access
3. **Handle authentication errors** gracefully
4. **Clear tokens on logout** to ensure security
5. **Cache user data** for better performance
6. **Use AuthGuard** for protected pages
7. **Verify authentication** in critical flows

## Future Enhancements

1. Implement automatic token refresh
2. Add biometric authentication support
3. Implement secure token storage for production
4. Add token validation and expiration handling
5. Implement offline authentication caching
6. Add multi-factor authentication support

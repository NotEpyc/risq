# Base URL Centralization Changes

## Summary
Successfully centralized all Railway backend URLs to use a single source of truth in `AuthService`.

## Changes Made

### 1. **Modified `lib/services/auth_service.dart`**
- Added a public static getter `baseUrl` to expose the base URL
- Kept the private `_baseUrl` constant with the Railway URL
- Other services can now access the URL via `AuthService.baseUrl`

```dart
class AuthService {
  // Railway app URL with HTTPS
  static const String _baseUrl = 'https://resqbackend-production.up.railway.app';
  
  // Public getter for base URL to be used by other services
  static String get baseUrl => _baseUrl;
  
  // ... rest of the class
}
```

### 2. **Updated `lib/services/health_service.dart`**
- Removed hardcoded Railway URL
- Added import for `auth_service.dart`
- Changed to use `AuthService.baseUrl`

```dart
import 'package:risq/services/auth_service.dart';

class HealthService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  // ... rest of the class
}
```

### 3. **Updated `lib/services/startup_service.dart`**
- Removed hardcoded Railway URL
- Changed to use `AuthService.baseUrl`
- Already had the auth_service import

```dart
class StartupService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  // ... rest of the class
}
```

### 4. **Updated `lib/services/speculation_service.dart`**
- Removed hardcoded Railway URL
- Changed to use `AuthService.baseUrl`
- Already had the auth_service import

```dart
class SpeculationService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  // ... rest of the class
}
```

### 5. **Updated `lib/services/decision_service.dart`**
- Removed hardcoded Railway URL
- Changed to use `AuthService.baseUrl`
- Already had the auth_service import

```dart
class DecisionService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  // ... rest of the class
}
```

## Benefits

### ✅ **Single Source of Truth**
- The Railway URL `https://resqbackend-production.up.railway.app` is now only defined in one place
- Easy to update for different environments (dev, staging, production)

### ✅ **Better Maintainability**
- When the Railway URL changes, only `auth_service.dart` needs to be updated
- Reduces risk of inconsistent URLs across services

### ✅ **No Breaking Changes**
- All service methods continue to work exactly the same
- No changes needed in UI code or other parts of the app

### ✅ **Clean Architecture**
- AuthService naturally becomes the central configuration point
- Other services depend on AuthService, which makes sense since they all need authentication

## Verification

### ✅ **Compilation Check**
```bash
flutter analyze lib/services/
# Result: No issues found!
```

### ✅ **No Hardcoded URLs**
Verified that `resqbackend-production.up.railway.app` only appears in:
- `lib/services/auth_service.dart` (the source)
- Documentation files (not actual code)

### ✅ **All Services Updated**
- ✅ health_service.dart
- ✅ startup_service.dart  
- ✅ speculation_service.dart
- ✅ decision_service.dart
- ✅ auth_service.dart (source)

## Future Usage

To change the base URL for all services (e.g., for different environments):

```dart
// In auth_service.dart - change only this line:
static const String _baseUrl = 'https://your-new-backend-url.com';
```

All other services will automatically use the new URL without any code changes.

## Files Modified
1. `lib/services/auth_service.dart` - Added public baseUrl getter
2. `lib/services/health_service.dart` - Updated to use AuthService.baseUrl
3. `lib/services/startup_service.dart` - Updated to use AuthService.baseUrl
4. `lib/services/speculation_service.dart` - Updated to use AuthService.baseUrl
5. `lib/services/decision_service.dart` - Updated to use AuthService.baseUrl

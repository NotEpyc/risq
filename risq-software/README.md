# RISQ - Startup Risk Assessment Platform

[![Flutter](https://img.shields.io/badge/Flutter-3.7.0-blue.svg)](https://flutter.dev/)
[![Dart](https://img.shields.io/badge/Dart-3.7.0-blue.svg)](https://dart.dev/)
[![License](https://img.shields.io/badge/License-Private-red.svg)](LICENSE)

RISQ is a comprehensive Flutter application designed to help startups assess and manage their business risks through AI-powered analysis. The platform provides intelligent risk scoring, personalized recommendations, and decision-making tools to support entrepreneurial ventures.

## 🚀 Features

### 🔐 Authentication & Onboarding
- **Secure Registration/Login**: JWT-based authentication with secure credential management
- **Interactive Onboarding**: Comprehensive startup profile creation with guided steps
- **Profile Management**: Detailed startup information collection and management

### 📊 Risk Assessment
- **AI-Powered Risk Scoring**: Intelligent analysis with confidence levels (0-10 scale)
- **Visual Risk Meter**: Interactive Syncfusion gauge displaying current risk status
- **Dynamic Risk Levels**: Color-coded risk categories (Low, Medium, High, Critical)

### 🔔 Notification System
- **Real-time Risk Alerts**: Immediate notifications for critical risk changes
- **AI Suggestions**: Personalized recommendations based on risk analysis
- **Priority-based Organization**: Smart sorting by risk priority and timestamp
- **Interactive Management**: Mark notifications as read/unread, pull-to-refresh

### 💡 Decision Support Tools
- **Speculation Engine**: Test potential business decisions before implementation
- **Impact Analysis**: Understand how decisions might affect your risk profile
- **Data Visualization**: Clear charts and graphs for decision insights

### 📱 Modern UI/UX
- **Responsive Design**: Optimized for all screen sizes and orientations
- **Material Design 3**: Modern, intuitive interface following Google's design principles
- **Smooth Animations**: Fluid transitions and engaging user interactions
- **Video Backgrounds**: Dynamic visual elements for enhanced user experience

## 🛠️ Technology Stack

### Frontend (Flutter)
- **Framework**: Flutter 3.7.0
- **Language**: Dart 3.7.0
- **UI Components**: Material Design 3, Custom widgets
- **State Management**: StatefulWidget with modern Flutter patterns
- **Navigation**: MaterialPageRoute with smooth transitions

### Key Dependencies
```yaml
# Core Dependencies
http: ^1.1.0                    # HTTP client for API communication
video_player: ^2.8.1           # Video background support
shared_preferences: ^2.2.2     # Local data persistence

# UI/UX Libraries
syncfusion_flutter_gauges: ^27.1.58  # Interactive risk meter
google_nav_bar: ^5.0.7               # Modern navigation bar
font_awesome_flutter: ^10.6.0        # Icon library
cupertino_icons: ^1.0.8              # iOS-style icons
```

### Backend Integration
- **API**: RESTful services hosted on Railway
- **Authentication**: JWT token-based security
- **Data Format**: JSON for all API communications
- **Endpoints**: Centralized URL management through AppConfig

## 📁 Project Structure

```
lib/
├── config/
│   └── app_config.dart          # Centralized configuration management
├── screens/
│   ├── login/
│   │   ├── signin_page.dart     # User authentication
│   │   └── signup_page.dart     # User registration
│   ├── pages/
│   │   ├── home_page.dart       # Main dashboard
│   │   ├── notifications_page.dart   # Notification center
│   │   ├── startup_profile_page.dart # Profile management
│   │   ├── onboarding_page.dart      # Initial setup
│   │   ├── speculation_page.dart     # Decision testing
│   │   ├── decision_page.dart        # Decision analysis
│   │   └── data_display_page.dart    # Data visualization
│   ├── main_navigation.dart     # Bottom navigation
│   └── splash_screen.dart       # App loading screen
├── services/
│   ├── auth_service.dart        # Authentication management
│   ├── notification_service.dart    # Notification handling
│   ├── startup_service.dart     # Profile data management
│   ├── speculation_service.dart     # Decision analysis
│   ├── decision_service.dart    # Decision processing
│   └── health_service.dart      # System health checks
├── theme/
│   ├── theme.dart              # App-wide theming
│   └── themed_text_selection.dart   # Custom text selection
├── utils/
│   └── responsive_utils.dart    # Responsive design utilities
├── widgets/
│   └── themed_text_field.dart   # Custom form components
└── main.dart                   # Application entry point

assets/
└── videos/
    └── bg_video.mp4            # Background video for auth screens

android/                        # Android-specific configurations
ios/                           # iOS-specific configurations
```

## 🚀 Getting Started

### Prerequisites
- **Flutter SDK**: 3.7.0 or higher
- **Dart SDK**: 3.7.0 or higher
- **Android Studio** or **VS Code** with Flutter extensions
- **Android SDK** (for Android development)
- **Xcode** (for iOS development, macOS only)

### Installation

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd risq
   ```

2. **Install dependencies**
   ```bash
   flutter pub get
   ```

3. **Configure environment**
   ```bash
   # Copy environment template
   cp .env.example .env
   
   # Edit .env with your configuration
   # BACKEND_URL=https://your-backend-url.railway.app
   ```

4. **Run the application**
   ```bash
   # Check for any issues
   flutter doctor
   
   # Run on connected device/emulator
   flutter run
   
   # Run in debug mode
   flutter run --debug
   
   # Run in release mode
   flutter run --release
   ```

### Build for Production

```bash
# Android APK
flutter build apk --release

# Android App Bundle
flutter build appbundle --release

# iOS (macOS only)
flutter build ios --release

# Web
flutter build web --release

# Windows
flutter build windows --release
```

## 🔧 Configuration

### Environment Variables
Create a `.env` file based on `.env.example`:

```env
# Backend API Configuration
BACKEND_URL=https://resqbackend-production.up.railway.app

# Development URLs (uncomment when needed)
# BACKEND_URL=http://localhost:3000
# BACKEND_URL=http://127.0.0.1:3000

# Other configuration
# API_TIMEOUT=30000
# DEBUG_MODE=true
```

### AppConfig Settings
Modify `lib/config/app_config.dart` for app-wide settings:

```dart
class AppConfig {
  static String get backendUrl => 'your-backend-url';
  static const int apiTimeout = 30000;
  static const bool debugMode = bool.fromEnvironment('DEBUG_MODE');
}
```

## 🔒 Security

### Protected Files
The following files are automatically ignored by Git for security:

- `android/local.properties` - Local Android SDK paths
- `android/key.properties` - Android signing keys
- `android/app/google-services.json` - Firebase Android config
- `ios/Runner/GoogleService-Info.plist` - Firebase iOS config
- `.env` - Environment variables
- Build artifacts and temporary files

### Best Practices
- ✅ No hardcoded API keys or credentials
- ✅ Environment-based configuration
- ✅ Secure JWT token handling
- ✅ Platform-specific credential protection
- ✅ Comprehensive `.gitignore` configuration

## 📱 Supported Platforms

| Platform | Status | Notes |
|----------|--------|-------|
| Android | ✅ Fully Supported | API 21+ (Android 5.0) |
| iOS | ✅ Fully Supported | iOS 12.0+ |
| Web | ⚠️ Limited Support | Basic functionality |
| Windows | ⚠️ Limited Support | Desktop features limited |
| macOS | ⚠️ Limited Support | Desktop features limited |
| Linux | ⚠️ Limited Support | Desktop features limited |

## 🎨 Design System

### Theme Colors
```dart
// Primary brand colors
authPrimaryColor: Color(0xFF2E3440)     // Dark blue-gray
authSecondaryColor: Color(0xFF3B4252)   // Medium blue-gray
authAccentColor: Color(0xFF5E81AC)      // Blue accent
authTertiaryColor: Color(0xFF88C0D0)    // Light blue

// UI colors
authBackgroundColor: Color(0xFFF8F9FA)  // Light background
authTextColor: Color(0xFF2E3440)        // Primary text
authErrorColor: Color(0xFFBF616A)       // Error red
```

### Responsive Design
- **Small Phone**: < 360px width
- **Phone**: 360px - 600px width
- **Tablet**: 600px - 1024px width
- **Desktop**: > 1024px width

## 🔌 API Integration

### Authentication Endpoints
```
POST /api/v1/auth/login     - User login
POST /api/v1/auth/signup    - User registration
POST /api/v1/auth/logout    - User logout
```

### Risk Assessment Endpoints
```
GET  /api/v1/risk/current   - Get current risk profile
POST /api/v1/risk/analyze   - Trigger risk analysis
```

### Profile Management Endpoints
```
GET  /api/v1/profile        - Get user profile
PUT  /api/v1/profile        - Update user profile
POST /api/v1/profile/setup  - Initial profile setup
```

## 🐛 Debugging

### Common Issues

1. **Video playback issues**
   ```bash
   # Ensure video file exists
   ls assets/videos/bg_video.mp4
   
   # Check pubspec.yaml assets configuration
   flutter clean && flutter pub get
   ```

2. **Network connectivity**
   ```bash
   # Check internet permissions (Android)
   # Verify backend URL in AppConfig
   flutter run --verbose
   ```

3. **Build issues**
   ```bash
   # Clean build artifacts
   flutter clean
   flutter pub get
   flutter build <platform>
   ```

### Development Tools

```bash
# Analyze code for issues
flutter analyze

# Run tests
flutter test

# Check dependencies
flutter pub deps

# Hot reload during development
# Press 'r' in terminal or save file in IDE
```

## 📈 Performance

### Optimization Strategies
- **Lazy loading**: Pages and data loaded on demand
- **Image optimization**: Compressed assets and caching
- **State management**: Efficient widget rebuilding
- **Network caching**: HTTP response caching
- **Build optimization**: Release builds with tree shaking

### Monitoring
- **Error tracking**: Comprehensive error handling
- **Performance metrics**: Built-in Flutter performance tools
- **Network monitoring**: API response time tracking
- **User analytics**: Navigation and feature usage

## 🤝 Contributing

### Development Workflow
1. **Fork** the repository
2. **Create** a feature branch (`git checkout -b feature/amazing-feature`)
3. **Commit** your changes (`git commit -m 'Add amazing feature'`)
4. **Push** to the branch (`git push origin feature/amazing-feature`)
5. **Open** a Pull Request

### Code Standards
- Follow [Dart style guide](https://dart.dev/guides/language/effective-dart/style)
- Use meaningful variable and function names
- Add comments for complex logic
- Write tests for new features
- Ensure responsive design compatibility

### Commit Message Format
```
type(scope): description

feat(auth): add biometric authentication
fix(notifications): resolve push notification bug
docs(readme): update installation instructions
```

## 📄 Documentation

### Additional Documentation
- [`SECURITY_CREDENTIALS.md`](SECURITY_CREDENTIALS.md) - Security implementation guide
- [`NOTIFICATION_INTEGRATION.md`](NOTIFICATION_INTEGRATION.md) - Notification system documentation
- [`BASE_URL_CENTRALIZATION.md`](BASE_URL_CENTRALIZATION.md) - URL management guide
- [`STARTUP_PROFILE_IMPLEMENTATION.md`](STARTUP_PROFILE_IMPLEMENTATION.md) - Profile feature guide

### API Documentation
- Backend API documentation available at backend repository
- Postman collection for API testing
- OpenAPI/Swagger specifications

## 📞 Support

### Getting Help
- **Issues**: Create a GitHub issue for bugs or feature requests
- **Discussions**: Use GitHub discussions for questions and ideas
- **Documentation**: Check the docs folder for detailed guides

### Team Contact
- **Development Team**: [Contact information]
- **Product Manager**: [Contact information]
- **Technical Support**: [Contact information]

## 📅 Roadmap

### Upcoming Features
- [ ] **Real-time Notifications**: Push notifications for critical alerts
- [ ] **Advanced Analytics**: Detailed risk trend analysis
- [ ] **Team Collaboration**: Multi-user startup profiles
- [ ] **Export Features**: PDF reports and data export
- [ ] **Mobile Optimizations**: Enhanced mobile experience
- [ ] **Offline Support**: Limited offline functionality

### Version History
- **v1.0.0** (Current) - Initial release with core features
- **v0.9.0** - Beta release with notification system
- **v0.8.0** - Alpha release with basic risk assessment

## 📜 License

This project is private and proprietary. All rights reserved.

---

**Built with ❤️ using Flutter**

For more information about Flutter development, visit the [official documentation](https://docs.flutter.dev/).

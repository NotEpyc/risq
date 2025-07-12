# Startup Profile Page Implementation

## Overview
A comprehensive startup profile page that displays and allows editing of startup information retrieved from the RISQ backend API. This page follows the same design patterns as the provided UserProfilePage example, but is specifically tailored for startup data.

## Features

### 📋 Data Display
- **Startup Information**: Name, description, industry, location
- **Company Details**: Funding stage, team size, founded date, website
- **Profile Metadata**: Shows when profile was created
- **Beautiful UI**: Gradient backgrounds, rounded cards, and consistent theming

### ✏️ Edit Functionality
- **Toggle Edit Mode**: Tap the edit icon in the app bar to enable editing
- **Field Validation**: Basic validation for required fields
- **Save Changes**: Currently shows "coming soon" message (API update endpoint not yet available)

### 🔒 Authentication
- **Logout Functionality**: Secure logout with confirmation dialog
- **Auto-redirect**: Redirects to login page after logout
- **Token Management**: Uses AuthService for authentication handling

### 🎨 UI Components
- **Funding Stage Dropdown**: Displays all available funding stages (idea, pre-seed, seed, series A, B, C, IPO)
- **Date Picker**: Founded date selection with date picker
- **Responsive Design**: Adapts to different screen sizes
- **Loading States**: Shows loading animation while fetching data
- **Error Handling**: Displays user-friendly error messages

## File Structure
```
lib/screens/pages/startup_profile_page.dart
```

## API Integration

### Data Retrieval
Uses `StartupService.getStartupProfile()` to fetch startup data from:
```
GET /api/v1/startup/profile
```

Expected response format:
```json
{
  "success": true,
  "message": "Startup profile retrieved successfully",
  "data": {
    "id": "660e8400-e29b-41d4-a716-446655440000",
    "name": "EduTech Solutions",
    "description": "AI-powered personalized learning platform for K-12 students",
    "industry": "Education Technology",
    "funding_stage": "seed",
    "location": "San Francisco, CA",
    "founded_date": "2024-01-15T00:00:00Z",
    "team_size": 5,
    "website": "https://edutech-solutions.com",
    "created_at": "2025-06-28T10:30:00Z",
    "updated_at": "2025-06-28T10:30:00Z"
  }
}
```

### Authentication
Uses `AuthService.logout()` for secure logout functionality.

## Navigation Integration

### From Home Page
The profile page is accessible from the home page by tapping the profile picture:

```dart
// In home_page.dart
GestureDetector(
  onTap: () {
    Navigator.push(
      context,
      MaterialPageRoute(
        builder: (context) => const StartupProfilePage(),
      ),
    );
  },
  child: // Profile picture widget
)
```

### Direct Navigation
You can navigate to the profile page from any screen:

```dart
Navigator.push(
  context,
  MaterialPageRoute(builder: (context) => const StartupProfilePage()),
);
```

## Dependencies
- **flutter/material.dart**: Core Flutter UI components
- **flutter/services.dart**: System UI overlay styles
- **risq/theme/theme.dart**: App theme and colors
- **risq/services/startup_service.dart**: API service for startup data
- **risq/services/auth_service.dart**: Authentication service
- **risq/screens/login/signin_page.dart**: Login page for logout redirect
- **risq/utils/responsive_utils.dart**: Responsive design utilities

## Styling

### Color Scheme
- **Primary Colors**: Uses AppTheme.authSecondaryColor, authAccentColor, authDarkBlue
- **Background**: Light gradient using auth theme colors
- **Text**: Dark blue for headings, black87 for body text
- **Form Fields**: White background with light borders

### Responsive Design
- **Text Sizes**: Adapts based on ResponsiveUtils methods
- **Spacing**: Consistent spacing using theme guidelines
- **Layout**: Single scroll view with proper padding

## Error Handling

### Network Errors
- Connection timeout messages
- Server unavailable notifications
- Authentication errors with auto-redirect

### User Feedback
- Loading spinners during API calls
- Success/error snackbars
- Confirmation dialogs for critical actions

## Future Enhancements

### Profile Updates
When the backend API supports profile updates, the save functionality can be implemented:

```dart
// Future implementation
Future<void> _saveChanges() async {
  final updateData = {
    'name': _nameController.text,
    'description': _descriptionController.text,
    // ... other fields
  };
  
  final result = await StartupService.updateProfile(updateData);
  // Handle response
}
```

### Profile Photo
Profile photo upload functionality can be added similar to the original UserProfilePage example.

### Additional Fields
More startup-specific fields can be added as needed:
- Business model details
- Revenue information
- Team member details
- Investor information

## Usage Example

```dart
import 'package:risq/screens/pages/startup_profile_page.dart';

// Navigate to profile page
void openProfile() {
  Navigator.push(
    context,
    MaterialPageRoute(
      builder: (context) => const StartupProfilePage(),
    ),
  );
}
```

## Testing

The page has been successfully compiled and integrated into the app. To test:

1. **Build the app**: `flutter build apk --debug`
2. **Run the app**: `flutter run`
3. **Navigate**: Tap profile picture on home page
4. **Test features**: Try edit mode, logout, etc.

## Notes

- The profile update functionality shows a "coming soon" message since the backend API doesn't currently support profile updates
- All form validation and UI interactions work properly
- The page follows the same design patterns as other pages in the app
- Error handling is comprehensive and user-friendly

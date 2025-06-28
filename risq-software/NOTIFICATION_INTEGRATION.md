# Notification System Integration

## Overview
This document describes the implementation of a notification system that displays risk score notifications fetched from the backend's risk profile JSON response.

## Implementation Details

### Files Modified/Created

#### 1. NotificationService (`lib/services/notification_service.dart`)
- **Purpose**: Fetches risk data from `/api/v1/risk/current` endpoint and generates notifications
- **Key Features**:
  - Uses centralized base URL from `AuthService.baseUrl`
  - Generates different types of notifications based on risk data:
    - Risk score updates
    - High/critical risk alerts
    - AI-generated suggestions
    - Risk factor notifications
  - Handles authentication tokens automatically
  - Robust error handling for network issues

#### 2. NotificationsPage (`lib/screens/pages/notifications_page.dart`)
- **Purpose**: UI component for displaying risk-based notifications
- **Key Features**:
  - Modern, responsive design with loading states
  - Notifications sorted by priority (High → Medium → Low)
  - Color-coded priority indicators
  - Mark as read/unread functionality
  - Empty state when no notifications exist
  - Pull-to-refresh functionality
  - Error handling with retry options

#### 3. HomePage Integration (`lib/screens/pages/home_page.dart`)
- **Purpose**: Connect notification button to notifications page
- **Changes Made**:
  - Added import for `NotificationsPage`
  - Updated notification button `onTap` to navigate to `NotificationsPage`
  - Removed placeholder "coming soon" message

### Notification Types Generated

1. **Risk Score Updates**
   - Displays current risk score and level
   - Generated for all risk assessments

2. **High Priority Alerts**
   - Only shown for high/critical risk levels
   - Immediate attention required

3. **AI Suggestions**
   - Based on backend AI recommendations
   - Actionable advice for risk mitigation

4. **Risk Factor Notifications**
   - Detailed breakdown of risk contributing factors
   - Educational and informational

### Data Flow

1. User taps notification bell in home page
2. Navigation to `NotificationsPage`
3. `NotificationsPage` calls `NotificationService.getRiskNotifications()`
4. Service fetches data from `/api/v1/risk/current`
5. Service generates notifications based on risk data
6. UI displays notifications with proper formatting and interactions

### Error Handling

- Network connectivity issues
- Authentication failures
- Backend API errors
- Invalid response data

### Security Considerations

- Uses Bearer token authentication
- Centralized URL management
- No sensitive data stored locally
- Secure HTTP requests

## Usage

### Accessing Notifications
- Users can tap the notification bell icon in the top-right corner of the home page
- This opens the dedicated notifications page

### Notification Interactions
- View all risk-related notifications
- Mark notifications as read/unread
- Pull down to refresh notifications
- Automatic priority-based sorting

## Future Enhancements

1. **Real-time Updates**: Push notifications for critical risk changes
2. **Persistent Storage**: Local storage for read/unread states
3. **Custom Filtering**: Filter by notification type or priority
4. **Notification History**: Archive and search past notifications
5. **Backend Sync**: Sync read states with backend for multi-device support

## Technical Notes

- All notifications are generated client-side based on backend risk data
- The notification badge (red dot) is always visible for UI consistency
- Responsive design supports various screen sizes
- Uses Material Design 3 principles for modern UI

## API Dependencies

- **Endpoint**: `/api/v1/risk/current`
- **Authentication**: Bearer token required
- **Response Format**: Standard risk profile JSON with score, level, suggestions, and factors

This implementation provides a complete notification system that seamlessly integrates with the existing app architecture while maintaining code quality and user experience standards.

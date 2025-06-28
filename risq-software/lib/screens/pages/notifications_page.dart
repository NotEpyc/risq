import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/services/notification_service.dart';

class NotificationsPage extends StatefulWidget {
  const NotificationsPage({super.key});

  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  List<Map<String, dynamic>> _notifications = [];
  bool _isLoading = true;
  bool _hasError = false;
  String _errorMessage = '';

  @override
  void initState() {
    super.initState();
    _loadNotifications();
  }

  Future<void> _loadNotifications() async {
    try {
      setState(() {
        _isLoading = true;
        _hasError = false;
      });

      final result = await NotificationService.getRiskNotifications();
      
      if (result['success']) {
        setState(() {
          _notifications = List<Map<String, dynamic>>.from(result['data'] ?? []);
          _isLoading = false;
        });
      } else {
        setState(() {
          _hasError = true;
          _errorMessage = result['message'] ?? 'Failed to load notifications';
          _isLoading = false;
        });
      }
    } catch (e) {
      setState(() {
        _hasError = true;
        _errorMessage = 'An error occurred while loading notifications';
        _isLoading = false;
      });
    }
  }

  Future<void> _markAsRead(String notificationId, int index) async {
    final success = await NotificationService.markAsRead(notificationId);
    if (success && mounted) {
      setState(() {
        _notifications[index]['isRead'] = true;
      });
    }
  }

  Future<void> _markAllAsRead() async {
    final success = await NotificationService.markAllAsRead();
    if (success && mounted) {
      setState(() {
        for (var notification in _notifications) {
          notification['isRead'] = true;
        }
      });
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(
          content: Text('All notifications marked as read'),
          backgroundColor: Colors.green,
        ),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        systemOverlayStyle: SystemUiOverlayStyle.dark,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back,
            color: AppTheme.authAccentColor,
            size: 24,
          ),
          onPressed: () => Navigator.of(context).pop(),
        ),
        title: Text(
          'Notifications',
          style: AppTheme.headingTextStyle.copyWith(
            color: AppTheme.authDarkBlue,
            fontSize: ResponsiveUtils.getHeadingSize(context),
            fontWeight: FontWeight.bold,
          ),
        ),
        actions: [
          if (_notifications.any((n) => !n['isRead']))
            TextButton(
              onPressed: _markAllAsRead,
              child: Text(
                'Mark all read',
                style: TextStyle(
                  color: AppTheme.authAccentColor,
                  fontSize: ResponsiveUtils.getSmallTextSize(context),
                  fontWeight: FontWeight.w600,
                ),
              ),
            ),
          IconButton(
            icon: Icon(
              Icons.refresh,
              color: AppTheme.authAccentColor,
              size: 24,
            ),
            onPressed: _loadNotifications,
          ),
        ],
      ),
      body: _buildBody(),
    );
  }

  Widget _buildBody() {
    if (_isLoading) {
      return Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            CircularProgressIndicator(
              valueColor: AlwaysStoppedAnimation<Color>(AppTheme.authAccentColor),
            ),
            const SizedBox(height: 16),
            Text(
              'Loading notifications...',
              style: TextStyle(
                color: Colors.grey[600],
                fontSize: ResponsiveUtils.getBodySize(context),
              ),
            ),
          ],
        ),
      );
    }

    if (_hasError) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.error_outline,
                size: 64,
                color: Colors.red[400],
              ),
              const SizedBox(height: 16),
              Text(
                'Failed to load notifications',
                style: AppTheme.headingTextStyle.copyWith(
                  color: Colors.red[700],
                  fontSize: ResponsiveUtils.getTitleSize(context),
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 8),
              Text(
                _errorMessage,
                style: TextStyle(
                  color: Colors.grey[600],
                  fontSize: ResponsiveUtils.getBodySize(context),
                ),
                textAlign: TextAlign.center,
              ),
              const SizedBox(height: 24),
              ElevatedButton.icon(
                onPressed: _loadNotifications,
                icon: const Icon(Icons.refresh, color: Colors.white),
                label: const Text(
                  'Retry',
                  style: TextStyle(color: Colors.white),
                ),
                style: ElevatedButton.styleFrom(
                  backgroundColor: AppTheme.authAccentColor,
                  padding: const EdgeInsets.symmetric(horizontal: 24, vertical: 12),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(25),
                  ),
                ),
              ),
            ],
          ),
        ),
      );
    }

    if (_notifications.isEmpty) {
      return Center(
        child: Padding(
          padding: const EdgeInsets.all(32),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                Icons.notifications_none,
                size: 64,
                color: Colors.grey[400],
              ),
              const SizedBox(height: 16),
              Text(
                'No notifications',
                style: AppTheme.headingTextStyle.copyWith(
                  color: Colors.grey[600],
                  fontSize: ResponsiveUtils.getTitleSize(context),
                ),
              ),
              const SizedBox(height: 8),
              Text(
                'You\'re all caught up! New notifications will appear here.',
                style: TextStyle(
                  color: Colors.grey[500],
                  fontSize: ResponsiveUtils.getBodySize(context),
                ),
                textAlign: TextAlign.center,
              ),
            ],
          ),
        ),
      );
    }

    return RefreshIndicator(
      onRefresh: _loadNotifications,
      color: AppTheme.authAccentColor,
      child: ListView.builder(
        padding: const EdgeInsets.all(16),
        itemCount: _notifications.length,
        itemBuilder: (context, index) {
          final notification = _notifications[index];
          return _buildNotificationCard(notification, index);
        },
      ),
    );
  }

  Widget _buildNotificationCard(Map<String, dynamic> notification, int index) {
    final bool isRead = notification['isRead'] ?? false;
    final String type = notification['type'] ?? 'info';
    final String priority = notification['priority'] ?? 'medium';
    
    return Container(
      margin: const EdgeInsets.only(bottom: 12),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isRead ? Colors.transparent : _getPriorityColor(priority).withOpacity(0.3),
          width: isRead ? 0 : 2,
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: InkWell(
        borderRadius: BorderRadius.circular(16),
        onTap: () {
          if (!isRead) {
            _markAsRead(notification['id'], index);
          }
          _showNotificationDetails(notification);
        },
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Row(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Icon
              Container(
                width: 48,
                height: 48,
                decoration: BoxDecoration(
                  color: _getTypeColor(type).withOpacity(0.1),
                  borderRadius: BorderRadius.circular(24),
                ),
                child: Icon(
                  _getTypeIcon(type),
                  color: _getTypeColor(type),
                  size: 24,
                ),
              ),
              const SizedBox(width: 16),
              
              // Content
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Title
                    Text(
                      notification['title'] ?? 'Notification',
                      style: TextStyle(
                        fontSize: ResponsiveUtils.getBodySize(context),
                        fontWeight: isRead ? FontWeight.w500 : FontWeight.w700,
                        color: isRead ? Colors.grey[700] : Colors.black87,
                      ),
                    ),
                    const SizedBox(height: 4),
                    
                    // Message
                    Text(
                      notification['message'] ?? '',
                      style: TextStyle(
                        fontSize: ResponsiveUtils.getSmallTextSize(context),
                        color: isRead ? Colors.grey[500] : Colors.grey[600],
                        height: 1.4,
                      ),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    const SizedBox(height: 8),
                    
                    // Timestamp and Priority
                    Row(
                      children: [
                        Text(
                          _formatTimestamp(notification['timestamp']),
                          style: TextStyle(
                            fontSize: ResponsiveUtils.getSmallTextSize(context) - 1,
                            color: Colors.grey[400],
                          ),
                        ),
                        const Spacer(),
                        if (priority == 'high')
                          Container(
                            padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
                            decoration: BoxDecoration(
                              color: Colors.red.withOpacity(0.1),
                              borderRadius: BorderRadius.circular(12),
                            ),
                            child: Text(
                              'High Priority',
                              style: TextStyle(
                                fontSize: ResponsiveUtils.getSmallTextSize(context) - 1,
                                color: Colors.red[700],
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                          ),
                      ],
                    ),
                  ],
                ),
              ),
              
              // Unread indicator
              if (!isRead)
                Container(
                  width: 8,
                  height: 8,
                  decoration: BoxDecoration(
                    color: _getPriorityColor(priority),
                    borderRadius: BorderRadius.circular(4),
                  ),
                ),
            ],
          ),
        ),
      ),
    );
  }

  void _showNotificationDetails(Map<String, dynamic> notification) {
    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      backgroundColor: Colors.transparent,
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.6,
        maxChildSize: 0.9,
        minChildSize: 0.3,
        builder: (context, scrollController) => Container(
          decoration: const BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
          ),
          child: Column(
            children: [
              // Handle
              Container(
                width: 40,
                height: 4,
                margin: const EdgeInsets.symmetric(vertical: 12),
                decoration: BoxDecoration(
                  color: Colors.grey[300],
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              
              // Content
              Expanded(
                child: SingleChildScrollView(
                  controller: scrollController,
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Title
                      Text(
                        notification['title'] ?? 'Notification',
                        style: AppTheme.headingTextStyle.copyWith(
                          fontSize: ResponsiveUtils.getTitleSize(context),
                          color: AppTheme.authDarkBlue,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      const SizedBox(height: 8),
                      
                      // Timestamp
                      Text(
                        _formatTimestamp(notification['timestamp']),
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                          color: Colors.grey[500],
                        ),
                      ),
                      const SizedBox(height: 16),
                      
                      // Message
                      Text(
                        notification['message'] ?? '',
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getBodySize(context),
                          color: Colors.black87,
                          height: 1.5,
                        ),
                      ),
                      
                      // Additional data based on type
                      if (notification['data'] != null)
                        _buildAdditionalData(notification['data'], notification['type']),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildAdditionalData(Map<String, dynamic> data, String type) {
    switch (type) {
      case 'suggestion':
        final suggestions = List<String>.from(data['suggestions'] ?? []);
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 24),
            Text(
              'AI Suggestions:',
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context),
                fontWeight: FontWeight.w600,
                color: AppTheme.authDarkBlue,
              ),
            ),
            const SizedBox(height: 12),
            ...suggestions.map((suggestion) => Container(
              margin: const EdgeInsets.only(bottom: 8),
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.amber.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: Colors.amber.withOpacity(0.3)),
              ),
              child: Row(
                children: [
                  Icon(Icons.lightbulb_outline, color: Colors.amber[700], size: 16),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      suggestion,
                      style: TextStyle(
                        fontSize: ResponsiveUtils.getSmallTextSize(context),
                        color: Colors.black87,
                      ),
                    ),
                  ),
                ],
              ),
            )).toList(),
          ],
        );
      
      case 'success':
        final positiveFactors = List<String>.from(data['positiveFactors'] ?? []);
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 24),
            Text(
              'Positive Factors:',
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context),
                fontWeight: FontWeight.w600,
                color: AppTheme.authDarkBlue,
              ),
            ),
            const SizedBox(height: 12),
            ...positiveFactors.map((factor) => Container(
              margin: const EdgeInsets.only(bottom: 8),
              padding: const EdgeInsets.all(12),
              decoration: BoxDecoration(
                color: Colors.green.withOpacity(0.1),
                borderRadius: BorderRadius.circular(8),
                border: Border.all(color: Colors.green.withOpacity(0.3)),
              ),
              child: Row(
                children: [
                  Icon(Icons.check_circle_outline, color: Colors.green[700], size: 16),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      factor,
                      style: TextStyle(
                        fontSize: ResponsiveUtils.getSmallTextSize(context),
                        color: Colors.black87,
                      ),
                    ),
                  ),
                ],
              ),
            )).toList(),
          ],
        );
      
      default:
        return const SizedBox.shrink();
    }
  }

  // Helper methods
  Color _getTypeColor(String type) {
    switch (type) {
      case 'alert': return Colors.red;
      case 'suggestion': return Colors.amber;
      case 'success': return Colors.green;
      case 'info': return Colors.blue;
      case 'risk_update': return AppTheme.authAccentColor;
      default: return Colors.grey;
    }
  }

  IconData _getTypeIcon(String type) {
    switch (type) {
      case 'alert': return Icons.warning;
      case 'suggestion': return Icons.lightbulb;
      case 'success': return Icons.check_circle;
      case 'info': return Icons.info;
      case 'risk_update': return Icons.analytics;
      default: return Icons.notifications;
    }
  }

  Color _getPriorityColor(String priority) {
    switch (priority) {
      case 'high': return Colors.red;
      case 'medium': return Colors.orange;
      case 'low': return Colors.blue;
      default: return Colors.grey;
    }
  }

  String _formatTimestamp(String? timestamp) {
    if (timestamp == null) return '';
    
    try {
      final dateTime = DateTime.parse(timestamp);
      final now = DateTime.now();
      final difference = now.difference(dateTime);
      
      if (difference.inMinutes < 1) {
        return 'Just now';
      } else if (difference.inMinutes < 60) {
        return '${difference.inMinutes}m ago';
      } else if (difference.inHours < 24) {
        return '${difference.inHours}h ago';
      } else if (difference.inDays < 7) {
        return '${difference.inDays}d ago';
      } else {
        return '${dateTime.day}/${dateTime.month}/${dateTime.year}';
      }
    } catch (e) {
      return '';
    }
  }
}

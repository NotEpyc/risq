import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';

class NotificationsListPage extends StatefulWidget {
  @override
  State<NotificationsListPage> createState() => _NotificationsListPageState();
}

class _NotificationsListPageState extends State<NotificationsListPage> {
  final List<NotificationItem> notifications = [
    NotificationItem(
      title: 'Risk Analysis Complete',
      message: 'Your startup risk profile has been updated with new insights.',
      time: '2 hours ago',
      isRead: false,
      icon: Icons.analytics,
      color: Colors.blue,
    ),
    NotificationItem(
      title: 'New Market Insight',
      message: 'AI education market shows 15% growth this quarter.',
      time: '1 day ago',
      isRead: true,
      icon: Icons.trending_up,
      color: Colors.green,
    ),
    NotificationItem(
      title: 'Speculation Reminder',
      message: 'Complete your pending speculation on product pricing.',
      time: '2 days ago',
      isRead: false,
      icon: Icons.psychology,
      color: Colors.orange,
    ),
    NotificationItem(
      title: 'Weekly Report Ready',
      message: 'Your weekly risk assessment report is now available.',
      time: '1 week ago',
      isRead: true,
      icon: Icons.description,
      color: Colors.purple,
    ),
    NotificationItem(
      title: 'Team Risk Update',
      message: 'Team risk score improved from 6.5 to 4.2.',
      time: '1 week ago',
      isRead: true,
      icon: Icons.people,
      color: Colors.green,
    ),
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        toolbarHeight: 60,
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.black87),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          'Notifications',
          style: TextStyle(
            color: Colors.black87,
            fontSize: ResponsiveUtils.getBodySize(context) + 2,
            fontWeight: FontWeight.w600,
          ),
        ),
        actions: [
          IconButton(
            icon: Icon(Icons.mark_email_read, color: AppTheme.authPrimaryColor),
            onPressed: () {
              setState(() {
                for (var notification in notifications) {
                  notification.isRead = true;
                }
              });
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('All notifications marked as read')),
              );
            },
          ),
        ],
      ),
      body: notifications.isEmpty 
          ? _buildEmptyState()
          : ListView.builder(
              padding: EdgeInsets.all(20),
              itemCount: notifications.length,
              itemBuilder: (context, index) {
                return _buildNotificationCard(notifications[index], index);
              },
            ),
    );
  }

  Widget _buildEmptyState() {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            Icons.notifications_none,
            size: 80,
            color: Colors.grey[400],
          ),
          SizedBox(height: 16),
          Text(
            'No notifications yet',
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context) + 2,
              fontWeight: FontWeight.w600,
              color: Colors.grey[600],
            ),
          ),
          SizedBox(height: 8),
          Text(
            'You\'ll see notifications here when there are updates',
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context),
              color: Colors.grey[500],
            ),
            textAlign: TextAlign.center,
          ),
        ],
      ),
    );
  }

  Widget _buildNotificationCard(NotificationItem notification, int index) {
    return Dismissible(
      key: Key('notification_$index'),
      direction: DismissDirection.endToStart,
      background: Container(
        margin: EdgeInsets.only(bottom: 12),
        padding: EdgeInsets.symmetric(horizontal: 20),
        decoration: BoxDecoration(
          color: Colors.red,
          borderRadius: BorderRadius.circular(16),
        ),
        alignment: Alignment.centerRight,
        child: Icon(Icons.delete, color: Colors.white),
      ),
      onDismissed: (direction) {
        setState(() {
          notifications.removeAt(index);
        });
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('Notification deleted'),
            action: SnackBarAction(
              label: 'Undo',
              onPressed: () {
                setState(() {
                  notifications.insert(index, notification);
                });
              },
            ),
          ),
        );
      },
      child: Container(
        margin: EdgeInsets.only(bottom: 12),
        child: Material(
          color: notification.isRead ? Colors.grey[50] : Colors.white,
          borderRadius: BorderRadius.circular(16),
          child: InkWell(
            onTap: () {
              setState(() {
                notification.isRead = true;
              });
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Opening: ${notification.title}')),
              );
            },
            borderRadius: BorderRadius.circular(16),
            child: Container(
              padding: EdgeInsets.all(16),
              decoration: BoxDecoration(
                borderRadius: BorderRadius.circular(16),
                border: notification.isRead 
                    ? null 
                    : Border.all(color: AppTheme.authPrimaryColor.withOpacity(0.3)),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.05),
                    blurRadius: 10,
                    offset: Offset(0, 5),
                  ),
                ],
              ),
              child: Row(
                children: [
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: notification.color.withOpacity(0.1),
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Icon(
                      notification.icon,
                      color: notification.color,
                      size: 20,
                    ),
                  ),
                  SizedBox(width: 16),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Row(
                          children: [
                            Expanded(
                              child: Text(
                                notification.title,
                                style: TextStyle(
                                  fontSize: ResponsiveUtils.getBodySize(context),
                                  fontWeight: notification.isRead ? FontWeight.w500 : FontWeight.w600,
                                  color: Colors.black87,
                                ),
                              ),
                            ),
                            if (!notification.isRead)
                              Container(
                                width: 8,
                                height: 8,
                                decoration: BoxDecoration(
                                  color: AppTheme.authPrimaryColor,
                                  borderRadius: BorderRadius.circular(4),
                                ),
                              ),
                          ],
                        ),
                        SizedBox(height: 4),
                        Text(
                          notification.message,
                          style: TextStyle(
                            fontSize: ResponsiveUtils.getSmallTextSize(context),
                            color: Colors.grey[600],
                            height: 1.3,
                          ),
                          maxLines: 2,
                          overflow: TextOverflow.ellipsis,
                        ),
                        SizedBox(height: 4),
                        Text(
                          notification.time,
                          style: TextStyle(
                            fontSize: ResponsiveUtils.getSmallTextSize(context) - 1,
                            color: Colors.grey[500],
                          ),
                        ),
                      ],
                    ),
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }
}

class NotificationItem {
  final String title;
  final String message;
  final String time;
  bool isRead;
  final IconData icon;
  final Color color;

  NotificationItem({
    required this.title,
    required this.message,
    required this.time,
    required this.isRead,
    required this.icon,
    required this.color,
  });
}

import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';

class ProfilePage extends StatefulWidget {
  final String userName;
  final String userEmail;

  const ProfilePage({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
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
          'Profile',
          style: TextStyle(
            color: Colors.black87,
            fontSize: ResponsiveUtils.getBodySize(context) + 2,
            fontWeight: FontWeight.w600,
          ),
        ),
        actions: [
          IconButton(
            icon: Icon(Icons.edit, color: AppTheme.authPrimaryColor),
            onPressed: () {
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(content: Text('Edit profile coming soon!')),
              );
            },
          ),
        ],
      ),
      body: SingleChildScrollView(
        padding: EdgeInsets.all(20),
        child: Column(
          children: [
            // Profile Header
            _buildProfileHeader(),
            SizedBox(height: 30),
            
            // Profile Stats
            _buildProfileStats(),
            SizedBox(height: 30),
            
            // Profile Options
            _buildProfileOptions(),
          ],
        ),
      ),
    );
  }

  Widget _buildProfileHeader() {
    return Container(
      padding: EdgeInsets.all(24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(20),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.08),
            blurRadius: 20,
            offset: Offset(0, 10),
          ),
        ],
      ),
      child: Column(
        children: [
          // Profile Picture
          Container(
            width: 80,
            height: 80,
            decoration: BoxDecoration(
              color: AppTheme.authAccentColor,
              borderRadius: BorderRadius.circular(40),
              border: Border.all(color: Colors.white, width: 4),
            ),
            child: Center(
              child: Text(
                widget.userName.isNotEmpty 
                    ? widget.userName[0].toUpperCase() 
                    : 'U',
                style: TextStyle(
                  color: Colors.white,
                  fontSize: 32,
                  fontWeight: FontWeight.bold,
                ),
              ),
            ),
          ),
          SizedBox(height: 16),
          
          // Name and Email
          Text(
            widget.userName.isNotEmpty ? widget.userName : 'User',
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context) + 4,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ),
          SizedBox(height: 4),
          Text(
            widget.userEmail,
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context),
              color: Colors.grey[600],
            ),
          ),
          SizedBox(height: 16),
          
          // Role Badge
          Container(
            padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: AppTheme.authPrimaryColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Text(
              'Startup Founder',
              style: TextStyle(
                color: AppTheme.authPrimaryColor,
                fontSize: ResponsiveUtils.getSmallTextSize(context),
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildProfileStats() {
    return Row(
      children: [
        Expanded(
          child: _buildStatCard('Risk Score', '7.2', Icons.warning_amber, Colors.orange),
        ),
        SizedBox(width: 16),
        Expanded(
          child: _buildStatCard('Speculations', '12', Icons.psychology, Colors.blue),
        ),
        SizedBox(width: 16),
        Expanded(
          child: _buildStatCard('Decisions', '8', Icons.fact_check, Colors.green),
        ),
      ],
    );
  }

  Widget _buildStatCard(String title, String value, IconData icon, Color color) {
    return Container(
      height: 130, // Increased height to prevent overflow
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: Offset(0, 5),
          ),
        ],
      ),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Container(
            width: 36,
            height: 36,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(18),
            ),
            child: Icon(icon, color: color, size: 18),
          ),
          SizedBox(height: 6),
          Text(
            value,
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context) + 1,
              fontWeight: FontWeight.bold,
              color: Colors.black87,
            ),
          ),
          SizedBox(height: 2),
          Text(
            title,
            style: TextStyle(
              fontSize: ResponsiveUtils.getSmallTextSize(context) - 1,
              color: Colors.grey[600],
            ),
            textAlign: TextAlign.center,
            maxLines: 1,
            overflow: TextOverflow.ellipsis,
          ),
        ],
      ),
    );
  }

  Widget _buildProfileOptions() {
    return Column(
      children: [
        _buildOptionCard(
          'Startup Information',
          'Update your company details',
          Icons.business,
          () {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Startup info coming soon!')),
            );
          },
        ),
        _buildOptionCard(
          'Risk Preferences',
          'Customize your risk assessment',
          Icons.tune,
          () {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Risk preferences coming soon!')),
            );
          },
        ),
        _buildOptionCard(
          'Notifications',
          'Manage your notification settings',
          Icons.notifications,
          () {
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => NotificationsPage(),
              ),
            );
          },
        ),
        _buildOptionCard(
          'Export Data',
          'Download your data and reports',
          Icons.download,
          () {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Export feature coming soon!')),
            );
          },
        ),
        _buildOptionCard(
          'Help & Support',
          'Get help and contact support',
          Icons.help,
          () {
            ScaffoldMessenger.of(context).showSnackBar(
              SnackBar(content: Text('Support coming soon!')),
            );
          },
        ),
        _buildOptionCard(
          'Sign Out',
          'Sign out of your account',
          Icons.logout,
          () {
            _showSignOutDialog();
          },
          isDestructive: true,
        ),
      ],
    );
  }

  Widget _buildOptionCard(String title, String subtitle, IconData icon, VoidCallback onTap, {bool isDestructive = false}) {
    return Container(
      margin: EdgeInsets.only(bottom: 12),
      child: Material(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(16),
          child: Container(
            padding: EdgeInsets.all(16),
            decoration: BoxDecoration(
              borderRadius: BorderRadius.circular(16),
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
                    color: isDestructive 
                        ? Colors.red.withOpacity(0.1)
                        : AppTheme.authPrimaryColor.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Icon(
                    icon, 
                    color: isDestructive ? Colors.red : AppTheme.authPrimaryColor, 
                    size: 20
                  ),
                ),
                SizedBox(width: 16),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        title,
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getBodySize(context),
                          fontWeight: FontWeight.w600,
                          color: isDestructive ? Colors.red : Colors.black87,
                        ),
                      ),
                      SizedBox(height: 2),
                      Text(
                        subtitle,
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                  ),
                ),
                Icon(
                  Icons.arrow_forward_ios,
                  size: 16,
                  color: Colors.grey[400],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  void _showSignOutDialog() {
    showDialog(
      context: context,
      builder: (BuildContext context) {
        return AlertDialog(
          title: Text('Sign Out'),
          content: Text('Are you sure you want to sign out?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context),
              child: Text('Cancel'),
            ),
            TextButton(
              onPressed: () {
                Navigator.pop(context);
                ScaffoldMessenger.of(context).showSnackBar(
                  SnackBar(content: Text('Sign out functionality coming soon!')),
                );
              },
              child: Text('Sign Out', style: TextStyle(color: Colors.red)),
            ),
          ],
        );
      },
    );
  }
}

class NotificationsPage extends StatefulWidget {
  @override
  State<NotificationsPage> createState() => _NotificationsPageState();
}

class _NotificationsPageState extends State<NotificationsPage> {
  bool pushNotifications = true;
  bool emailNotifications = true;
  bool riskAlerts = true;
  bool weeklyReports = false;

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      appBar: AppBar(
        backgroundColor: Colors.white,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.black87),
          onPressed: () => Navigator.pop(context),
        ),
        title: Text(
          'Notifications',
          style: TextStyle(
            color: Colors.black87,
            fontSize: ResponsiveUtils.getBodySize(context) + 4,
            fontWeight: FontWeight.w600,
          ),
        ),
      ),
      body: ListView(
        padding: EdgeInsets.all(20),
        children: [
          _buildNotificationSection(
            'Push Notifications',
            'Get notified on your device',
            pushNotifications,
            (value) => setState(() => pushNotifications = value),
          ),
          _buildNotificationSection(
            'Email Notifications',
            'Receive updates via email',
            emailNotifications,
            (value) => setState(() => emailNotifications = value),
          ),
          _buildNotificationSection(
            'Risk Alerts',
            'Get notified of important risk changes',
            riskAlerts,
            (value) => setState(() => riskAlerts = value),
          ),
          _buildNotificationSection(
            'Weekly Reports',
            'Receive weekly summary reports',
            weeklyReports,
            (value) => setState(() => weeklyReports = value),
          ),
        ],
      ),
    );
  }

  Widget _buildNotificationSection(String title, String subtitle, bool value, Function(bool) onChanged) {
    return Container(
      margin: EdgeInsets.only(bottom: 12),
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
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
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  title,
                  style: TextStyle(
                    fontSize: ResponsiveUtils.getBodySize(context),
                    fontWeight: FontWeight.w600,
                    color: Colors.black87,
                  ),
                ),
                SizedBox(height: 2),
                Text(
                  subtitle,
                  style: TextStyle(
                    fontSize: ResponsiveUtils.getSmallTextSize(context),
                    color: Colors.grey[600],
                  ),
                ),
              ],
            ),
          ),
          Switch(
            value: value,
            onChanged: onChanged,
            activeColor: AppTheme.authPrimaryColor,
          ),
        ],
      ),
    );
  }
}

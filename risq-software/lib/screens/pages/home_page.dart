import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/screens/pages/startup_profile_page.dart';
import 'package:risq/screens/pages/notifications_page.dart';
import 'package:syncfusion_flutter_gauges/gauges.dart';

class HomePage extends StatefulWidget {
  final String userName;
  final String userEmail;

  const HomePage({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      body: SafeArea(
        child: Column(
          children: [
            // Top Header with Profile and Notifications
            Container(
              padding: EdgeInsets.all(20),
              decoration: BoxDecoration(
                color: Colors.white,
              ),
              child: Row(
                children: [
                  // Profile Picture
                  GestureDetector(
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => const StartupProfilePage(),
                        ),
                      );
                    },
                    child: Container(
                      width: 50,
                      height: 50,
                      decoration: BoxDecoration(
                        color: AppTheme.authAccentColor,
                        borderRadius: BorderRadius.circular(25),
                        border: Border.all(color: Colors.white, width: 2),
                      ),
                      child: Center(
                        child: Text(
                          widget.userName.isNotEmpty 
                              ? widget.userName[0].toUpperCase() 
                              : 'U',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: ResponsiveUtils.getBodySize(context) + 2,
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  ),
                  
                  Spacer(),
                  
                  // Notification Button
                  GestureDetector(
                    onTap: () {
                      Navigator.push(
                        context,
                        MaterialPageRoute(
                          builder: (context) => const NotificationsPage(),
                        ),
                      );
                    },
                    child: Container(
                      width: 50,
                      height: 50,
                      decoration: BoxDecoration(
                        color: Colors.grey[100],
                        borderRadius: BorderRadius.circular(25),
                      ),
                      child: Stack(
                        children: [
                          Center(
                            child: Icon(
                              Icons.notifications_outlined,
                              color: Colors.grey[700],
                              size: 24,
                            ),
                          ),
                          // Notification badge
                          Positioned(
                            right: 12,
                            top: 12,
                            child: Container(
                              width: 8,
                              height: 8,
                              decoration: BoxDecoration(
                                color: Colors.red,
                                borderRadius: BorderRadius.circular(4),
                              ),
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
                ],
              ),
            ),
            
            // Main Content
            Expanded(
              child: SingleChildScrollView(
                padding: EdgeInsets.all(20),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    // Risk Score Meter
                    _buildRiskScoreMeter(),
                    SizedBox(height: 30),
                    
                    // Quick Actions
                    _buildQuickActionsSection(),
                    SizedBox(height: 20),
                    
                    // Recent Activity
                    _buildRecentActivitySection(),
                    SizedBox(height: 20),
                    
                    // AI Suggestions
                    _buildAISuggestionsSection(),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRiskScoreMeter() {
    return Container(
      width: double.infinity,
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
          Text(
            'Risk Score',
            style: TextStyle(
              color: Colors.grey[700],
              fontSize: ResponsiveUtils.getBodySize(context) + 2,
              fontWeight: FontWeight.w600,
            ),
          ),
          SizedBox(height: 20),
          Container(
            height: 200,
            child: SfRadialGauge(
              axes: <RadialAxis>[
                RadialAxis(
                  minimum: 0,
                  maximum: 10,
                  showLabels: true,
                  showTicks: true,
                  radiusFactor: 0.9,
                  axisLineStyle: AxisLineStyle(
                    thickness: 0.15,
                    cornerStyle: CornerStyle.bothCurve,
                    color: Colors.grey[300],
                    thicknessUnit: GaugeSizeUnit.factor,
                  ),
                  majorTickStyle: MajorTickStyle(
                    length: 12,
                    thickness: 2,
                    color: Colors.grey[400],
                  ),
                  minorTickStyle: MinorTickStyle(
                    length: 6,
                    thickness: 1,
                    color: Colors.grey[300],
                  ),
                  axisLabelStyle: GaugeTextStyle(
                    color: Colors.grey[600],
                    fontSize: 12,
                    fontWeight: FontWeight.w500,
                  ),
                  ranges: <GaugeRange>[
                    GaugeRange(
                      startValue: 0,
                      endValue: 7.2, // Only fill up to the current risk score
                      color: AppTheme.authAccentColor,
                      startWidth: 12,
                      endWidth: 12,
                    ),
                  ],
                  annotations: <GaugeAnnotation>[
                    GaugeAnnotation(
                      widget: Column(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Text(
                            '7.2',
                            style: TextStyle(
                              fontSize: 28,
                              fontWeight: FontWeight.bold,
                              color: AppTheme.authAccentColor,
                            ),
                          ),
                          Text(
                            'Medium Risk',
                            style: TextStyle(
                              fontSize: 14,
                              fontWeight: FontWeight.w500,
                              color: Colors.grey[600],
                            ),
                          ),
                        ],
                      ),
                      angle: 90,
                      positionFactor: 0.5,
                    ),
                  ],
                ),
              ],
            ),
          ),
          SizedBox(height: 20),
          Container(
            padding: EdgeInsets.symmetric(horizontal: 16, vertical: 8),
            decoration: BoxDecoration(
              color: AppTheme.authAccentColor.withOpacity(0.1),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Text(
              '85% Confidence Level',
              style: TextStyle(
                color: AppTheme.authAccentColor,
                fontSize: ResponsiveUtils.getSmallTextSize(context),
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildQuickActionsSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Quick Actions',
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context) + 2,
            color: Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
        SizedBox(height: 16),
        Row(
          children: [
            Expanded(
              child: _buildActionCard(
                icon: Icons.psychology,
                title: 'New Speculation',
                subtitle: 'Test a decision',
                color: Colors.blue,
                onTap: () {
                  // TODO: Navigate to speculation page
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Use the Speculation tab below!')),
                  );
                },
              ),
            ),
            SizedBox(width: 16),
            Expanded(
              child: _buildActionCard(
                icon: Icons.analytics,
                title: 'View Report',
                subtitle: 'Detailed analysis',
                color: Colors.green,
                onTap: () {
                  // TODO: Navigate to reports
                  ScaffoldMessenger.of(context).showSnackBar(
                    SnackBar(content: Text('Reports coming soon!')),
                  );
                },
              ),
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildActionCard({
    required IconData icon,
    required String title,
    required String subtitle,
    required Color color,
    required VoidCallback onTap,
  }) {
    return GestureDetector(
      onTap: onTap,
      child: Container(
        padding: EdgeInsets.all(20),
        decoration: BoxDecoration(
          color: Colors.white,
          borderRadius: BorderRadius.circular(16),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.05),
              blurRadius: 8,
              offset: Offset(0, 2),
            ),
          ],
        ),
        child: Column(
          children: [
            Container(
              width: 50,
              height: 50,
              decoration: BoxDecoration(
                color: color.withOpacity(0.1),
                borderRadius: BorderRadius.circular(25),
              ),
              child: Icon(icon, color: color, size: 24),
            ),
            SizedBox(height: 12),
            Text(
              title,
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context),
                fontWeight: FontWeight.w600,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: 4),
            Text(
              subtitle,
              style: TextStyle(
                fontSize: ResponsiveUtils.getSmallTextSize(context),
                color: Colors.grey[600],
              ),
              textAlign: TextAlign.center,
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildRecentActivitySection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'Recent Activity',
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context) + 2,
            color: Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
        SizedBox(height: 16),
        _buildActivityItem(
          icon: Icons.add_circle_outline,
          title: 'Profile Created',
          subtitle: 'Your startup profile was successfully created',
          time: 'Today',
          color: Colors.green,
        ),
        _buildActivityItem(
          icon: Icons.analytics_outlined,
          title: 'Risk Analysis Complete',
          subtitle: 'AI analysis generated your risk profile',
          time: 'Today',
          color: Colors.blue,
        ),
      ],
    );
  }

  Widget _buildActivityItem({
    required IconData icon,
    required String title,
    required String subtitle,
    required String time,
    required Color color,
  }) {
    return Container(
      margin: EdgeInsets.only(bottom: 12),
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.03),
            blurRadius: 6,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Container(
            width: 40,
            height: 40,
            decoration: BoxDecoration(
              color: color.withOpacity(0.1),
              borderRadius: BorderRadius.circular(20),
            ),
            child: Icon(icon, color: color, size: 20),
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
                    color: Colors.black87,
                  ),
                ),
                SizedBox(height: 4),
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
          Text(
            time,
            style: TextStyle(
              fontSize: ResponsiveUtils.getSmallTextSize(context),
              color: Colors.grey[500],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildAISuggestionsSection() {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          'AI Suggestions',
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context) + 2,
            color: Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
        SizedBox(height: 16),
        _buildSuggestionCard(
          'Focus on unique AI differentiation in marketing',
          Icons.lightbulb_outline,
        ),
        _buildSuggestionCard(
          'Develop partnerships with educational institutions early',
          Icons.handshake_outlined,
        ),
        _buildSuggestionCard(
          'Consider pilot programs to validate product-market fit',
          Icons.science_outlined,
        ),
      ],
    );
  }

  Widget _buildSuggestionCard(String suggestion, IconData icon) {
    return Container(
      margin: EdgeInsets.only(bottom: 12),
      padding: EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
        border: Border.all(color: Colors.amber.withOpacity(0.3)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.03),
            blurRadius: 6,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Row(
        children: [
          Icon(icon, color: Colors.amber[700], size: 24),
          SizedBox(width: 16),
          Expanded(
            child: Text(
              suggestion,
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context),
                color: Colors.black87,
                height: 1.4,
              ),
            ),
          ),
        ],
      ),
    );
  }
}

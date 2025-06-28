import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/screens/pages/startup_profile_page.dart';

class DecisionPage extends StatefulWidget {
  final String userName;
  final String userEmail;

  const DecisionPage({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<DecisionPage> createState() => _DecisionPageState();
}

class _DecisionPageState extends State<DecisionPage> with TickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 2, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  // Mock data for decisions
  final List<Map<String, dynamic>> _confirmedDecisions = [
    {
      'id': '1',
      'description': 'Hire 3 senior engineers to accelerate product development',
      'category': 'hiring',
      'riskChange': 0.4,
      'confidence': 0.82,
      'date': 'Today',
      'status': 'confirmed',
    },
    {
      'id': '2',
      'description': 'Partner with educational institutions for pilot programs',
      'category': 'strategy',
      'riskChange': -0.2,
      'confidence': 0.78,
      'date': 'Yesterday',
      'status': 'confirmed',
    },
  ];

  final List<Map<String, dynamic>> _speculativeDecisions = [
    {
      'id': '3',
      'description': 'Launch paid marketing campaign on Google Ads',
      'category': 'marketing',
      'riskChange': -0.3,
      'confidence': 0.75,
      'date': '2 days ago',
      'status': 'speculative',
    },
    {
      'id': '4',
      'description': 'Raise Series A funding round',
      'category': 'funding',
      'riskChange': 0.6,
      'confidence': 0.85,
      'date': '3 days ago',
      'status': 'speculative',
    },
  ];

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.grey[50],
      body: SafeArea(
        child: Column(
          children: [
            // Top Header
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
                      width: 40,
                      height: 40,
                      decoration: BoxDecoration(
                        color: AppTheme.authAccentColor,
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: Center(
                        child: Text(
                          widget.userName.isNotEmpty 
                              ? widget.userName[0].toUpperCase() 
                              : 'U',
                          style: TextStyle(
                            color: Colors.white,
                            fontSize: ResponsiveUtils.getBodySize(context),
                            fontWeight: FontWeight.bold,
                          ),
                        ),
                      ),
                    ),
                  ),
                  SizedBox(width: 12),
                  Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Decision History',
                        style: AppTheme.headingTextStyle.copyWith(
                          fontSize: ResponsiveUtils.getBodySize(context) + 4,
                          color: Colors.black87,
                        ),
                      ),
                      Text(
                        'Track your business decisions',
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                          color: Colors.grey[600],
                        ),
                      ),
                    ],
                  ),
                  Spacer(),
                  // Notification Button
                  Container(
                    width: 40,
                    height: 40,
                    decoration: BoxDecoration(
                      color: Colors.grey[100],
                      borderRadius: BorderRadius.circular(20),
                    ),
                    child: Icon(
                      Icons.notifications_outlined,
                      color: Colors.grey[700],
                      size: 20,
                    ),
                  ),
                ],
              ),
            ),
            
            // Tab Bar
            Container(
              color: Colors.white,
              child: TabBar(
                controller: _tabController,
                labelColor: AppTheme.authAccentColor,
                unselectedLabelColor: Colors.grey[600],
                indicatorColor: AppTheme.authAccentColor,
                labelStyle: TextStyle(
                  fontSize: ResponsiveUtils.getBodySize(context),
                  fontWeight: FontWeight.w600,
                ),
                tabs: [
                  Tab(text: 'Confirmed (${_confirmedDecisions.length})'),
                  Tab(text: 'Speculative (${_speculativeDecisions.length})'),
                ],
              ),
            ),
            
            // Tab Content
            Expanded(
              child: TabBarView(
                controller: _tabController,
                children: [
                  _buildDecisionList(_confirmedDecisions, isConfirmed: true),
                  _buildDecisionList(_speculativeDecisions, isConfirmed: false),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDecisionList(List<Map<String, dynamic>> decisions, {required bool isConfirmed}) {
    if (decisions.isEmpty) {
      return _buildEmptyState(isConfirmed);
    }

    return ListView.builder(
      padding: EdgeInsets.all(20),
      itemCount: decisions.length,
      itemBuilder: (context, index) {
        final decision = decisions[index];
        return _buildDecisionCard(decision, isConfirmed);
      },
    );
  }

  Widget _buildEmptyState(bool isConfirmed) {
    return Center(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            isConfirmed ? Icons.check_circle_outline : Icons.psychology_outlined,
            size: 64,
            color: Colors.grey[400],
          ),
          SizedBox(height: 16),
          Text(
            isConfirmed ? 'No confirmed decisions yet' : 'No speculative decisions yet',
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context) + 2,
              fontWeight: FontWeight.w600,
              color: Colors.grey[600],
            ),
          ),
          SizedBox(height: 8),
          Text(
            isConfirmed 
                ? 'Decisions you confirm will appear here'
                : 'Test decisions in the Speculation tab',
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

  Widget _buildDecisionCard(Map<String, dynamic> decision, bool isConfirmed) {
    final riskChange = decision['riskChange'] as double;
    final isPositiveChange = riskChange > 0;
    
    return Container(
      margin: EdgeInsets.only(bottom: 16),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        border: Border.all(
          color: isConfirmed 
              ? Colors.green.withOpacity(0.3) 
              : Colors.orange.withOpacity(0.3),
        ),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 8,
            offset: Offset(0, 2),
          ),
        ],
      ),
      child: Padding(
        padding: EdgeInsets.all(20),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Header with status badge
            Row(
              children: [
                Container(
                  padding: EdgeInsets.symmetric(horizontal: 12, vertical: 6),
                  decoration: BoxDecoration(
                    color: isConfirmed 
                        ? Colors.green.withOpacity(0.1) 
                        : Colors.orange.withOpacity(0.1),
                    borderRadius: BorderRadius.circular(20),
                  ),
                  child: Row(
                    mainAxisSize: MainAxisSize.min,
                    children: [
                      Icon(
                        isConfirmed ? Icons.check_circle : Icons.help_outline,
                        size: 16,
                        color: isConfirmed ? Colors.green[700] : Colors.orange[700],
                      ),
                      SizedBox(width: 4),
                      Text(
                        isConfirmed ? 'Confirmed' : 'Speculative',
                        style: TextStyle(
                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                          fontWeight: FontWeight.w600,
                          color: isConfirmed ? Colors.green[700] : Colors.orange[700],
                        ),
                      ),
                    ],
                  ),
                ),
                Spacer(),
                _buildCategoryChip(decision['category']),
              ],
            ),
            SizedBox(height: 16),
            
            // Decision description
            Text(
              decision['description'],
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context) + 1,
                fontWeight: FontWeight.w600,
                color: Colors.black87,
                height: 1.4,
              ),
            ),
            SizedBox(height: 16),
            
            // Risk metrics
            Row(
              children: [
                Expanded(
                  child: _buildMetricItem(
                    'Risk Change',
                    '${isPositiveChange ? '+' : ''}${riskChange.toStringAsFixed(1)}',
                    isPositiveChange ? Colors.red : Colors.green,
                  ),
                ),
                Expanded(
                  child: _buildMetricItem(
                    'Confidence',
                    '${(decision['confidence'] * 100).toInt()}%',
                    Colors.blue,
                  ),
                ),
                Expanded(
                  child: _buildMetricItem(
                    'Date',
                    decision['date'],
                    Colors.grey[600]!,
                  ),
                ),
              ],
            ),
            
            if (!isConfirmed) ...[
              SizedBox(height: 16),
              Row(
                children: [
                  Expanded(
                    child: OutlinedButton(
                      onPressed: () {
                        _showDecisionDetails(decision);
                      },
                      style: OutlinedButton.styleFrom(
                        side: BorderSide(color: AppTheme.authAccentColor),
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: Text(
                        'View Details',
                        style: TextStyle(color: AppTheme.authAccentColor),
                      ),
                    ),
                  ),
                  SizedBox(width: 12),
                  Expanded(
                    child: ElevatedButton(
                      onPressed: () {
                        _confirmDecision(decision);
                      },
                      style: ElevatedButton.styleFrom(
                        backgroundColor: AppTheme.authAccentColor,
                        shape: RoundedRectangleBorder(
                          borderRadius: BorderRadius.circular(8),
                        ),
                      ),
                      child: Text(
                        'Confirm',
                        style: TextStyle(color: Colors.white),
                      ),
                    ),
                  ),
                ],
              ),
            ],
          ],
        ),
      ),
    );
  }

  Widget _buildCategoryChip(String category) {
    final categoryData = _getCategoryData(category);
    return Container(
      padding: EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: categoryData['color'].withOpacity(0.1),
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(
            categoryData['icon'],
            size: 14,
            color: categoryData['color'],
          ),
          SizedBox(width: 4),
          Text(
            categoryData['label'],
            style: TextStyle(
              fontSize: ResponsiveUtils.getSmallTextSize(context) - 1,
              fontWeight: FontWeight.w500,
              color: categoryData['color'],
            ),
          ),
        ],
      ),
    );
  }

  Map<String, dynamic> _getCategoryData(String category) {
    switch (category) {
      case 'hiring':
        return {'label': 'Hiring', 'icon': Icons.people, 'color': Colors.blue};
      case 'marketing':
        return {'label': 'Marketing', 'icon': Icons.campaign, 'color': Colors.green};
      case 'strategy':
        return {'label': 'Strategy', 'icon': Icons.business, 'color': Colors.purple};
      case 'funding':
        return {'label': 'Funding', 'icon': Icons.attach_money, 'color': Colors.orange};
      default:
        return {'label': 'Other', 'icon': Icons.help_outline, 'color': Colors.grey};
    }
  }

  Widget _buildMetricItem(String label, String value, Color color) {
    return Column(
      children: [
        Text(
          label,
          style: TextStyle(
            fontSize: ResponsiveUtils.getSmallTextSize(context),
            color: Colors.grey[600],
          ),
        ),
        SizedBox(height: 4),
        Text(
          value,
          style: TextStyle(
            fontSize: ResponsiveUtils.getBodySize(context),
            fontWeight: FontWeight.w600,
            color: color,
          ),
        ),
      ],
    );
  }

  void _showDecisionDetails(Map<String, dynamic> decision) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: Text('Decision Details'),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              decision['description'],
              style: TextStyle(
                fontSize: ResponsiveUtils.getBodySize(context),
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            Text(
              'AI Reasoning:',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            SizedBox(height: 8),
            Text(
              'This decision analysis is based on current market conditions, team capacity, and financial projections. The AI considers various risk factors to provide this assessment.',
              style: TextStyle(color: Colors.grey[700]),
            ),
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Close'),
          ),
        ],
      ),
    );
  }

  void _confirmDecision(Map<String, dynamic> decision) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: Text('Confirm Decision'),
        content: Text(
          'Are you sure you want to confirm this decision? This will update your risk profile and move it to confirmed decisions.',
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Cancel'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
              setState(() {
                _speculativeDecisions.remove(decision);
                _confirmedDecisions.insert(0, {
                  ...decision,
                  'status': 'confirmed',
                  'date': 'Today',
                });
              });
              ScaffoldMessenger.of(context).showSnackBar(
                SnackBar(
                  content: Text('Decision confirmed! Risk profile updated.'),
                  backgroundColor: AppTheme.authAccentColor,
                ),
              );
            },
            style: ElevatedButton.styleFrom(
              backgroundColor: AppTheme.authAccentColor,
            ),
            child: Text('Confirm', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
  }
}

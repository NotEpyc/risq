import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/screens/main_navigation.dart';

class DataDisplayPage extends StatelessWidget {
  final Map<String, dynamic> submittedData;

  const DataDisplayPage({
    super.key,
    required this.submittedData,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(
          'Submitted Data',
          style: AppTheme.headingTextStyle.copyWith(
            color: Colors.white,
            fontSize: ResponsiveUtils.getHeadingSize(context) * 0.8,
          ),
        ),
        backgroundColor: AppTheme.authAccentColor,
        elevation: 0,
        leading: IconButton(
          icon: Icon(Icons.arrow_back, color: Colors.white),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      backgroundColor: Colors.grey[50],
      body: SingleChildScrollView(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            // Success message
            Container(
              width: double.infinity,
              padding: EdgeInsets.all(16),
              margin: EdgeInsets.only(bottom: 24),
              decoration: BoxDecoration(
                color: Colors.green[50],
                border: Border.all(color: Colors.green[200]!),
                borderRadius: BorderRadius.circular(12),
              ),
              child: Row(
                children: [
                  Icon(Icons.check_circle, color: Colors.green[600], size: 24),
                  SizedBox(width: 12),
                  Expanded(
                    child: Text(
                      'Your startup profile has been successfully submitted!',
                      style: TextStyle(
                        color: Colors.green[700],
                        fontSize: ResponsiveUtils.getBodySize(context),
                        fontWeight: FontWeight.w500,
                      ),
                    ),
                  ),
                ],
              ),
            ),

            // Company Information
            _buildSection(
              context,
              'Company Information',
              [
                _buildDataRow('Company Name', submittedData['name']),
                _buildDataRow('Description', submittedData['description']),
                _buildDataRow('Industry', submittedData['industry']),
                _buildDataRow('Sector', submittedData['sector']),
                _buildDataRow('Location', submittedData['location']),
                _buildDataRow('Website', submittedData['website']),
                _buildDataRow('Founded Date', submittedData['founded_date']),
              ],
            ),

            // Business Model & Market
            _buildSection(
              context,
              'Business Model & Market',
              [
                _buildDataRow('Business Model', submittedData['business_model']),
                _buildDataRow('Target Market', submittedData['target_market']),
                _buildDataRow('Competitive Advantage', submittedData['competitive_advantage']),
                _buildDataRow('Go-to-Market Strategy', submittedData['go_to_market_strategy']),
                _buildDataRow('Funding Stage', submittedData['funding_stage']),
                _buildDataRow('Team Size', submittedData['team_size']?.toString()),
                _buildArrayDataRow('Revenue Streams', submittedData['revenue_streams']),
              ],
            ),

            // Technical & Implementation
            _buildSection(
              context,
              'Technical & Implementation',
              [
                _buildDataRow('Implementation Plan', submittedData['implementation_plan']),
                _buildDataRow('Development Timeline', submittedData['development_timeline']),
                _buildArrayDataRow('Technology Stack', submittedData['technology_stack']),
              ],
            ),

            // Financial Information
            _buildSection(
              context,
              'Financial Information',
              [
                _buildDataRow('Initial Investment', '\$${submittedData['initial_investment']?.toString() ?? '0'}'),
                _buildDataRow('Monthly Burn Rate', '\$${submittedData['monthly_burn_rate']?.toString() ?? '0'}'),
                _buildDataRow('Projected Revenue', '\$${submittedData['projected_revenue']?.toString() ?? '0'}'),
                _buildDataRow('Funding Requirement', '\$${submittedData['funding_requirement']?.toString() ?? '0'}'),
              ],
            ),

            // Founder Information
            if (submittedData['founder_details'] != null && submittedData['founder_details'].isNotEmpty)
              _buildFounderSection(context, submittedData['founder_details'][0]),

            SizedBox(height: 32),

            // Action buttons
            Row(
              children: [
                Expanded(
                  child: OutlinedButton(
                    onPressed: () {
                      Navigator.popUntil(context, (route) => route.isFirst);
                    },
                    style: OutlinedButton.styleFrom(
                      side: BorderSide(color: AppTheme.authAccentColor),
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(30),
                      ),
                      padding: EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: Text(
                      'Back to Home',
                      style: TextStyle(color: AppTheme.authAccentColor),
                    ),
                  ),
                ),
                SizedBox(width: 16),
                Expanded(
                  child: ElevatedButton(
                    onPressed: () {
                      // Get founder name and email from submitted data
                      String userName = 'User';
                      String userEmail = 'user@example.com';
                      
                      if (submittedData['founder_details'] != null && 
                          submittedData['founder_details'].isNotEmpty) {
                        final founder = submittedData['founder_details'][0];
                        userName = founder['name'] ?? 'User';
                        userEmail = founder['email'] ?? 'user@example.com';
                      }
                      
                      Navigator.pushReplacement(
                        context,
                        MaterialPageRoute(
                          builder: (context) => MainNavigation(
                            userName: userName,
                            userEmail: userEmail,
                          ),
                        ),
                      );
                    },
                    style: ElevatedButton.styleFrom(
                      backgroundColor: AppTheme.authAccentColor,
                      shape: RoundedRectangleBorder(
                        borderRadius: BorderRadius.circular(30),
                      ),
                      padding: EdgeInsets.symmetric(vertical: 16),
                    ),
                    child: Text(
                      'Continue to App',
                      style: TextStyle(color: Colors.white),
                    ),
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildSection(BuildContext context, String title, List<Widget> children) {
    return Container(
      width: double.infinity,
      margin: EdgeInsets.only(bottom: 24),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(12),
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
            Text(
              title,
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 4,
                color: AppTheme.authAccentColor,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            ...children,
          ],
        ),
      ),
    );
  }

  Widget _buildDataRow(String label, dynamic value) {
    if (value == null || value.toString().isEmpty) return SizedBox.shrink();
    
    return Padding(
      padding: EdgeInsets.only(bottom: 12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(
              fontWeight: FontWeight.w500,
              color: Colors.grey[700],
              fontSize: 14,
            ),
          ),
          SizedBox(height: 4),
          Text(
            value.toString(),
            style: TextStyle(
              color: Colors.black87,
              fontSize: 16,
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildArrayDataRow(String label, dynamic value) {
    if (value == null || (value is List && value.isEmpty)) return SizedBox.shrink();
    
    List<dynamic> items = value is List ? value : [value];
    
    return Padding(
      padding: EdgeInsets.only(bottom: 12),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            label,
            style: TextStyle(
              fontWeight: FontWeight.w500,
              color: Colors.grey[700],
              fontSize: 14,
            ),
          ),
          SizedBox(height: 4),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: items.map((item) => Container(
              padding: EdgeInsets.symmetric(horizontal: 12, vertical: 6),
              decoration: BoxDecoration(
                color: AppTheme.authAccentColor.withOpacity(0.1),
                borderRadius: BorderRadius.circular(20),
                border: Border.all(color: AppTheme.authAccentColor.withOpacity(0.3)),
              ),
              child: Text(
                item.toString(),
                style: TextStyle(
                  color: AppTheme.authAccentColor,
                  fontSize: 14,
                  fontWeight: FontWeight.w500,
                ),
              ),
            )).toList(),
          ),
        ],
      ),
    );
  }

  Widget _buildFounderSection(BuildContext context, Map<String, dynamic> founder) {
    return _buildSection(
      context,
      'Founder Information',
      [
        _buildDataRow('Name', founder['name']),
        _buildDataRow('Email', founder['email']),
        _buildDataRow('Role', founder['role']),
        _buildDataRow('LinkedIn', founder['linkedin_url']),
        _buildArrayDataRow('Skills', founder['skills']),
        _buildArrayDataRow('Achievements', founder['achievements']),
        
        // Education
        if (founder['education'] != null && founder['education'].isNotEmpty) ...[
          SizedBox(height: 16),
          Text(
            'Education',
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.grey[800],
              fontSize: 16,
            ),
          ),
          SizedBox(height: 8),
          ...founder['education'].map<Widget>((edu) => Container(
            margin: EdgeInsets.only(bottom: 12),
            padding: EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.grey[50],
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.grey[200]!),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (edu['degree'] != null && edu['degree'].isNotEmpty)
                  Text(
                    edu['degree'],
                    style: TextStyle(fontWeight: FontWeight.w500, fontSize: 15),
                  ),
                if (edu['institution'] != null && edu['institution'].isNotEmpty)
                  Text(
                    edu['institution'],
                    style: TextStyle(color: Colors.grey[600], fontSize: 14),
                  ),
                if (edu['graduation_year'] != null && edu['graduation_year'].isNotEmpty)
                  Text(
                    'Graduated: ${edu['graduation_year']}',
                    style: TextStyle(color: Colors.grey[600], fontSize: 14),
                  ),
              ],
            ),
          )).toList(),
        ],
        
        // Experience
        if (founder['experience'] != null && founder['experience'].isNotEmpty) ...[
          SizedBox(height: 16),
          Text(
            'Professional Experience',
            style: TextStyle(
              fontWeight: FontWeight.w600,
              color: Colors.grey[800],
              fontSize: 16,
            ),
          ),
          SizedBox(height: 8),
          ...founder['experience'].map<Widget>((exp) => Container(
            margin: EdgeInsets.only(bottom: 12),
            padding: EdgeInsets.all(12),
            decoration: BoxDecoration(
              color: Colors.grey[50],
              borderRadius: BorderRadius.circular(8),
              border: Border.all(color: Colors.grey[200]!),
            ),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                if (exp['position'] != null && exp['position'].isNotEmpty)
                  Text(
                    exp['position'],
                    style: TextStyle(fontWeight: FontWeight.w500, fontSize: 15),
                  ),
                if (exp['company'] != null && exp['company'].isNotEmpty)
                  Text(
                    exp['company'],
                    style: TextStyle(color: Colors.grey[600], fontSize: 14),
                  ),
                if (exp['start_date'] != null && exp['start_date'].isNotEmpty) ...[
                  SizedBox(height: 4),
                  Text(
                    '${exp['start_date']} - ${exp['end_date']?.isNotEmpty == true ? exp['end_date'] : 'Present'}',
                    style: TextStyle(color: Colors.grey[600], fontSize: 13),
                  ),
                ],
                if (exp['description'] != null && exp['description'].isNotEmpty) ...[
                  SizedBox(height: 8),
                  Text(
                    exp['description'],
                    style: TextStyle(color: Colors.grey[700], fontSize: 14),
                  ),
                ],
              ],
            ),
          )).toList(),
        ],
      ],
    );
  }
}

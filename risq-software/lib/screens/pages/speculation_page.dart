import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/services/speculation_service.dart';
import 'package:risq/screens/pages/startup_profile_page.dart';

class SpeculationPage extends StatefulWidget {
  final String userName;
  final String userEmail;

  const SpeculationPage({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<SpeculationPage> createState() => _SpeculationPageState();
}

class _SpeculationPageState extends State<SpeculationPage> {
  final GlobalKey<FormState> _formKey = GlobalKey<FormState>();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _contextController = TextEditingController();
  final TextEditingController _timelineController = TextEditingController();
  final TextEditingController _budgetController = TextEditingController();
  
  String _selectedCategory = 'hiring';
  bool _isLoading = false;

  final List<Map<String, dynamic>> _categories = [
    {'value': 'hiring', 'label': 'Team & Hiring', 'icon': Icons.people},
    {'value': 'funding', 'label': 'Funding & Investment', 'icon': Icons.attach_money},
    {'value': 'product', 'label': 'Product Development', 'icon': Icons.build},
    {'value': 'marketing', 'label': 'Marketing & Sales', 'icon': Icons.campaign},
    {'value': 'operations', 'label': 'Operations', 'icon': Icons.settings},
    {'value': 'strategy', 'label': 'Business Strategy', 'icon': Icons.business},
    {'value': 'legal', 'label': 'Legal & Compliance', 'icon': Icons.gavel},
    {'value': 'other', 'label': 'Other', 'icon': Icons.help_outline},
  ];

  @override
  void dispose() {
    _descriptionController.dispose();
    _contextController.dispose();
    _timelineController.dispose();
    _budgetController.dispose();
    super.dispose();
  }

  Future<void> _submitSpeculation() async {
    if (!_formKey.currentState!.validate()) return;

    setState(() {
      _isLoading = true;
    });

    try {
      // Call the actual API
      final result = await SpeculationService.submitSpeculation(
        description: _descriptionController.text,
        category: _selectedCategory,
        context: _contextController.text,
        timeline: _timelineController.text,
        budget: _budgetController.text,
        // TODO: Add auth token when authentication is implemented
        // authToken: await AuthService.getAuthToken(),
      );

      if (mounted) {
        setState(() {
          _isLoading = false;
        });

        if (result['success']) {
          // Show result dialog with real data
          _showSpeculationResults(result['data']);
        } else {
          // Show error
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(result['message'] ?? 'Failed to analyze speculation'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    } catch (e) {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
        
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(
            content: Text('An error occurred. Please try again.'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  void _showSpeculationResults([Map<String, dynamic>? analysisData]) {
    // Use real data if available, otherwise use mock data
    final riskScore = analysisData?['riskScore']?.toString() ?? '7.8';
    final confidence = analysisData?['confidence']?.toString() ?? '82%';
    final analysis = analysisData?['analysis'] ?? 
        'This decision may increase operational complexity and burn rate. Consider implementing gradual hiring with clear performance metrics.';
    final recommendations = analysisData?['recommendations'] as List<dynamic>? ?? [];
    
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
        title: Row(
          children: [
            Icon(Icons.psychology, color: AppTheme.authAccentColor),
            SizedBox(width: 8),
            Text('AI Analysis Results'),
          ],
        ),
        content: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            _buildResultItem('Previous Risk Score', '7.2', Colors.orange),
            SizedBox(height: 8),
            _buildResultItem('Projected Risk Score', riskScore, 
                double.tryParse(riskScore) != null && double.parse(riskScore) > 7.2 ? Colors.red : Colors.green),
            SizedBox(height: 8),
            _buildResultItem('Risk Change', 
                double.tryParse(riskScore) != null ? 
                    (double.parse(riskScore) - 7.2 > 0 ? '+${(double.parse(riskScore) - 7.2).toStringAsFixed(1)}' : '${(double.parse(riskScore) - 7.2).toStringAsFixed(1)}') : 
                    '+0.6', 
                double.tryParse(riskScore) != null && double.parse(riskScore) > 7.2 ? Colors.red : Colors.green),
            SizedBox(height: 8),
            _buildResultItem('AI Confidence', confidence, Colors.blue),
            SizedBox(height: 16),
            Text(
              'AI Reasoning:',
              style: TextStyle(fontWeight: FontWeight.bold),
            ),
            SizedBox(height: 8),
            Text(
              analysis,
              style: TextStyle(color: Colors.grey[700]),
            ),
            if (recommendations.isNotEmpty) ...[
              SizedBox(height: 16),
              Text(
                'Recommendations:',
                style: TextStyle(fontWeight: FontWeight.bold),
              ),
              SizedBox(height: 8),
              ...recommendations.map((rec) => Padding(
                padding: EdgeInsets.only(bottom: 4),
                child: Text(
                  '• $rec',
                  style: TextStyle(color: Colors.grey[700], fontSize: 12),
                ),
              )).toList(),
            ],
          ],
        ),
        actions: [
          TextButton(
            onPressed: () => Navigator.pop(context),
            child: Text('Save as Speculation'),
          ),
          ElevatedButton(
            onPressed: () {
              Navigator.pop(context);
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
            child: Text('Confirm Decision', style: TextStyle(color: Colors.white)),
          ),
        ],
      ),
    );
  }

  Widget _buildResultItem(String label, String value, Color color) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        Text(label),
        Text(
          value,
          style: TextStyle(
            fontWeight: FontWeight.bold,
            color: color,
          ),
        ),
      ],
    );
  }

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
                        'Decision Speculation',
                        style: AppTheme.headingTextStyle.copyWith(
                          fontSize: ResponsiveUtils.getBodySize(context) + 4,
                          color: Colors.black87,
                        ),
                      ),
                      Text(
                        'Test your business decisions',
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
            
            // Main Content
            Expanded(
              child: SingleChildScrollView(
                padding: EdgeInsets.all(20),
                child: Form(
                  key: _formKey,
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      // Info Card
                      Container(
                        width: double.infinity,
                        padding: EdgeInsets.all(20),
                        decoration: BoxDecoration(
                          color: Colors.blue[50],
                          borderRadius: BorderRadius.circular(12),
                          border: Border.all(color: Colors.blue[200]!),
                        ),
                        child: Row(
                          children: [
                            Icon(Icons.info_outline, color: Colors.blue[700]),
                            SizedBox(width: 12),
                            Expanded(
                              child: Text(
                                'Describe a decision you\'re considering and get AI-powered risk analysis before committing.',
                                style: TextStyle(
                                  color: Colors.blue[700],
                                  fontSize: ResponsiveUtils.getBodySize(context),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      SizedBox(height: 24),
                      
                      // Decision Category
                      Text(
                        'Decision Category',
                        style: AppTheme.labelTextStyle.copyWith(
                          fontSize: ResponsiveUtils.getBodySize(context),
                          color: Colors.black87,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                      SizedBox(height: 12),
                      Container(
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
                        child: DropdownButtonFormField<String>(
                          value: _selectedCategory,
                          decoration: InputDecoration(
                            contentPadding: EdgeInsets.symmetric(horizontal: 16, vertical: 12),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide.none,
                            ),
                          ),
                          items: _categories.map((category) {
                            return DropdownMenuItem<String>(
                              value: category['value'],
                              child: Row(
                                children: [
                                  Icon(category['icon'], size: 20, color: Colors.grey[600]),
                                  SizedBox(width: 12),
                                  Text(category['label']),
                                ],
                              ),
                            );
                          }).toList(),
                          onChanged: (value) {
                            setState(() {
                              _selectedCategory = value!;
                            });
                          },
                        ),
                      ),
                      SizedBox(height: 20),
                      
                      // Decision Description
                      _buildInputField(
                        label: 'Decision Description',
                        controller: _descriptionController,
                        hintText: 'e.g., Hire 3 senior engineers to accelerate product development',
                        maxLines: 3,
                        validator: (value) {
                          if (value == null || value.isEmpty) {
                            return 'Please describe your decision';
                          }
                          return null;
                        },
                      ),
                      
                      // Context
                      _buildInputField(
                        label: 'Context (Optional)',
                        controller: _contextController,
                        hintText: 'Why are you considering this decision?',
                        maxLines: 3,
                      ),
                      
                      // Timeline
                      _buildInputField(
                        label: 'Timeline (Optional)',
                        controller: _timelineController,
                        hintText: 'e.g., 2 months, Q2 2025',
                      ),
                      
                      // Budget
                      _buildInputField(
                        label: 'Budget (Optional)',
                        controller: _budgetController,
                        hintText: 'e.g., 180000',
                        keyboardType: TextInputType.number,
                        prefix: '\$ ',
                      ),
                      
                      SizedBox(height: 32),
                      
                      // Submit Button
                      Container(
                        width: double.infinity,
                        height: 56,
                        child: ElevatedButton(
                          onPressed: _isLoading ? null : _submitSpeculation,
                          style: ElevatedButton.styleFrom(
                            backgroundColor: AppTheme.authAccentColor,
                            shape: RoundedRectangleBorder(
                              borderRadius: BorderRadius.circular(16),
                            ),
                            elevation: 0,
                          ),
                          child: _isLoading
                              ? CircularProgressIndicator(color: Colors.white)
                              : Row(
                                  mainAxisAlignment: MainAxisAlignment.center,
                                  children: [
                                    Icon(Icons.psychology, color: Colors.white),
                                    SizedBox(width: 8),
                                    Text(
                                      'Analyze Decision',
                                      style: TextStyle(
                                        color: Colors.white,
                                        fontSize: ResponsiveUtils.getBodySize(context),
                                        fontWeight: FontWeight.w600,
                                      ),
                                    ),
                                  ],
                                ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildInputField({
    required String label,
    required TextEditingController controller,
    String? hintText,
    int maxLines = 1,
    TextInputType? keyboardType,
    String? Function(String?)? validator,
    String? prefix,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context),
            color: Colors.black87,
            fontWeight: FontWeight.w600,
          ),
        ),
        SizedBox(height: 8),
        Container(
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
          child: TextFormField(
            controller: controller,
            validator: validator,
            maxLines: maxLines,
            keyboardType: keyboardType,
            style: AppTheme.inputTextStyle.copyWith(
              fontSize: ResponsiveUtils.getBodySize(context),
              color: Colors.black87,
            ),
            decoration: InputDecoration(
              hintText: hintText,
              prefixText: prefix,
              hintStyle: TextStyle(
                color: Colors.grey[500],
                fontSize: ResponsiveUtils.getBodySize(context),
              ),
              contentPadding: EdgeInsets.all(16),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide.none,
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: AppTheme.authAccentColor, width: 2),
              ),
              errorBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
                borderSide: BorderSide(color: Colors.red, width: 1),
              ),
            ),
          ),
        ),
        SizedBox(height: 16),
      ],
    );
  }
}

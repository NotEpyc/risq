import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/services/startup_service.dart';
import 'package:risq/services/auth_service.dart';
import 'package:risq/screens/login/signin_page.dart';
import 'package:risq/utils/responsive_utils.dart';

class StartupProfilePage extends StatefulWidget {
  const StartupProfilePage({super.key});

  @override
  State<StartupProfilePage> createState() => _StartupProfilePageState();
}

class _StartupProfilePageState extends State<StartupProfilePage> 
    with TickerProviderStateMixin {
  
  // Text controllers for form fields
  final TextEditingController _nameController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _industryController = TextEditingController();
  final TextEditingController _locationController = TextEditingController();
  final TextEditingController _teamSizeController = TextEditingController();
  final TextEditingController _websiteController = TextEditingController();

  // State variables
  bool _isEditing = false;
  bool _isLoading = true;
  String? _selectedFundingStage;
  DateTime? _foundedDate;
  Map<String, dynamic> _startupData = {};

  // Funding stage options from API documentation
  final List<String> _fundingStageOptions = [
    'idea',
    'pre_seed', 
    'seed',
    'series_a',
    'series_b',
    'series_c',
    'ipo'
  ];

  // Display names for funding stages
  final Map<String, String> _fundingStageDisplayNames = {
    'idea': 'Idea Stage',
    'pre_seed': 'Pre-Seed',
    'seed': 'Seed',
    'series_a': 'Series A',
    'series_b': 'Series B', 
    'series_c': 'Series C',
    'ipo': 'IPO',
  };

  @override
  void initState() {
    super.initState();
    _loadStartupProfile();
  }

  Future<void> _loadStartupProfile() async {
    try {
      final result = await StartupService.getStartupProfile();
      
      if (result['success']) {
        final startupData = result['data']['data'] ?? result['data'];
        
        if (mounted) {
          setState(() {
            _startupData = startupData;
            _nameController.text = startupData['name'] ?? '';
            _descriptionController.text = startupData['description'] ?? '';
            _industryController.text = startupData['industry'] ?? '';
            _locationController.text = startupData['location'] ?? '';
            _teamSizeController.text = startupData['team_size']?.toString() ?? '';
            _websiteController.text = startupData['website'] ?? '';
            _selectedFundingStage = startupData['funding_stage'];
            
            // Parse founded date
            if (startupData['founded_date'] != null) {
              try {
                _foundedDate = DateTime.parse(startupData['founded_date']);
              } catch (e) {
                debugPrint('Error parsing founded date: $e');
              }
            }
            
            _isLoading = false;
          });
        }
      } else {
        if (mounted) {
          setState(() => _isLoading = false);
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(
              content: Text(result['message'] ?? 'Failed to load startup profile'),
              backgroundColor: Colors.red,
            ),
          );
        }
      }
    } catch (e) {
      debugPrint('Error loading startup profile: $e');
      if (mounted) {
        setState(() => _isLoading = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Failed to load startup profile'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _saveChanges() async {
    if (!_isEditing) return;

    // Basic validation
    if (_nameController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Startup name is required')),
      );
      return;
    }

    if (_descriptionController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('Description is required')),
      );
      return;
    }

    // Show loading dialog
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const AlertDialog(
        content: Row(
          children: [
            CircularProgressIndicator(),
            SizedBox(width: 16),
            Text('Updating profile...'),
          ],
        ),
      ),
    );

    try {
      // Note: The current API doesn't have an update endpoint, 
      // so this is a placeholder for when that's implemented
      // For now, we'll just show a message
      
      await Future.delayed(const Duration(seconds: 1)); // Simulate API call
      
      if (mounted) {
        Navigator.pop(context); // Close loading dialog
        setState(() => _isEditing = false);
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Profile update feature coming soon!'),
            backgroundColor: Colors.orange,
          ),
        );
      }
    } catch (e) {
      if (mounted) {
        Navigator.pop(context); // Close loading dialog
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text('Failed to update profile'),
            backgroundColor: Colors.red,
          ),
        );
      }
    }
  }

  Future<void> _signOut() async {
    try {
      final confirm = await showDialog<bool>(
        context: context,
        builder: (context) => AlertDialog(
          title: const Text('Confirm Logout'),
          content: const Text('Are you sure you want to logout?'),
          actions: [
            TextButton(
              onPressed: () => Navigator.pop(context, false),
              child: const Text(
                'Cancel',
                style: TextStyle(color: Colors.grey),
              ),
            ),
            TextButton(
              onPressed: () => Navigator.pop(context, true),
              child: Text(
                'Logout',
                style: TextStyle(color: AppTheme.authAccentColor),
              ),
            ),
          ],
        ),
      );

      if (confirm == true && mounted) {
        await AuthService.logout();
        if (mounted) {
          Navigator.of(context).pushAndRemoveUntil(
            MaterialPageRoute(builder: (context) => const SignInPage()),
            (route) => false,
          );
        }
      }
    } catch (e) {
      ScaffoldMessenger.of(context).showSnackBar(
        SnackBar(content: Text('Failed to log out: ${e.toString()}')),
      );
    }
  }

  String _formatFundingStage(String? stage) {
    if (stage == null) return '-';
    return _fundingStageDisplayNames[stage] ?? stage.toUpperCase();
  }

  String _getTimeSinceCreation() {
    if (_startupData['created_at'] == null) return '-';
    
    try {
      final createdAt = DateTime.parse(_startupData['created_at']);
      final now = DateTime.now();
      final difference = now.difference(createdAt);
      
      if (difference.inDays > 0) {
        return '${difference.inDays} days ago';
      } else if (difference.inHours > 0) {
        return '${difference.inHours} hours ago';
      } else {
        return '${difference.inMinutes} minutes ago';
      }
    } catch (e) {
      return '-';
    }
  }

  @override
  Widget build(BuildContext context) {
    if (_isLoading) {
      return Scaffold(
        backgroundColor: Colors.white,
        body: Container(
          decoration: BoxDecoration(
            gradient: LinearGradient(
              begin: Alignment.topLeft,
              end: Alignment.bottomRight,
              colors: [
                AppTheme.authSecondaryColor.withOpacity(0.1),
                AppTheme.authPrimaryColor.withOpacity(0.1),
              ],
            ),
          ),
          child: Center(
            child: Column(
              mainAxisAlignment: MainAxisAlignment.center,
              children: [
                CircularProgressIndicator(
                  valueColor: AlwaysStoppedAnimation<Color>(AppTheme.authAccentColor),
                  strokeWidth: 3,
                ),
                const SizedBox(height: 16),
                Text(
                  'Loading startup profile...',
                  style: AppTheme.regularTextStyle.copyWith(
                    color: AppTheme.authAccentColor,
                    fontSize: ResponsiveUtils.getBodySize(context),
                  ),
                ),
              ],
            ),
          ),
        ),
      );
    }

    return Scaffold(
      extendBodyBehindAppBar: true,
      appBar: AppBar(
        backgroundColor: Colors.transparent,
        elevation: 0,
        systemOverlayStyle: SystemUiOverlayStyle.dark,
        leading: IconButton(
          icon: Icon(
            Icons.arrow_back,
            color: AppTheme.authAccentColor,
            size: 28,
          ),
          onPressed: () => Navigator.of(context).pop(),
        ),
        actions: [
          if (_isEditing)
            IconButton(
              icon: Icon(Icons.save, color: AppTheme.authAccentColor),
              onPressed: _saveChanges,
            )
          else
            IconButton(
              icon: Icon(Icons.edit, color: AppTheme.authAccentColor),
              onPressed: () => setState(() => _isEditing = true),
            ),
        ],
      ),
      body: Container(
        decoration: BoxDecoration(
          gradient: LinearGradient(
            begin: Alignment.topLeft,
            end: Alignment.bottomRight,
            colors: [
              AppTheme.authSecondaryColor.withOpacity(0.05),
              AppTheme.authPrimaryColor.withOpacity(0.05),
            ],
          ),
        ),
        child: SingleChildScrollView(
          padding: const EdgeInsets.fromLTRB(16, 100, 16, 16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Startup Header Section
              Center(
                child: Column(
                  children: [
                    // Startup Logo/Icon
                    Container(
                      padding: const EdgeInsets.all(20),
                      decoration: BoxDecoration(
                        gradient: LinearGradient(
                          colors: [
                            AppTheme.authSecondaryColor,
                            AppTheme.authAccentColor,
                          ],
                        ),
                        shape: BoxShape.circle,
                        boxShadow: [
                          BoxShadow(
                            color: AppTheme.authAccentColor.withOpacity(0.3),
                            blurRadius: 15,
                            offset: const Offset(0, 5),
                          ),
                        ],
                      ),
                      child: const Icon(
                        Icons.business,
                        size: 40,
                        color: Colors.white,
                      ),
                    ),
                    const SizedBox(height: 16),
                    
                    // Startup Name
                    Text(
                      _startupData['name'] ?? 'Startup Name',
                      style: AppTheme.headingTextStyle.copyWith(
                        fontSize: ResponsiveUtils.getHeadingSize(context) * 1.2,
                        color: AppTheme.authDarkBlue,
                        fontWeight: FontWeight.bold,
                      ),
                      textAlign: TextAlign.center,
                    ),
                    const SizedBox(height: 8),
                    
                    // Profile Created Badge
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                      decoration: BoxDecoration(
                        color: AppTheme.authAccentColor.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(20),
                        border: Border.all(
                          color: AppTheme.authAccentColor.withOpacity(0.3),
                        ),
                      ),
                      child: Row(
                        mainAxisSize: MainAxisSize.min,
                        children: [
                          Icon(
                            Icons.access_time,
                            color: AppTheme.authAccentColor,
                            size: 16,
                          ),
                          const SizedBox(width: 8),
                          Text(
                            'Created ${_getTimeSinceCreation()}',
                            style: TextStyle(
                              color: AppTheme.authAccentColor,
                              fontSize: ResponsiveUtils.getSmallTextSize(context),
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ),
              
              const SizedBox(height: 32),
              
              // Basic Information Section
              _buildSectionHeader('Basic Information'),
              const SizedBox(height: 16),
              
              _buildTextField(
                controller: _nameController,
                label: 'Startup Name',
                enabled: _isEditing,
                readOnly: !_isEditing,
                icon: Icons.business_outlined,
              ),
              const SizedBox(height: 16),
              
              _buildTextField(
                controller: _descriptionController,
                label: 'Description',
                enabled: _isEditing,
                readOnly: !_isEditing,
                icon: Icons.description_outlined,
                maxLines: 3,
                minLines: 1,
                keyboardType: TextInputType.multiline,
              ),
              const SizedBox(height: 16),
              
              _buildTextField(
                controller: _industryController,
                label: 'Industry',
                enabled: _isEditing,
                readOnly: !_isEditing,
                icon: Icons.category_outlined,
              ),
              const SizedBox(height: 16),
              
              // Company Details Section
              _buildSectionHeader('Company Details'),
              const SizedBox(height: 16),
              
              // Funding Stage and Team Size Row
              Row(
                children: [
                  Expanded(
                    flex: 2,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Padding(
                          padding: const EdgeInsets.only(left: 4, bottom: 5),
                          child: Text(
                            'Funding Stage',
                            style: AppTheme.labelTextStyle.copyWith(
                              fontSize: ResponsiveUtils.getSmallTextSize(context),
                              color: Colors.black87,
                            ),
                          ),
                        ),
                        Container(
                          decoration: BoxDecoration(
                            color: Colors.white,
                            borderRadius: BorderRadius.circular(30),
                            border: Border.all(
                              color: AppTheme.authSecondaryColor.withOpacity(0.3),
                              width: 1,
                            ),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withOpacity(0.04),
                                blurRadius: 4,
                                offset: const Offset(0, 2),
                              ),
                            ],
                          ),
                          child: DropdownButtonFormField<String>(
                            value: _selectedFundingStage,
                            decoration: InputDecoration(
                              prefixIcon: Icon(
                                Icons.trending_up_outlined,
                                color: AppTheme.authAccentColor,
                                size: 20,
                              ),
                              contentPadding: const EdgeInsets.symmetric(
                                horizontal: 20, 
                                vertical: 14
                              ),
                              border: OutlineInputBorder(
                                borderRadius: BorderRadius.circular(30),
                                borderSide: BorderSide.none,
                              ),
                              filled: true,
                              fillColor: Colors.white,
                              enabled: _isEditing,
                            ),
                            items: _fundingStageOptions.map((String stage) {
                              return DropdownMenuItem<String>(
                                value: stage,
                                child: Text(
                                  _formatFundingStage(stage),
                                  style: TextStyle(
                                    fontSize: ResponsiveUtils.getBodySize(context),
                                    color: _isEditing ? Colors.black87 : Colors.grey[600],
                                    fontWeight: FontWeight.normal,
                                  ),
                                ),
                              );
                            }).toList(),
                            onChanged: _isEditing ? (String? newValue) {
                              setState(() => _selectedFundingStage = newValue);
                            } : null,
                            icon: Icon(
                              Icons.arrow_drop_down,
                              color: _isEditing ? AppTheme.authAccentColor : Colors.grey,
                            ),
                            dropdownColor: Colors.white,
                          ),
                        ),
                      ],
                    ),
                  ),
                  const SizedBox(width: 16),
                  Expanded(
                    child: _buildTextField(
                      controller: _teamSizeController,
                      label: 'Team Size',
                      enabled: _isEditing,
                      readOnly: !_isEditing,
                      icon: Icons.group_outlined,
                      keyboardType: TextInputType.number,
                    ),
                  ),
                ],
              ),
              
              const SizedBox(height: 16),
              
              _buildTextField(
                controller: _locationController,
                label: 'Location',
                enabled: _isEditing,
                readOnly: !_isEditing,
                icon: Icons.location_on_outlined,
              ),
              const SizedBox(height: 16),
              
              // Founded Date
              Padding(
                padding: const EdgeInsets.only(left: 4, bottom: 5),
                child: Text(
                  'Founded Date',
                  style: AppTheme.labelTextStyle.copyWith(
                    fontSize: ResponsiveUtils.getSmallTextSize(context),
                    color: Colors.black87,
                  ),
                ),
              ),
              GestureDetector(
                onTap: _isEditing ? () async {
                  final DateTime? picked = await showDatePicker(
                    context: context,
                    initialDate: _foundedDate ?? DateTime.now(),
                    firstDate: DateTime(1900),
                    lastDate: DateTime.now(),
                  );
                  if (picked != null && mounted) {
                    setState(() => _foundedDate = picked);
                  }
                } : null,
                child: Container(
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.circular(30),
                    border: Border.all(
                      color: AppTheme.authSecondaryColor.withOpacity(0.3),
                      width: 1,
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.04),
                        blurRadius: 4,
                        offset: const Offset(0, 2),
                      ),
                    ],
                  ),
                  child: ListTile(
                    leading: Icon(
                      Icons.calendar_today,
                      color: AppTheme.authAccentColor,
                      size: 20,
                    ),
                    title: Text(
                      _foundedDate != null 
                        ? "${_foundedDate!.day}/${_foundedDate!.month}/${_foundedDate!.year}"
                        : "Select Date",
                      style: TextStyle(
                        color: _isEditing ? Colors.black87 : Colors.grey[600],
                        fontSize: ResponsiveUtils.getBodySize(context),
                      ),
                    ),
                    trailing: _isEditing ? Icon(
                      Icons.arrow_drop_down,
                      color: AppTheme.authAccentColor,
                    ) : null,
                    contentPadding: const EdgeInsets.symmetric(horizontal: 20),
                  ),
                ),
              ),
              
              const SizedBox(height: 16),
              
              _buildTextField(
                controller: _websiteController,
                label: 'Website',
                enabled: _isEditing,
                readOnly: !_isEditing,
                icon: Icons.language_outlined,
                keyboardType: TextInputType.url,
              ),
              
              const SizedBox(height: 32),
              
              // Logout Button
              Container(
                margin: const EdgeInsets.symmetric(horizontal: 16),
                child: TextButton(
                  onPressed: _signOut,
                  style: TextButton.styleFrom(
                    backgroundColor: Colors.white,
                    padding: const EdgeInsets.symmetric(vertical: 16),
                    shape: RoundedRectangleBorder(
                      borderRadius: BorderRadius.circular(30),
                      side: BorderSide(color: AppTheme.authAccentColor),
                    ),
                  ),
                  child: Row(
                    mainAxisAlignment: MainAxisAlignment.center,
                    children: [
                      Icon(
                        Icons.logout,
                        color: AppTheme.authAccentColor,
                        size: 20,
                      ),
                      const SizedBox(width: 8),
                      Text(
                        'Logout',
                        style: TextStyle(
                          color: AppTheme.authAccentColor,
                          fontSize: ResponsiveUtils.getBodySize(context),
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              
              const SizedBox(height: 16),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSectionHeader(String title) {
    return Padding(
      padding: const EdgeInsets.only(left: 4),
      child: Text(
        title,
        style: AppTheme.headingTextStyle.copyWith(
          fontSize: ResponsiveUtils.getTitleSize(context),
          color: AppTheme.authDarkBlue,
          fontWeight: FontWeight.w600,
        ),
      ),
    );
  }

  Widget _buildTextField({
    required TextEditingController controller,
    required String label,
    required IconData icon,
    bool enabled = true,
    bool readOnly = false,
    TextInputType? keyboardType,
    int maxLines = 1,
    int? minLines,
    String? prefixText,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 4, bottom: 5),
          child: Text(
            label,
            style: AppTheme.labelTextStyle.copyWith(
              fontSize: ResponsiveUtils.getSmallTextSize(context),
              color: Colors.black87,
            ),
          ),
        ),
        Container(
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(30),
            border: Border.all(
              color: AppTheme.authSecondaryColor.withOpacity(0.3),
              width: 1,
            ),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.04),
                blurRadius: 4,
                offset: const Offset(0, 2),
              ),
            ],
          ),
          child: TextFormField(
            controller: controller,
            enabled: enabled,
            readOnly: readOnly,
            keyboardType: keyboardType ?? TextInputType.text,
            maxLines: maxLines,
            minLines: minLines,
            style: TextStyle(
              color: enabled ? Colors.black87 : Colors.grey[600],
              fontSize: ResponsiveUtils.getBodySize(context),
            ),
            decoration: InputDecoration(
              prefixIcon: Icon(
                icon,
                color: AppTheme.authAccentColor,
                size: 20,
              ),
              prefixText: prefixText,
              contentPadding: const EdgeInsets.symmetric(
                horizontal: 20, 
                vertical: 16
              ),
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide.none,
              ),
              filled: true,
              fillColor: Colors.white,
              hintText: enabled ? 'Enter $label' : null,
              hintStyle: TextStyle(
                color: Colors.grey[400],
                fontSize: ResponsiveUtils.getBodySize(context),
              ),
            ),
            onChanged: _isEditing ? (value) {
              // Only allow changes when editing
            } : null,
          ),
        ),
      ],
    );
  }

  @override
  void dispose() {
    _nameController.dispose();
    _descriptionController.dispose();
    _industryController.dispose();
    _locationController.dispose();
    _teamSizeController.dispose();
    _websiteController.dispose();
    super.dispose();
  }
}

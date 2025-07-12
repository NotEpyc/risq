import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:video_player/video_player.dart';
import 'dart:math' as math;
import 'dart:convert';
import 'package:risq/screens/main_navigation.dart';

class OnboardingPage extends StatefulWidget {
  final String userName;
  final String userEmail;

  const OnboardingPage({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<OnboardingPage> createState() => _OnboardingPageState();
}

class _OnboardingPageState extends State<OnboardingPage> {
  final PageController _pageController = PageController();
  int _currentStep = 0;
  final int _totalSteps = 4;
  
  // Form keys for each step
  final List<GlobalKey<FormState>> _formKeys = List.generate(4, (index) => GlobalKey<FormState>());
  
  String _errorMessage = '';
  bool _isLoading = false;
  
  // Video player
  VideoPlayerController? _videoController;
  bool _videoFailed = false;

  // Random gradient for button
  late Alignment _gradientBegin;
  late Alignment _gradientEnd;

  // Step 1: Company Basic Info
  final TextEditingController _companyNameController = TextEditingController();
  final TextEditingController _descriptionController = TextEditingController();
  final TextEditingController _industryController = TextEditingController();
  final TextEditingController _sectorController = TextEditingController();
  final TextEditingController _locationController = TextEditingController();
  final TextEditingController _websiteController = TextEditingController();
  
  // Step 2: Business Model & Market
  final TextEditingController _businessModelController = TextEditingController();
  final TextEditingController _targetMarketController = TextEditingController();
  final TextEditingController _competitiveAdvantageController = TextEditingController();
  final TextEditingController _goToMarketController = TextEditingController();
  
  String _fundingStage = 'Seed';
  int _teamSize = 1;
  List<String> _revenueStreams = [];
  
  // Step 3: Technical & Implementation
  final TextEditingController _implementationPlanController = TextEditingController();
  final TextEditingController _developmentTimelineController = TextEditingController();
  
  List<String> _technologyStack = [];
  
  // Step 4: Financial & Founder Info
  final TextEditingController _initialInvestmentController = TextEditingController();
  final TextEditingController _monthlyBurnRateController = TextEditingController();
  final TextEditingController _projectedRevenueController = TextEditingController();
  final TextEditingController _fundingRequirementController = TextEditingController();
  
  // Founder Details
  final TextEditingController _founderNameController = TextEditingController();
  final TextEditingController _founderEmailController = TextEditingController();
  final TextEditingController _founderRoleController = TextEditingController();
  final TextEditingController _founderLinkedInController = TextEditingController();
  final TextEditingController _founderSkillsController = TextEditingController();
  final TextEditingController _founderAchievementsController = TextEditingController();
  
  // Experience list for multiple experiences
  List<Map<String, dynamic>> _experiences = [
    {
      'company': TextEditingController(),
      'position': TextEditingController(),
      'start_date': TextEditingController(),
      'end_date': TextEditingController(),
      'description': TextEditingController(),
    }
  ];

  // Education list for multiple education entries
  List<Map<String, dynamic>> _educations = [];

  void _addEducation() {
    setState(() {
      _educations.add({
        'degree': TextEditingController(),
        'institution': TextEditingController(),
        'graduation_year': TextEditingController(),
      });
    });
  }

  void _removeEducation(int index) {
    if (_educations.length > 1) {
      setState(() {
        // Dispose controllers before removing
        (_educations[index]['degree'] as TextEditingController).dispose();
        (_educations[index]['institution'] as TextEditingController).dispose();
        (_educations[index]['graduation_year'] as TextEditingController).dispose();
        _educations.removeAt(index);
      });
    }
  }

  @override
  void initState() {
    super.initState();
    _initializeVideo();
    _initializeGradient();
    _founderEmailController.text = widget.userEmail;
    _founderNameController.text = widget.userName;
    
    // Add initial education entry
    _addEducation();
    
    // No authentication verification needed for dummy data
    print('Using dummy data - no authentication required');
  }

  void _initializeVideo() {
    try {
      _videoController = VideoPlayerController.asset('assets/videos/bg_video.mp4');
      _videoController!.initialize().then((_) {
        if (mounted) {
          setState(() {});
          _videoController!.play();
          _videoController!.setLooping(true);
        }
      }).catchError((error) {
        if (mounted) {
          setState(() {
            _videoFailed = true;
          });
        }
      });
    } catch (e) {
      if (mounted) {
        setState(() {
          _videoFailed = true;
        });
      }
    }
  }

  void _initializeGradient() {
    double angle = math.Random().nextDouble() * 2 * math.pi;
    _gradientBegin = Alignment(math.cos(angle), math.sin(angle));
    _gradientEnd = Alignment(-math.cos(angle), -math.sin(angle));
  }

  @override
  void dispose() {
    // Step 1 controllers
    _companyNameController.dispose();
    _descriptionController.dispose();
    _industryController.dispose();
    _sectorController.dispose();
    _locationController.dispose();
    _websiteController.dispose();
    
    // Step 2 controllers
    _businessModelController.dispose();
    _targetMarketController.dispose();
    _competitiveAdvantageController.dispose();
    _goToMarketController.dispose();
    
    //
    _implementationPlanController.dispose();
    _developmentTimelineController.dispose();
    
    // Step 4 controllers
    _initialInvestmentController.dispose();
    _monthlyBurnRateController.dispose();
    _projectedRevenueController.dispose();
    _fundingRequirementController.dispose();
    
    // Founder controllers
    _founderNameController.dispose();
    _founderEmailController.dispose();
    _founderRoleController.dispose();
    _founderLinkedInController.dispose();
    _founderSkillsController.dispose();
    _founderAchievementsController.dispose();
    
    // Dispose experience controllers
    for (var experience in _experiences) {
      (experience['company'] as TextEditingController).dispose();
      (experience['position'] as TextEditingController).dispose();
      (experience['start_date'] as TextEditingController).dispose();
      (experience['end_date'] as TextEditingController).dispose();
      (experience['description'] as TextEditingController).dispose();
    }
    
    // Dispose education controllers
    for (var education in _educations) {
      (education['degree'] as TextEditingController).dispose();
      (education['institution'] as TextEditingController).dispose();
      (education['graduation_year'] as TextEditingController).dispose();
    }
    
    _videoController?.dispose();
    _pageController.dispose();
    super.dispose();
  }

  // Helper method to convert team size to appropriate range string
  String _getTeamSizeRange(int teamSize) {
    if (teamSize <= 1) return "1";
    if (teamSize <= 5) return "1-5";
    if (teamSize <= 10) return "6-10";
    if (teamSize <= 20) return "11-20";
    if (teamSize <= 50) return "21-50";
    return "50+";
  }

  void _nextStep() {
    print('=== FORM VALIDATION ===');
    print('Current step: $_currentStep');
    print('Total steps: $_totalSteps');
    print('Validating form for step: $_currentStep');
    
    final isValid = _formKeys[_currentStep].currentState!.validate();
    print('Form validation result: $isValid');
    
    if (isValid) {
      if (_currentStep < _totalSteps - 1) {
        print('Moving to next step...');
        setState(() {
          _currentStep++;
        });
        _pageController.nextPage(
          duration: Duration(milliseconds: 300),
          curve: Curves.easeInOut,
        );
      } else {
        print('Last step - submitting onboarding...');
        _submitOnboarding();
      }
    } else {
      print('❌ Form validation failed for step $_currentStep');
      setState(() {
        _errorMessage = 'Please fill in all required fields correctly.';
      });
    }
  }

  void _previousStep() {
    if (_currentStep > 0) {
      setState(() {
        _currentStep--;
      });
      _pageController.previousPage(
        duration: Duration(milliseconds: 300),
        curve: Curves.easeInOut,
      );
    }
  }

  Future<void> _submitOnboarding() async {
    setState(() {
      _errorMessage = '';
      _isLoading = true;
    });

    try {
      print('=== DUMMY ONBOARDING SUBMISSION ===');
      print('User: ${widget.userName} (${widget.userEmail})');
      
      // Create the JSON payload with validation
      final payload = {
        "name": _companyNameController.text.trim(),
        "description": _descriptionController.text.trim(),
        "industry": _industryController.text.trim(),
        "sector": _sectorController.text.trim(),
        "funding_stage": _fundingStage,
        "location": _locationController.text.trim(),
        "founded_date": DateTime.now().toUtc().toIso8601String().split('T')[0],
        "team_size": _getTeamSizeRange(_teamSize),
        "website": _websiteController.text.trim().isEmpty ? null : _websiteController.text.trim(),
        "business_model": _businessModelController.text.trim(),
        "revenue_streams": _revenueStreams.where((s) => s.isNotEmpty).toList(),
        "target_market": _targetMarketController.text.trim(),
        "competitive_advantage": _competitiveAdvantageController.text.trim(),
        "implementation_plan": _implementationPlanController.text.trim(),
        "technology_stack": _technologyStack.where((s) => s.isNotEmpty).toList(),
        "development_timeline": _developmentTimelineController.text.trim(),
        "go_to_market_strategy": _goToMarketController.text.trim(),
        "initial_investment": double.tryParse(_initialInvestmentController.text) ?? 0.0,
        "monthly_burn_rate": double.tryParse(_monthlyBurnRateController.text) ?? 0.0,
        "projected_revenue": double.tryParse(_projectedRevenueController.text) ?? 0.0,
        "funding_requirement": double.tryParse(_fundingRequirementController.text) ?? 0.0,
        "founder_details": [
          {
            "name": _founderNameController.text.trim(),
            "email": _founderEmailController.text.trim(),
            "role": _founderRoleController.text.trim(),
            "linkedin_url": _founderLinkedInController.text.trim().isEmpty ? null : _founderLinkedInController.text.trim(),
            "education": _educations.map((education) => {
              'degree': (education['degree'] as TextEditingController).text.trim(),
              'institution': (education['institution'] as TextEditingController).text.trim(),
              'graduation_year': (education['graduation_year'] as TextEditingController).text.trim(),
            }).where((e) => e['degree']!.isNotEmpty || e['institution']!.isNotEmpty).toList(),
            "experience": _experiences.map((exp) => {
              "company": (exp['company'] as TextEditingController).text.trim(),
              "position": (exp['position'] as TextEditingController).text.trim(), 
              "start_date": _convertToISODate((exp['start_date'] as TextEditingController).text.trim()),
              "end_date": _convertToISODate((exp['end_date'] as TextEditingController).text.trim()),
              "description": (exp['description'] as TextEditingController).text.trim(),
            }).where((exp) => exp['company']!.isNotEmpty || exp['description']!.isNotEmpty).toList(),
            "skills": _founderSkillsController.text.trim().split(',').map((s) => s.trim()).where((s) => s.isNotEmpty).toList(),
            "achievements": _founderAchievementsController.text.trim().split('\n').where((e) => e.trim().isNotEmpty).toList(),
          }
        ]
      };

      // Remove null values from payload
      payload.removeWhere((key, value) => value == null);

      // Validate required fields
      final founderDetails = (payload["founder_details"] as List)[0];
      final requiredFields = {
        'name': payload["name"],
        'description': payload["description"],
        'industry': payload["industry"],
        'sector': payload["sector"],
        'founder_name': founderDetails["name"],
        'founder_email': founderDetails["email"],
        'founder_role': founderDetails["role"],
      };
      
      final emptyFields = <String>[];
      requiredFields.forEach((key, value) {
        if (value == null || value.toString().trim().isEmpty) {
          emptyFields.add(key);
        }
      });
      
      if (emptyFields.isNotEmpty) {
        print('❌ Empty required fields: $emptyFields');
        setState(() {
          _errorMessage = 'Required fields are empty: ${emptyFields.join(", ")}';
          _isLoading = false;
        });
        return;
      }
      
      // Validate numeric fields
      final numericFields = ['initial_investment', 'monthly_burn_rate', 'projected_revenue', 'funding_requirement'];
      for (final field in numericFields) {
        final value = payload[field];
        if (value != null && value is! num) {
          print('❌ $field is not a number: $value (${value.runtimeType})');
          setState(() {
            _errorMessage = 'Invalid numeric value for $field';
            _isLoading = false;
          });
          return;
        }
      }

      // Print the dummy payload
      print('=== DUMMY STARTUP DATA ===');
      print(json.encode(payload));
      print('=== END DUMMY DATA ===');
      
      // Simulate server processing delay
      await Future.delayed(Duration(milliseconds: 1500));
      
      print('✅ Dummy onboarding completed successfully!');
      
      // Navigate to main application
      if (mounted) {
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => MainNavigation(
              userName: widget.userName,
              userEmail: widget.userEmail,
            ),
          ),
        );
      }
    } catch (e) {
      print('Dummy onboarding error: $e');
      setState(() {
        _errorMessage = 'Unable to process your information: ${e.toString()}';
      });
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  void _addExperience() {
    setState(() {
      _experiences.add({
        'company': TextEditingController(),
        'position': TextEditingController(),
        'start_date': TextEditingController(),
        'end_date': TextEditingController(),
        'description': TextEditingController(),
      });
    });
  }

  void _removeExperience(int index) {
    if (_experiences.length > 1) {
      setState(() {
        // Dispose controllers before removing
        ((_experiences[index]['company']) as TextEditingController).dispose();
        ((_experiences[index]['position']) as TextEditingController).dispose();
        ((_experiences[index]['start_date']) as TextEditingController).dispose();
        ((_experiences[index]['end_date']) as TextEditingController).dispose();
        ((_experiences[index]['description']) as TextEditingController).dispose();
        _experiences.removeAt(index);
      });
    }
  }

  Widget _buildExperienceCard(int index) {
    final experience = _experiences[index];
    return Card(
      margin: EdgeInsets.only(bottom: 16),
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    'Experience ${index + 1}',
                    style: AppTheme.labelTextStyle.copyWith(
                      fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                      color: Colors.black87,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                if (_experiences.length > 1)
                  IconButton(
                    onPressed: () => _removeExperience(index),
                    icon: Icon(Icons.delete_outline, color: Colors.red[400]),
                    iconSize: 20,
                  ),
              ],
            ),
            SizedBox(height: 12),
            
            _buildInputField(
              label: 'Company',
              controller: experience['company'] as TextEditingController,
              hintText: 'e.g., Google, Microsoft',
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter company name';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Position',
              controller: experience['position'] as TextEditingController,
              hintText: 'e.g., Senior Software Engineer',
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter position';
                }
                return null;
              },
            ),
            
            Row(
              children: [
                Expanded(
                  child: _buildDateField(
                    label: 'Start Date',
                    controller: experience['start_date'] as TextEditingController,
                    hintText: 'Select start date',
                    validator: (value) {
                      if (index == 0 && (value == null || value.isEmpty)) {
                        return 'Please select start date';
                      }
                      return null;
                    },
                  ),
                ),
                SizedBox(width: 16),
                Expanded(
                  child: _buildDateField(
                    label: 'End Date',
                    controller: experience['end_date'] as TextEditingController,
                    hintText: 'Select end date (leave blank if current)',
                    validator: (value) {
                      // Optional field
                      return null;
                    },
                    allowEmpty: true,
                  ),
                ),
              ],
            ),
            
            _buildInputField(
              label: 'Description',
              controller: experience['description'] as TextEditingController,
              hintText: 'Describe your role and achievements',
              maxLines: 3,
              minLines: 2,
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter job description';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildDateField({
    required String label,
    required TextEditingController controller,
    required String? Function(String?) validator,
    String? hintText,
    bool allowEmpty = false,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
            color: Colors.black87,
          ),
        ),
        SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
        TextFormField(
          controller: controller,
          validator: validator,
          readOnly: true,
          style: AppTheme.inputTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
            color: Colors.black87,
          ),
          decoration: InputDecoration(
            hintText: hintText,
            hintStyle: TextStyle(
              color: Colors.grey[600],
              fontSize: ResponsiveUtils.getBodySize(context) * 0.85,
            ),
            contentPadding: EdgeInsets.symmetric(
              horizontal: 20,
              vertical: ResponsiveUtils.isSmallPhone(context) ? 14 : 17,
            ),
            filled: true,
            fillColor: Colors.white,
            suffixIcon: IconButton(
              icon: Icon(Icons.calendar_today, color: AppTheme.authAccentColor, size: 20),
              onPressed: () => _selectDate(controller),
              padding: EdgeInsets.zero,
              constraints: BoxConstraints(minWidth: 48, minHeight: 48),
            ),
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightErrorColor.withOpacity(0.5), width: 1),
            ),
            focusedErrorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightErrorColor, width: 1.5),
            ),
          ),
        ),
        SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
      ],
    );
  }

  Future<void> _selectDate(TextEditingController controller) async {
    final DateTime? picked = await showDatePicker(
      context: context,
      initialDate: DateTime.now(),
      firstDate: DateTime(1990),
      lastDate: DateTime.now(),
      builder: (context, child) {
        return Theme(
          data: Theme.of(context).copyWith(
            colorScheme: ColorScheme.light(
              primary: AppTheme.authAccentColor,
              onPrimary: Colors.white,
              surface: Colors.white,
              onSurface: Colors.black87,
            ),
          ),
          child: child!,
        );
      },
    );
    
    if (picked != null) {
      setState(() {
        controller.text = '${picked.day.toString().padLeft(2, '0')}-${picked.month.toString().padLeft(2, '0')}-${picked.year.toString().substring(2)}';
      });
    }
  }

  String _convertToISODate(String dateString) {
    if (dateString.isEmpty) return '';
    
    try {
      // Parse DD-MM-YY format
      final parts = dateString.split('-');
      if (parts.length == 3) {
        final day = parts[0];
        final month = parts[1];
        final yearShort = parts[2];
        
        // Convert YY to YYYY (assuming years 00-30 are 2000-2030, 31-99 are 1931-1999)
        final year = int.parse(yearShort) <= 30 
            ? '20$yearShort' 
            : '19$yearShort';
            
        return '$year-$month-$day'; // Convert to YYYY-MM-DD
      }
    } catch (e) {
      // If parsing fails, return the original string
    }
    return dateString;
  }

  // Helper method to build input field
  Widget _buildInputField({
    required String label,
    required TextEditingController controller,
    required String? Function(String?) validator,
    String? hintText,
    int maxLines = 1,
    int? minLines,
    TextInputType? keyboardType,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          label,
          style: AppTheme.labelTextStyle.copyWith(
            fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
            color: Colors.black87,
          ),
        ),
        SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
        TextFormField(
          controller: controller,
          validator: validator,
          maxLines: maxLines,
          minLines: minLines,
          keyboardType: keyboardType,
          style: AppTheme.inputTextStyle.copyWith(
            fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
            color: Colors.black87,
          ),
          decoration: InputDecoration(
            hintText: hintText,
            hintStyle: TextStyle(
              color: Colors.grey[600],
              fontSize: ResponsiveUtils.getBodySize(context) * 0.85,
            ),
            contentPadding: EdgeInsets.symmetric(
              horizontal: 20,
              vertical: ResponsiveUtils.isSmallPhone(context) ? 14 : 17,
            ),
            filled: true,
            fillColor: Colors.white,
            border: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
            ),
            enabledBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
            ),
            focusedBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
            ),
            errorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightErrorColor.withOpacity(0.5), width: 1),
            ),
            focusedErrorBorder: OutlineInputBorder(
              borderRadius: BorderRadius.circular(30),
              borderSide: BorderSide(color: AppTheme.lightErrorColor, width: 1.5),
            ),
          ),
        ),
        SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
      ],
    );
  }

  Widget _buildStep1() {
    return SingleChildScrollView(
      padding: EdgeInsets.all(24),
      child: Form(
        key: _formKeys[0],
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Company Information',
              style: AppTheme.headingTextStyle.copyWith(
                fontSize: ResponsiveUtils.getHeadingSize(context) * 0.8,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: 24),
            
            _buildInputField(
              label: 'Company Name',
              controller: _companyNameController,
              hintText: 'Enter your company name',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter company name';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Description',
              controller: _descriptionController,
              hintText: 'Describe your business in detail',
              maxLines: 4,
              minLines: 3,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please provide a description';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Industry',
              controller: _industryController,
              hintText: 'e.g., Technology, Healthcare, Finance',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please specify your industry';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Sector',
              controller: _sectorController,
              hintText: 'e.g., FinTech, HealthTech, EdTech',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please specify your sector';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Location',
              controller: _locationController,
              hintText: 'e.g., San Francisco, CA',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please specify your location';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Website',
              controller: _websiteController,
              hintText: 'https://yourcompany.com',
              keyboardType: TextInputType.url,
              validator: (value) {
                if (value != null && value.isNotEmpty) {
                  if (!value.startsWith('http://') && !value.startsWith('https://')) {
                    return 'Please enter a valid URL';
                  }
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildStep2() {
    return SingleChildScrollView(
      padding: EdgeInsets.all(24),
      child: Form(
        key: _formKeys[1],
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Business Model & Market',
              style: AppTheme.headingTextStyle.copyWith(
                fontSize: ResponsiveUtils.getHeadingSize(context) * 0.8,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: 24),
            
            _buildInputField(
              label: 'Business Model',
              controller: _businessModelController,
              hintText: 'e.g., B2B SaaS, B2C Marketplace, Freemium',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please specify your business model';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Target Market',
              controller: _targetMarketController,
              hintText: 'Describe your target customers',
              maxLines: 3,
              minLines: 2,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please describe your target market';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Competitive Advantage',
              controller: _competitiveAdvantageController,
              hintText: 'What makes you unique?',
              maxLines: 3,
              minLines: 2,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please describe your competitive advantage';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Go-to-Market Strategy',
              controller: _goToMarketController,
              hintText: 'How will you acquire customers?',
              maxLines: 3,
              minLines: 2,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please describe your go-to-market strategy';
                }
                return null;
              },
            ),
            
            // Funding Stage Dropdown
            Text(
              'Funding Stage',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
            DropdownButtonFormField<String>(
              value: _fundingStage,
              decoration: InputDecoration(
                contentPadding: EdgeInsets.symmetric(horizontal: 20, vertical: 17),
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
                ),
              ),
              items: ['Pre-Seed', 'Seed', 'Series A', 'Series B', 'Series C+'].map((String value) {
                return DropdownMenuItem<String>(
                  value: value,
                  child: Text(value),
                );
              }).toList(),
              onChanged: (String? newValue) {
                setState(() {
                  _fundingStage = newValue!;
                });
              },
            ),
            SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
            
            // Team Size
            Text(
              'Team Size',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
            Row(
              children: [
                Expanded(
                  child: Slider(
                    value: _teamSize.toDouble(),
                    min: 1,
                    max: 50,
                    divisions: 49,
                    activeColor: AppTheme.authAccentColor,
                    inactiveColor: AppTheme.lightDividerColor,
                    thumbColor: AppTheme.authAccentColor,
                    label: _teamSize.toString(),
                    onChanged: (double value) {
                      setState(() {
                        _teamSize = value.round();
                      });
                    },
                  ),
                ),
                Container(
                  width: 60,
                  padding: EdgeInsets.symmetric(horizontal: 12, vertical: 8),
                  decoration: BoxDecoration(
                    border: Border.all(color: AppTheme.lightDividerColor),
                    borderRadius: BorderRadius.circular(8),
                  ),
                  child: Text(
                    _teamSize.toString(),
                    textAlign: TextAlign.center,
                    style: TextStyle(
                      fontSize: ResponsiveUtils.getBodySize(context),
                      fontWeight: FontWeight.w500,
                    ),
                  ),
                ),
              ],
            ),
            SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
            
            // Revenue Streams (simplified as text input)
            Text(
              'Revenue Streams',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
            TextFormField(
              maxLines: 2,
              style: AppTheme.inputTextStyle.copyWith(
                fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
                color: Colors.black87,
              ),
              decoration: InputDecoration(
                hintText: 'e.g., Monthly Subscription, Enterprise Licenses (comma-separated)',
                hintStyle: TextStyle(
                  color: Colors.grey[600],
                  fontSize: ResponsiveUtils.getBodySize(context) * 0.85,
                ),
                contentPadding: EdgeInsets.symmetric(horizontal: 20, vertical: 17),
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
                ),
              ),
              onChanged: (value) {
                _revenueStreams = value.split(',').map((s) => s.trim()).where((s) => s.isNotEmpty).toList();
              },
            ),
            SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
          ],
        ),
      ),
    );
  }

  Widget _buildStep3() {
    return SingleChildScrollView(
      padding: EdgeInsets.all(24),
      child: Form(
        key: _formKeys[2],
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Technical & Implementation',
              style: AppTheme.headingTextStyle.copyWith(
                fontSize: ResponsiveUtils.getHeadingSize(context) * 0.8,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: 24),
            
            _buildInputField(
              label: 'Implementation Plan',
              controller: _implementationPlanController,
              hintText: 'Describe your development phases and timeline',
              maxLines: 4,
              minLines: 3,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please describe your implementation plan';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Development Timeline',
              controller: _developmentTimelineController,
              hintText: 'e.g., 12 months to full market release',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please specify your development timeline';
                }
                return null;
              },
            ),
            
            // Technology Stack (simplified as text input)
            Text(
              'Technology Stack',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.6),
            TextFormField(
              maxLines: 3,
              minLines: 2,
              style: AppTheme.inputTextStyle.copyWith(
                fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
                color: Colors.black87,
              ),
              decoration: InputDecoration(
                hintText: 'e.g., React, Node.js, PostgreSQL, Docker (comma-separated)',
                hintStyle: TextStyle(
                  color: Colors.grey[600],
                  fontSize: ResponsiveUtils.getBodySize(context) * 0.85,
                ),
                contentPadding: EdgeInsets.symmetric(horizontal: 20, vertical: 17),
                filled: true,
                fillColor: Colors.white,
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                enabledBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.lightDividerColor, width: 1),
                ),
                focusedBorder: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(30),
                  borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
                ),
              ),
              onChanged: (value) {
                _technologyStack = value.split(',').map((s) => s.trim()).where((s) => s.isNotEmpty).toList();
              },
            ),
            SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
          ],
        ),
      ),
    );
  }

  Widget _buildStep4() {
    return SingleChildScrollView(
      padding: EdgeInsets.all(24),
      child: Form(
        key: _formKeys[3],
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              'Financial & Founder Information',
              style: AppTheme.headingTextStyle.copyWith(
                fontSize: ResponsiveUtils.getHeadingSize(context) * 0.8,
                color: Colors.black87,
              ),
            ),
            SizedBox(height: 24),
            
            // Financial Information
            Text(
              'Financial Information',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                color: Colors.black87,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            
            _buildInputField(
              label: 'Initial Investment (\$)',
              controller: _initialInvestmentController,
              hintText: 'e.g., 500000',
              keyboardType: TextInputType.number,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter initial investment amount';
                }
                if (double.tryParse(value) == null) {
                  return 'Please enter a valid number';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Monthly Burn Rate (\$)',
              controller: _monthlyBurnRateController,
              hintText: 'e.g., 45000',
              keyboardType: TextInputType.number,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter monthly burn rate';
                }
                if (double.tryParse(value) == null) {
                  return 'Please enter a valid number';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Projected Annual Revenue (\$)',
              controller: _projectedRevenueController,
              hintText: 'e.g., 1200000',
              keyboardType: TextInputType.number,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter projected revenue';
                }
                if (double.tryParse(value) == null) {
                  return 'Please enter a valid number';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Funding Requirement (\$)',
              controller: _fundingRequirementController,
              hintText: 'e.g., 2000000',
              keyboardType: TextInputType.number,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter funding requirement';
                }
                if (double.tryParse(value) == null) {
                  return 'Please enter a valid number';
                }
                return null;
              },
            ),
            
            SizedBox(height: 24),
            
            // Founder Information
            Text(
              'Founder Information',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                color: Colors.black87,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            
            _buildInputField(
              label: 'Founder Name',
              controller: _founderNameController,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter founder name';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Founder Email',
              controller: _founderEmailController,
              keyboardType: TextInputType.emailAddress,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter founder email';
                }
                if (!value.contains('@')) {
                  return 'Please enter a valid email';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Founder Role',
              controller: _founderRoleController,
              hintText: 'e.g., CEO & Founder',
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter founder role';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'LinkedIn URL',
              controller: _founderLinkedInController,
              hintText: 'https://linkedin.com/in/yourprofile',
              keyboardType: TextInputType.url,
              validator: (value) {
                // Optional field
                return null;
              },
            ),
            
            // Education section
            Text(
              'Education',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                color: Colors.black87,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            
            // Education cards
            ...List.generate(_educations.length, (index) => _buildEducationCard(index)),
            
            // Add Education button
            Container(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: _addEducation,
                icon: Icon(Icons.add, color: AppTheme.authAccentColor),
                label: Text(
                  'Add Another Education',
                  style: TextStyle(color: AppTheme.authAccentColor),
                ),
                style: OutlinedButton.styleFrom(
                  side: BorderSide(color: AppTheme.authAccentColor),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(30),
                  ),
                  padding: EdgeInsets.symmetric(vertical: 12),
                ),
              ),
            ),
            SizedBox(height: 24),
            
            // Professional Experience Section
            Text(
              'Professional Experience',
              style: AppTheme.labelTextStyle.copyWith(
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                color: Colors.black87,
                fontWeight: FontWeight.w600,
              ),
            ),
            SizedBox(height: 16),
            
            // Experience cards
            ...List.generate(_experiences.length, (index) => _buildExperienceCard(index)),
            
            // Add Experience button
            Container(
              width: double.infinity,
              child: OutlinedButton.icon(
                onPressed: _addExperience,
                icon: Icon(Icons.add, color: AppTheme.authAccentColor),
                label: Text(
                  'Add Another Experience',
                  style: TextStyle(color: AppTheme.authAccentColor),
                ),
                style: OutlinedButton.styleFrom(
                  side: BorderSide(color: AppTheme.authAccentColor),
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(30),
                  ),
                  padding: EdgeInsets.symmetric(vertical: 12),
                ),
              ),
            ),
            SizedBox(height: 24),
            
            _buildInputField(
              label: 'Skills',
              controller: _founderSkillsController,
              hintText: 'e.g., AI/ML, Product Strategy, Team Leadership (comma-separated)',
              maxLines: 2,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter your skills';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Key Achievements',
              controller: _founderAchievementsController,
              hintText: 'Enter each achievement on a new line',
              maxLines: 4,
              minLines: 3,
              validator: (value) {
                if (value == null || value.isEmpty) {
                  return 'Please enter key achievements';
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildEducationCard(int index) {
    final education = _educations[index];
    return Card(
      margin: EdgeInsets.only(bottom: 16),
      elevation: 2,
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
      child: Padding(
        padding: EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    'Education ${index + 1}',
                    style: AppTheme.labelTextStyle.copyWith(
                      fontSize: ResponsiveUtils.getSmallTextSize(context) + 2,
                      color: Colors.black87,
                      fontWeight: FontWeight.w600,
                    ),
                  ),
                ),
                if (_educations.length > 1)
                  IconButton(
                    onPressed: () => _removeEducation(index),
                    icon: Icon(Icons.delete_outline, color: Colors.red[400]),
                    iconSize: 20,
                  ),
              ],
            ),
            SizedBox(height: 12),
            
            _buildInputField(
              label: 'Degree/Qualification',
              controller: education['degree'] as TextEditingController,
              hintText: 'e.g., Bachelor of Computer Science',
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter degree/qualification';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Institution',
              controller: education['institution'] as TextEditingController,
              hintText: 'e.g., MIT, Stanford University',
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter institution';
                }
                return null;
              },
            ),
            
            _buildInputField(
              label: 'Graduation Year',
              controller: education['graduation_year'] as TextEditingController,
              hintText: 'e.g., 2020',
              keyboardType: TextInputType.number,
              validator: (value) {
                if (index == 0 && (value == null || value.isEmpty)) {
                  return 'Please enter graduation year';
                }
                if (value != null && value.isNotEmpty) {
                  final year = int.tryParse(value);
                  if (year == null || year < 1950 || year > DateTime.now().year + 10) {
                    return 'Please enter a valid year';
                  }
                }
                return null;
              },
            ),
          ],
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.transparent,
      extendBodyBehindAppBar: true,
      resizeToAvoidBottomInset: false,
      body: Stack(
        children: [
          // Video background
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            child: _videoController != null && _videoController!.value.isInitialized && !_videoFailed
                ? AspectRatio(
                    aspectRatio: _videoController!.value.aspectRatio,
                    child: VideoPlayer(_videoController!),
                  )
                : Container(
                    height: MediaQuery.of(context).size.height * 0.1,
                    decoration: BoxDecoration(
                      gradient: LinearGradient(
                        begin: Alignment.topLeft,
                        end: Alignment.bottomRight,
                        colors: [
                          AppTheme.authSecondaryColor,
                          AppTheme.authPrimaryColor,
                          AppTheme.authTertiaryColor,
                          AppTheme.authDarkBlue,
                        ],
                        stops: [0.0, 0.3, 0.7, 1.0],
                      ),
                    ),
                  ),
          ),
          
          // Opacity overlay
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            height: _videoController != null && _videoController!.value.isInitialized && !_videoFailed
                ? MediaQuery.of(context).size.width / _videoController!.value.aspectRatio
                : MediaQuery.of(context).size.height * 0.2,
            child: Container(
              color: Colors.black.withOpacity(0.2),
            ),
          ),
          
          Column(
            children: [
              // Header with progress
              SafeArea(
                child: Padding(
                  padding: EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        'Welcome ${widget.userName}!',
                        style: AppTheme.headingTextStyle.copyWith(
                          color: Colors.white,
                          fontSize: ResponsiveUtils.getHeadingSize(context),
                        ),
                      ),
                      SizedBox(height: 8),
                      Text(
                        'Let\'s set up your startup profile',
                        style: AppTheme.regularTextStyle.copyWith(
                          color: Colors.white70,
                          fontSize: ResponsiveUtils.getBodySize(context),
                        ),
                      ),
                      SizedBox(height: 16),
                      // Progress indicator
                      LinearProgressIndicator(
                        value: (_currentStep + 1) / _totalSteps,
                        backgroundColor: Colors.white.withOpacity(0.3),
                        valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                      ),
                      SizedBox(height: 8),
                      Text(
                        'Step ${_currentStep + 1} of $_totalSteps',
                        style: AppTheme.regularTextStyle.copyWith(
                          color: Colors.white70,
                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
              
              // Form content
              Expanded(
                child: Container(
                  width: double.infinity,
                  decoration: BoxDecoration(
                    color: Colors.white,
                    borderRadius: BorderRadius.only(
                      topLeft: Radius.circular(30),
                      topRight: Radius.circular(30),
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.1),
                        blurRadius: 10,
                        offset: const Offset(0, -5),
                      ),
                    ],
                  ),
                  child: Column(
                    children: [
                      Expanded(
                        child: PageView(
                          controller: _pageController,
                          physics: NeverScrollableScrollPhysics(),
                          children: [
                            _buildStep1(),
                            _buildStep2(),
                            _buildStep3(),
                            _buildStep4(),
                          ],
                        ),
                      ),
                      
                      // Navigation buttons
                      Container(
                        padding: EdgeInsets.all(24),
                        child: Row(
                          children: [
                            if (_currentStep > 0)
                              Expanded(
                                child: OutlinedButton(
                                  onPressed: _previousStep,
                                  style: OutlinedButton.styleFrom(
                                    side: BorderSide(color: AppTheme.authAccentColor),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(30),
                                    ),
                                    padding: EdgeInsets.symmetric(vertical: 16),
                                  ),
                                  child: Text(
                                    'Previous',
                                    style: AppTheme.buttonTextStyle.copyWith(
                                      color: AppTheme.authAccentColor,
                                    ),
                                  ),
                                ),
                              ),
                            if (_currentStep > 0) SizedBox(width: 16),
                            Expanded(
                              flex: _currentStep == 0 ? 1 : 1,
                              child: Container(
                                decoration: BoxDecoration(
                                  gradient: LinearGradient(
                                    begin: _gradientBegin,
                                    end: _gradientEnd,
                                    colors: [
                                      AppTheme.authPrimaryColor,
                                      AppTheme.authSecondaryColor,
                                      AppTheme.authTertiaryColor,
                                    ],
                                  ),
                                  borderRadius: BorderRadius.circular(30),
                                  boxShadow: [
                                    BoxShadow(
                                      color: AppTheme.authAccentColor.withOpacity(0.3),
                                      blurRadius: 8,
                                      offset: Offset(0, 4),
                                    ),
                                  ],
                                ),
                                child: ElevatedButton(
                                  onPressed: _isLoading ? null : _nextStep,
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: Colors.transparent,
                                    shadowColor: Colors.transparent,
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(30),
                                    ),
                                    padding: EdgeInsets.symmetric(vertical: 16),
                                  ),
                                  child: _isLoading
                                      ? SizedBox(
                                          height: 20,
                                          width: 20,
                                          child: CircularProgressIndicator(
                                            strokeWidth: 2,
                                            valueColor: AlwaysStoppedAnimation<Color>(Colors.white),
                                          ),
                                        )
                                      : Text(
                                          _currentStep == _totalSteps - 1 ? 'Complete Setup' : 'Next',
                                          style: AppTheme.buttonTextStyle.copyWith(
                                            color: Colors.white,
                                          ),
                                        ),
                                ),
                              ),
                            ),
                          ],
                        ),
                      ),
                      
                      if (_errorMessage.isNotEmpty)
                        Padding(
                          padding: EdgeInsets.fromLTRB(24, 0, 24, 24),
                          child: Text(
                            _errorMessage,
                            style: TextStyle(
                              color: AppTheme.lightErrorColor,
                              fontSize: ResponsiveUtils.getSmallTextSize(context),
                            ),
                          ),
                        ),
                    ],
                  ),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}

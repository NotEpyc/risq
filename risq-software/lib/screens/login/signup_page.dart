import 'package:flutter/material.dart';
import 'package:video_player/video_player.dart';
import 'package:risq/theme/theme.dart';
import 'dart:math' as math;
import 'signin_page.dart';
import '../pages/onboarding_page.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/services/storage_service.dart';

class SignupPage extends StatefulWidget {
  const SignupPage({super.key});

  @override
  _SignupPageState createState() => _SignupPageState();
}

class _SignupPageState extends State<SignupPage> with TickerProviderStateMixin {
  // Add controllers
  final _nameController = TextEditingController();
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  final _confirmPasswordController = TextEditingController();
  
  // Add state variables
  bool _agreeToTerms = false;
  bool _isLoading = false;
  String _errorMessage = '';
  bool _passwordVisible = false;
  bool _confirmPasswordVisible = false;
  
  // Animation controllers for slide-up effect
  AnimationController? _animationController;
  Animation<double>? _containerAnimation;
  Animation<double>? _fadeAnimation;
  
  // Video controller for background
  VideoPlayerController? _videoController;
  bool _videoFailed = false;
  
  // Random gradient direction
  late Alignment _gradientBegin;
  late Alignment _gradientEnd;
  
  @override
  void initState() {
    super.initState();
    
    // Initialize random gradient direction
    _initializeRandomGradient();
    
    // Initialize video controller
    _videoController = VideoPlayerController.asset('assets/videos/bg_video.mp4')
      ..initialize().then((_) {
        print('Video initialized successfully');
        if (mounted) {
          _videoController!.setLooping(true);
          _videoController!.setVolume(0.0); // Mute the video
          _videoController!.play();
          print('Video started playing');
          setState(() {});
        }
      }).catchError((error) {
        print('Video initialization error: $error');
        if (mounted) {
          setState(() {
            _videoFailed = true;
          });
        }
      });
    
    // Initialize animation controller
    _animationController = AnimationController(
      duration: const Duration(milliseconds: 800),
      vsync: this,
    );
    
    // Create animation for the bottom container
    _containerAnimation = Tween<double>(begin: 1.0, end: 0.0).animate(
      CurvedAnimation(
        parent: _animationController!,
        curve: Curves.easeOut,
      ),
    );
    
    // Create fade-in animation for form elements
    _fadeAnimation = Tween<double>(begin: 0.0, end: 1.0).animate(
      CurvedAnimation(
        parent: _animationController!,
        curve: Interval(0.4, 1.0, curve: Curves.easeIn),
      ),
    );
    
    // Start animation after a short delay
    Future.delayed(Duration(milliseconds: 100), () {
      if (mounted) {
        _animationController!.forward();
      }    });
  }

  void _initializeRandomGradient() {
    final random = math.Random();
    
    // Generate random angle between 0 and 2π
    final angle = random.nextDouble() * 2 * math.pi;
    
    // Convert angle to alignment coordinates
    final x = math.cos(angle);
    final y = math.sin(angle);
    
    _gradientBegin = Alignment(x, y);
    _gradientEnd = Alignment(-x, -y);
  }

  @override
  void dispose() {
    _animationController?.dispose();
    _videoController?.dispose();
    _nameController.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    _confirmPasswordController.dispose();
    super.dispose();
  }
  
  // Validate form
  bool _validateForm() {
    if (_nameController.text.trim().isEmpty) {
      setState(() => _errorMessage = 'Please enter your name');
      return false;
    }
    if (_emailController.text.trim().isEmpty) {
      setState(() => _errorMessage = 'Please enter your email');
      return false;
    }
    if (_passwordController.text.isEmpty) {
      setState(() => _errorMessage = 'Please enter a password');
      return false;
    }
    if (_passwordController.text != _confirmPasswordController.text) {
      setState(() => _errorMessage = 'Passwords do not match');
      return false;
    }
    if (!_agreeToTerms) {
      setState(() => _errorMessage = 'Please agree to the Terms of Service');
      return false;
    }
    return true;
  }
  
  Future<void> _signUpWithEmailAndPassword() async {
    // Clear any previous errors
    setState(() => _errorMessage = '');
    
    // Validate form
    if (!_validateForm()) {
      print('SignupPage: Form validation failed');
      return;
    }
    
    setState(() => _isLoading = true);
    
    try {
      print('SignupPage: Starting signup process...');
      print('SignupPage: Name: "${_nameController.text.trim()}"');
      print('SignupPage: Email: "${_emailController.text.trim()}"');
      print('SignupPage: Password length: ${_passwordController.text.length}');
      
      // Check if fields are actually filled
      final name = _nameController.text.trim();
      final email = _emailController.text.trim();
      final password = _passwordController.text;
      
      if (name.isEmpty || email.isEmpty || password.isEmpty) {
        print('SignupPage: Error - required fields are empty after trim');
        setState(() {
          _errorMessage = 'Please fill in all required fields.';
        });
        return;
      }
      
      print('SignupPage: Calling StorageService.saveUserData...');
      
      // Save user data to local storage
      final success = await StorageService.saveUserData(
        name: name,
        email: email,
        password: password,
      );
      
      print('SignupPage: Storage result: $success');
      
      if (success) {
        print('SignupPage: Data saved successfully, navigating to onboarding...');
        
        // Simulate a brief delay for user feedback
        await Future.delayed(Duration(milliseconds: 800));
        
        // Navigate to onboarding page
        if (mounted) {
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(
              builder: (context) => OnboardingPage(
                userName: name,
                userEmail: email,
              ),
            ),
          );
        }
      } else {
        print('SignupPage: Failed to save data to storage');
        setState(() {
          _errorMessage = 'Failed to save account data. Please try again.';
        });
      }
    } catch (e) {
      print('SignupPage: Exception during signup: $e');
      if (e is Error) {
        print('SignupPage: Stack trace: ${e.stackTrace}');
      }
      setState(() {
        _errorMessage = 'Unable to create account. Please try again.';
      });
    } finally {
      if (mounted) {
        setState(() => _isLoading = false);
      }
    }
  }

  @override
  Widget build(BuildContext context) {    
    return Scaffold(
      backgroundColor: Colors.transparent,
      extendBodyBehindAppBar: true,
      resizeToAvoidBottomInset: true,
      body: Stack(
        fit: StackFit.expand,
        children: [
          // Video background - positioned at the top of the screen
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
                    height: MediaQuery.of(context).size.height * 0.5,
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
          
          // Opacity overlay - only over the video area
          Positioned(
            top: 0,
            left: 0,
            right: 0,
            height: _videoController != null && _videoController!.value.isInitialized && !_videoFailed
                ? MediaQuery.of(context).size.width / _videoController!.value.aspectRatio
                : MediaQuery.of(context).size.height * 0.5,
            child: Container(
              color: Colors.black.withOpacity(0.2),
            ),
          ),
          
          // Main content with animation
          SafeArea(
            child: LayoutBuilder(
              builder: (context, constraints) {
                // Increase bottom container height for better spacing
                final bottomContainerHeight = constraints.maxHeight * 
                    (ResponsiveUtils.getBottomContainerRatio(context) + 0.03);
                
                return Stack(
                  children: [
                    // Bottom container with animation
                    AnimatedBuilder(
                      animation: _animationController!,
                      builder: (context, child) {
                        return Positioned(
                          bottom: -bottomContainerHeight * _containerAnimation!.value,
                          left: 0,
                          right: 0,
                          height: bottomContainerHeight,
                          child: child!,
                        );
                      },
                      child: Container(
                        padding: EdgeInsets.all(ResponsiveUtils.getMediumSpace(context)),
                        width: double.infinity,
                        constraints: BoxConstraints(
                          maxWidth: ResponsiveUtils.isDesktop(context) ? 450 : double.infinity,
                        ),
                        decoration: BoxDecoration(
                          color: AppTheme.authBackgroundColor,
                          borderRadius: BorderRadius.only(
                            topLeft: Radius.circular(20),
                            topRight: Radius.circular(20),
                          ),
                          boxShadow: [
                            BoxShadow(
                              color: Colors.black.withOpacity(0.1),
                              spreadRadius: 1,
                              blurRadius: 10,
                              offset: Offset(0, -2),
                            ),
                          ],
                        ),
                        child: SingleChildScrollView(
                          physics: ClampingScrollPhysics(),
                          child: AnimatedBuilder(
                            animation: _fadeAnimation!,
                            builder: (context, child) {
                              return Opacity(
                                opacity: _fadeAnimation!.value,
                                child: child,
                              );
                            },
                            child: Column(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context)),
                                Text(
                                  "Sign Up",
                                  style: TextStyle(
                                    fontSize: ResponsiveUtils.getHeadingSize(context),
                                    fontWeight: FontWeight.bold,
                                    color: AppTheme.authPrimaryColor,
                                  ),
                                ),
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                _buildTextField(
                                  label: 'Founder Name',
                                  icon: Icons.person,
                                  controller: _nameController,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                _buildTextField(
                                  label: 'Email',
                                  icon: Icons.email,
                                  controller: _emailController,
                                  keyboardType: TextInputType.emailAddress,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                _buildTextField(
                                  label: 'Password',
                                  icon: Icons.lock,
                                  isPassword: true,
                                  controller: _passwordController,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                _buildTextField(
                                  label: 'Confirm Password',
                                  icon: Icons.lock,
                                  isPassword: true,
                                  isConfirmPassword: true,
                                  controller: _confirmPasswordController,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                
                                // Terms and conditions checkbox with proper alignment
                                Row(
                                  crossAxisAlignment: CrossAxisAlignment.center,
                                  children: [
                                    Transform.scale(
                                      scale: ResponsiveUtils.isSmallPhone(context) ? 0.9 : 1.0,
                                      child: Checkbox(
                                        value: _agreeToTerms,
                                        onChanged: (value) {
                                          setState(() {
                                            _agreeToTerms = value ?? false;
                                          });
                                        },
                                        activeColor: AppTheme.authAccentColor,
                                        checkColor: Colors.white,
                                        // Add these properties for better visibility
                                        side: BorderSide(
                                          color: AppTheme.authDividerColor.withOpacity(0.8), 
                                          width: 1.5,
                                        ),
                                        shape: RoundedRectangleBorder(
                                          borderRadius: BorderRadius.circular(4),
                                        ),
                                        materialTapTargetSize: ResponsiveUtils.isSmallPhone(context) 
                                            ? MaterialTapTargetSize.shrinkWrap 
                                            : MaterialTapTargetSize.padded,
                                      ),
                                    ),
                                    Expanded(
                                      child: Padding(
                                        // Reduce top padding
                                        padding: EdgeInsets.only(top: 1), // From 2 to 1
                                        child: RichText(
                                          text: TextSpan(
                                            style: TextStyle(
                                              fontSize: ResponsiveUtils.getSmallTextSize(context) - 2,
                                              color: AppTheme.authTextColor.withOpacity(0.8),
                                              height: 1.3, // Add line height to improve readability
                                            ),
                                            children: [
                                              TextSpan(text: 'I agree to the '),
                                              TextSpan(
                                                text: 'Terms of Service',
                                                style: TextStyle(
                                                  color: AppTheme.authAccentColor,
                                                  fontWeight: FontWeight.w500,
                                                  fontSize: ResponsiveUtils.getSmallTextSize(context) - 2,
                                                ),
                                              ),
                                              TextSpan(text: ' and '),
                                              TextSpan(
                                                text: 'Privacy Policy',
                                                style: TextStyle(
                                                  color: AppTheme.authAccentColor,
                                                  fontWeight: FontWeight.w500,
                                                  fontSize: ResponsiveUtils.getSmallTextSize(context) - 2,
                                                ),
                                              ),
                                            ],
                                          ),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                
                                // Error message (if any)
                                if (_errorMessage.isNotEmpty)
                                  Padding(
                                    padding: EdgeInsets.symmetric(
                                        vertical: ResponsiveUtils.getSmallSpace(context)),
                                    child: Text(
                                      _errorMessage,
                                      style: TextStyle(
                                          color: AppTheme.authErrorColor,
                                          fontSize: ResponsiveUtils.getSmallTextSize(context)),
                                      textAlign: TextAlign.center,
                                    ),
                                  ),
                                
                                // Create Account button with randomized static gradient
                                Container(
                                  width: double.infinity,
                                  height: ResponsiveUtils.getButtonHeight(context),
                                  decoration: BoxDecoration(
                                    gradient: LinearGradient(
                                      colors: AppTheme.authButtonGradient,
                                      begin: _gradientBegin,
                                      end: _gradientEnd,
                                    ),
                                    borderRadius: BorderRadius.circular(30),
                                    boxShadow: [
                                      BoxShadow(
                                        color: AppTheme.authButtonGradient[0].withOpacity(0.3),
                                        spreadRadius: 1,
                                        blurRadius: 8,
                                        offset: Offset(0, 3),
                                      ),
                                    ],
                                  ),
                                  child: ElevatedButton(
                                    onPressed: _isLoading ? null : _signUpWithEmailAndPassword,
                                    style: ElevatedButton.styleFrom(
                                      backgroundColor: Colors.transparent,
                                      disabledBackgroundColor: Colors.transparent,
                                      shadowColor: Colors.transparent,
                                      elevation: 0,
                                      shape: RoundedRectangleBorder(
                                        borderRadius: BorderRadius.circular(30),
                                      ),
                                    ),
                                    child: _isLoading 
                                      ? SizedBox(
                                          height: ResponsiveUtils.isSmallPhone(context) ? 18 : 22,
                                          width: ResponsiveUtils.isSmallPhone(context) ? 18 : 22, 
                                          child: CircularProgressIndicator(
                                            color: Colors.white,
                                            strokeWidth: ResponsiveUtils.isSmallPhone(context) ? 1.5 : 2,
                                          ),
                                        )
                                      : Text(
                                          'Create Account', 
                                          style: TextStyle(
                                            fontSize: ResponsiveUtils.getBodySize(context),
                                            fontWeight: FontWeight.bold,
                                            color: Colors.white,
                                          ),
                                        ),
                                  ),
                                ),
                                
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                
                                // Sign In link
                                Container(
                                  margin: EdgeInsets.only(top: ResponsiveUtils.getSmallSpace(context) * 0.5),
                                  padding: EdgeInsets.symmetric(
                                      vertical: ResponsiveUtils.getSmallSpace(context) * 0.5),
                                  decoration: BoxDecoration(
                                    border: Border(
                                      top: BorderSide(
                                        color: AppTheme.authDividerColor,
                                        width: 1,
                                      ),
                                    ),
                                  ),
                                  child: Row(
                                    mainAxisAlignment: MainAxisAlignment.center,
                                    children: [
                                      Text(
                                        "Already have an account? ",
                                        style: TextStyle(
                                          color: AppTheme.authTextColor,
                                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                                        ),
                                      ),
                                      TextButton(
                                        onPressed: () {
                                          Navigator.pushReplacement(
                                            context,
                                            MaterialPageRoute(builder: (context) => SignInPage()),
                                          );
                                        },
                                        style: TextButton.styleFrom(
                                          padding: EdgeInsets.symmetric(
                                              horizontal: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                          minimumSize: Size(0, 0),
                                        ),
                                        child: Text(
                                          "Sign In",
                                          style: TextStyle(
                                            color: AppTheme.authAccentColor,
                                            fontWeight: FontWeight.bold,
                                            fontSize: ResponsiveUtils.getSmallTextSize(context),
                                          ),
                                        ),
                                      ),
                                    ],
                                  ),
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context) * 1.5),
                              ],
                            ),
                          ),
                        ),
                      ),
                    ),
                  ],
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildTextField({
    required String label,
    required IconData icon,
    required TextEditingController controller,
    bool isPassword = false,
    bool isConfirmPassword = false,
    TextInputType keyboardType = TextInputType.text,
  }) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        // Label above the field
        Padding(
          padding: EdgeInsets.only(
            left: 4, 
            bottom: ResponsiveUtils.getSmallSpace(context) * 0.6
          ),
          child: Text(
            label,
            style: TextStyle(
              fontWeight: FontWeight.bold,
              color: AppTheme.authTextColor,
              fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
            ),
          ),
        ),
        // Text field with rounded styling
        Container(
          // Reduce height by about 5-10%
          height: ResponsiveUtils.getInputHeight(context) * 0.95,
          decoration: BoxDecoration(
            color: Colors.white,
            borderRadius: BorderRadius.circular(30),
            boxShadow: [
              BoxShadow(
                color: Colors.black.withOpacity(0.04),
                blurRadius: 4,
                offset: Offset(0, 2),
              ),
            ],
          ),
          child: TextFormField(
            controller: controller,
            obscureText: isPassword && (isConfirmPassword ? !_confirmPasswordVisible : !_passwordVisible),
            keyboardType: keyboardType,
            style: TextStyle(
              fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
              color: AppTheme.authTextColor,
            ),
            decoration: InputDecoration(
              contentPadding: EdgeInsets.symmetric(
                vertical: ResponsiveUtils.isSmallPhone(context) ? 14 : 17,
                horizontal: ResponsiveUtils.isSmallPhone(context) ? 16 : 20,
              ),
              filled: true,
              fillColor: Colors.white,
              hintText: label != 'Confirm Password' ? 'Enter your $label' : 'Enter your Password',
              hintStyle: TextStyle(
                color: Colors.grey,
                fontSize: ResponsiveUtils.getSmallTextSize(context) + 1,
              ),
              prefixIcon: Icon(
                icon, 
                color: AppTheme.authAccentColor, 
                size: ResponsiveUtils.getIconSize(context) * 0.9,
              ),
              prefixIconConstraints: BoxConstraints(
                minWidth: ResponsiveUtils.isSmallPhone(context) ? 35 : 45
              ),
              suffixIcon: isPassword
                ? IconButton(
                    icon: Icon(
                      (isConfirmPassword ? _confirmPasswordVisible : _passwordVisible) 
                          ? Icons.visibility 
                          : Icons.visibility_off,
                      color: AppTheme.authAccentColor,
                      size: ResponsiveUtils.getIconSize(context) * 0.9,
                    ),
                    onPressed: () {
                      setState(() {
                        if (isConfirmPassword) {
                          _confirmPasswordVisible = !_confirmPasswordVisible;
                        } else {
                          _passwordVisible = !_passwordVisible;
                        }
                      });
                    },
                  )
                : null,
              // Border styling exactly matching signin page
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authDividerColor, width: 1),
              ),
              enabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authDividerColor, width: 1),
              ),
              focusedBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authAccentColor, width: 1.5),
              ),
              errorBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authErrorColor.withOpacity(0.5), width: 1),
              ),
              focusedErrorBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authErrorColor, width: 1.5),
              ),
              disabledBorder: OutlineInputBorder(
                borderRadius: BorderRadius.circular(30),
                borderSide: BorderSide(color: AppTheme.authDividerColor.withOpacity(0.5), width: 1),
              ),
            ),
            cursorColor: AppTheme.authAccentColor,
          ),
        ),
      ],
    );
  }
}

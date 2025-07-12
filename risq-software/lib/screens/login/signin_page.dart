import 'package:flutter/material.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/services/storage_service.dart';
import 'package:video_player/video_player.dart';
import 'package:font_awesome_flutter/font_awesome_flutter.dart';
import 'package:risq/theme/theme.dart';
import 'dart:math' as math;
import 'signup_page.dart';
import '../main_navigation.dart';

class SignInPage extends StatefulWidget {
  const SignInPage({super.key});

  @override
  _SignInPageState createState() => _SignInPageState();
}

class _SignInPageState extends State<SignInPage> with TickerProviderStateMixin {
  // Add controllers
  final _emailController = TextEditingController();
  final _passwordController = TextEditingController();
  
  // Add state variables
  bool _rememberMe = false;
  bool _isLoading = false;
  String _errorMessage = '';
  bool _passwordVisible = false;

  // Animation controllers for slide-up effect
  AnimationController? _animationController;
  Animation<double>? _containerAnimation;
  Animation<double>? _fadeAnimation;
  
  // Animation controller for gradient dissolution effect
  AnimationController? _gradientController;
  
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
    
    // Initialize gradient animation controller
    _gradientController = AnimationController(
      duration: const Duration(seconds: 4),
      vsync: this,
    )..repeat();
    
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
      }
    });
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
    _gradientController?.dispose();
    _videoController?.dispose();
    _emailController.dispose();
    _passwordController.dispose();
    super.dispose();
  }

  Future<void> _signInWithEmailAndPassword() async {
    // Clear previous error messages
    setState(() {
      _isLoading = true;
      _errorMessage = '';
    });
    
    // Validate input fields
    if (_emailController.text.trim().isEmpty) {
      setState(() {
        _errorMessage = 'Please enter your email';
        _isLoading = false;
      });
      return;
    }
    
    if (_passwordController.text.isEmpty) {
      setState(() {
        _errorMessage = 'Please enter your password';
        _isLoading = false;
      });
      return;
    }
    
    try {
      // Check if user has an account
      final hasAccount = await StorageService.hasAccount();
      if (!hasAccount) {
        setState(() {
          _errorMessage = 'No account found. Please sign up first.';
          _isLoading = false;
        });
        return;
      }

      // Verify login credentials using stored data
      final userData = await StorageService.verifyLogin(
        email: _emailController.text.trim(),
        password: _passwordController.text,
      );
      
      if (userData != null) {
        // Login successful - navigate to main navigation
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(
            builder: (context) => MainNavigation(
              userName: userData['name'] ?? '',
              userEmail: userData['email'] ?? '',
            ),
          ),
        );
      } else {
        setState(() {
          _errorMessage = 'Invalid email or password. Please try again.';
        });
      }
    } catch (e) {
      setState(() {
        _errorMessage = 'Unable to sign in. Please try again.';
      });
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
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
          
          // Main content
          SafeArea(
            child: LayoutBuilder(
              builder: (context, constraints) {
                // Reduce bottom container height further
                final bottomContainerHeight = constraints.maxHeight * 
                    (ResponsiveUtils.getBottomContainerRatio(context) - 0.05);
                
                return Stack(
                  children: [
                    // Bottom container
                    AnimatedBuilder(
                      animation: _animationController ?? const AlwaysStoppedAnimation(0),
                      builder: (context, child) {
                        return Positioned(
                          bottom: -bottomContainerHeight * (_containerAnimation?.value ?? 0.0),
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
                          child: AnimatedBuilder(
                            animation: _fadeAnimation ?? const AlwaysStoppedAnimation(0),
                            builder: (context, child) {
                              return Opacity(
                                opacity: _fadeAnimation?.value ?? 1.0,
                                child: child,
                              );
                            },
                            child: Column(
                              mainAxisSize: MainAxisSize.min,
                              children: [
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context)),
                                Text(
                                  "Sign In",
                                  style: TextStyle(
                                    fontSize: ResponsiveUtils.getHeadingSize(context),
                                    fontWeight: FontWeight.bold,
                                    color: AppTheme.authPrimaryColor,
                                  ),
                                ),
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                _buildTextField(
                                  label: 'Email',
                                  controller: _emailController,
                                  icon: Icons.email,
                                  keyboardType: TextInputType.emailAddress,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context)),
                                _buildTextField(
                                  label: 'Password',
                                  controller: _passwordController,
                                  icon: Icons.lock,
                                  isPassword: true,
                                ),
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context)),
                                
                                // Remember me & Forgot password
                                Row(
                                  mainAxisAlignment: MainAxisAlignment.spaceBetween,
                                  children: [
                                    Row(
                                      children: [
                                        Transform.scale(
                                          scale: ResponsiveUtils.isSmallPhone(context) ? 0.9 : 1.0,
                                          child: Checkbox(
                                            value: _rememberMe,
                                            onChanged: (value) {
                                              setState(() {
                                                _rememberMe = value ?? false;
                                              });
                                            },
                                            activeColor: AppTheme.authAccentColor,
                                            checkColor: Colors.white, // Explicit white check mark
                                            // Add these properties for the outline
                                            side: BorderSide(
                                              color: AppTheme.authDividerColor.withOpacity(0.8), 
                                              width: 1.5,
                                            ),
                                            shape: RoundedRectangleBorder(
                                              borderRadius: BorderRadius.circular(4),
                                            ),
                                          ),
                                        ),
                                        Text(
                                          'Remember me',
                                          style: TextStyle(
                                            fontSize: ResponsiveUtils.getSmallTextSize(context),
                                            color: AppTheme.authTextColor,
                                          ),
                                        ),
                                      ],
                                    ),
                                    TextButton(
                                      onPressed: () {},
                                      style: TextButton.styleFrom(
                                        padding: EdgeInsets.symmetric(
                                            horizontal: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                        minimumSize: Size(0, 0),
                                      ),
                                      child: Text(
                                        'Forgot Password?',
                                        style: TextStyle(
                                          color: AppTheme.authAccentColor,
                                          fontWeight: FontWeight.w500,
                                          fontSize: ResponsiveUtils.getSmallTextSize(context),
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
                                
                                // Sign In button with randomized static gradient
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
                                    onPressed: _isLoading ? null : _signInWithEmailAndPassword,
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
                                          'Sign In', 
                                          style: TextStyle(
                                            fontSize: ResponsiveUtils.getBodySize(context),
                                            fontWeight: FontWeight.bold,
                                            color: Colors.white,
                                          ),
                                        ),
                                  ),
                                ),
                                
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                
                                // Divider with "OR"
                                Row(
                                  children: [
                                    Expanded(
                                      child: Divider(
                                        color: AppTheme.authDividerColor,
                                        thickness: 1,
                                      ),
                                    ),
                                    Padding(
                                      padding: EdgeInsets.symmetric(horizontal: ResponsiveUtils.getSmallSpace(context)),
                                      child: Text(
                                        'OR',
                                        style: TextStyle(
                                          color: AppTheme.authTextColor.withOpacity(0.6),
                                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                                          fontWeight: FontWeight.w500,
                                        ),
                                      ),
                                    ),
                                    Expanded(
                                      child: Divider(
                                        color: AppTheme.authDividerColor,
                                        thickness: 1,
                                      ),
                                    ),
                                  ],
                                ),
                                
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                
                                // Social Sign In Buttons - Side by Side
                                Row(
                                  children: [
                                    // Google Sign In Button
                                    Expanded(
                                      child: Container(
                                        height: ResponsiveUtils.getButtonHeight(context),
                                        decoration: BoxDecoration(
                                          color: AppTheme.authAccentColor,
                                          borderRadius: BorderRadius.circular(30),
                                          boxShadow: [
                                            BoxShadow(
                                              color: AppTheme.authAccentColor.withOpacity(0.3),
                                              spreadRadius: 1,
                                              blurRadius: 8,
                                              offset: Offset(0, 2),
                                            ),
                                          ],
                                        ),
                                        child: ElevatedButton.icon(
                                          onPressed: () {
                                            // TODO: Implement Google Sign In
                                            print('Google Sign In pressed');
                                          },
                                          style: ElevatedButton.styleFrom(
                                            backgroundColor: AppTheme.authAccentColor,
                                            foregroundColor: Colors.white,
                                            shadowColor: Colors.transparent,
                                            elevation: 0,
                                            shape: RoundedRectangleBorder(
                                              borderRadius: BorderRadius.circular(30),
                                            ),
                                          ),
                                          icon: FaIcon(
                                            FontAwesomeIcons.google,
                                            size: ResponsiveUtils.getIconSize(context) * 0.8,
                                            color: Colors.white,
                                          ),
                                          label: Text(
                                            'Google',
                                            style: TextStyle(
                                              fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
                                              fontWeight: FontWeight.w600,
                                              color: Colors.white,
                                            ),
                                          ),
                                        ),
                                      ),
                                    ),
                                    
                                    SizedBox(width: ResponsiveUtils.getSmallSpace(context)),
                                    
                                    // Facebook Sign In Button
                                    Expanded(
                                      child: Container(
                                        height: ResponsiveUtils.getButtonHeight(context),
                                        decoration: BoxDecoration(
                                          color: AppTheme.authAccentColor,
                                          borderRadius: BorderRadius.circular(30),
                                          boxShadow: [
                                            BoxShadow(
                                              color: AppTheme.authAccentColor.withOpacity(0.3),
                                              spreadRadius: 1,
                                              blurRadius: 8,
                                              offset: Offset(0, 3),
                                            ),
                                          ],
                                        ),
                                        child: ElevatedButton.icon(
                                          onPressed: () {
                                            // TODO: Implement Facebook Sign In
                                            print('Facebook Sign In pressed');
                                          },
                                          style: ElevatedButton.styleFrom(
                                            backgroundColor: AppTheme.authAccentColor,
                                            foregroundColor: Colors.white,
                                            shadowColor: Colors.transparent,
                                            elevation: 0,
                                            shape: RoundedRectangleBorder(
                                              borderRadius: BorderRadius.circular(30),
                                            ),
                                          ),
                                          icon: FaIcon(
                                            FontAwesomeIcons.facebookF,
                                            size: ResponsiveUtils.getIconSize(context) * 0.8,
                                            color: Colors.white,
                                          ),
                                          label: Text(
                                            'Facebook',
                                            style: TextStyle(
                                              fontSize: ResponsiveUtils.getBodySize(context) * 0.9,
                                              fontWeight: FontWeight.w600,
                                              color: Colors.white,
                                            ),
                                          ),
                                        ),
                                      ),
                                    ),
                                  ],
                                ),
                                
                                SizedBox(height: ResponsiveUtils.getMediumSpace(context)),
                                
                                // Sign Up link
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
                                        "Don't have an account? ",
                                        style: TextStyle(
                                          color: AppTheme.authTextColor,
                                          fontSize: ResponsiveUtils.getSmallTextSize(context),
                                        ),
                                      ),
                                      TextButton(
                                        onPressed: () {
                                          Navigator.push(
                                            context,
                                            MaterialPageRoute(builder: (context) => SignupPage()),
                                          );
                                        },
                                        style: TextButton.styleFrom(
                                          padding: EdgeInsets.symmetric(
                                              horizontal: ResponsiveUtils.getSmallSpace(context) * 0.8),
                                          minimumSize: Size(0, 0),
                                        ),
                                        child: Text(
                                          "Sign Up",
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
                                SizedBox(height: ResponsiveUtils.getSmallSpace(context)),
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
    required TextEditingController controller,
    required IconData icon,
    bool isPassword = false,
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
          height: ResponsiveUtils.getInputHeight(context),
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
            obscureText: isPassword && !_passwordVisible,
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
              hintText: 'Enter your $label',
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
                      _passwordVisible ? Icons.visibility : Icons.visibility_off,
                      color: AppTheme.authAccentColor,
                      size: ResponsiveUtils.getIconSize(context) * 0.9,
                    ),
                    onPressed: () {
                      setState(() {
                        _passwordVisible = !_passwordVisible;
                      });
                    },
                  )
                : null,
              // Use fixed light theme styles
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
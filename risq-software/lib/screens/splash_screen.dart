import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/utils/responsive_utils.dart';
import 'package:risq/services/health_service.dart';
import 'package:risq/screens/login/signin_page.dart';

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> 
    with SingleTickerProviderStateMixin {
  
  late AnimationController _animationController;
  late Animation<double> _fadeAnimation;
  late Animation<double> _scaleAnimation;
  
  String _statusMessage = 'Initializing...';
  bool _isError = false;
  bool _showRetry = false;

  @override
  void initState() {
    super.initState();
    _initializeAnimations();
    _performHealthCheck();
  }

  void _initializeAnimations() {
    _animationController = AnimationController(
      duration: Duration(milliseconds: 1500),
      vsync: this,
    );

    _fadeAnimation = Tween<double>(
      begin: 0.0,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.easeInOut,
    ));

    _scaleAnimation = Tween<double>(
      begin: 0.8,
      end: 1.0,
    ).animate(CurvedAnimation(
      parent: _animationController,
      curve: Curves.elasticOut,
    ));

    _animationController.forward();
  }

  Future<void> _performHealthCheck() async {
    await Future.delayed(Duration(milliseconds: 800)); // Give animations time to start
    
    setState(() {
      _statusMessage = 'Connecting to server...';
    });

    try {
      final healthResult = await HealthService.checkHealth();
      
      if (healthResult['success'] == true) {
        setState(() {
          _statusMessage = 'Server connection established ✓';
          _isError = false;
        });
        
        // Wait a moment to show success message
        await Future.delayed(Duration(milliseconds: 1000));
        
        // Navigate to login page
        if (mounted) {
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(builder: (context) => SignInPage()),
          );
        }
      } else {
        _handleHealthCheckError(healthResult['message'] ?? 'Unknown error');
      }
    } catch (e) {
      _handleHealthCheckError(e.toString());
    }
  }

  void _handleHealthCheckError(String errorMessage) {
    setState(() {
      _statusMessage = 'Connection failed: $errorMessage';
      _isError = true;
      _showRetry = true;
    });
  }

  void _retryHealthCheck() {
    setState(() {
      _isError = false;
      _showRetry = false;
      _statusMessage = 'Retrying connection...';
    });
    _performHealthCheck();
  }

  @override
  void dispose() {
    _animationController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: Container(
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
        child: SafeArea(
          child: Center(
            child: AnimatedBuilder(
              animation: _animationController,
              builder: (context, child) {
                return FadeTransition(
                  opacity: _fadeAnimation,
                  child: ScaleTransition(
                    scale: _scaleAnimation,
                    child: Column(
                      mainAxisAlignment: MainAxisAlignment.center,
                      children: [
                        // App Logo/Title
                        Container(
                          padding: EdgeInsets.all(32),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.1),
                            borderRadius: BorderRadius.circular(20),
                            border: Border.all(
                              color: Colors.white.withOpacity(0.2),
                              width: 2,
                            ),
                          ),
                          child: Column(
                            children: [
                              Icon(
                                Icons.trending_up,
                                size: 80,
                                color: Colors.white,
                              ),
                              SizedBox(height: 16),
                              Text(
                                'RisQ',
                                style: AppTheme.headingTextStyle.copyWith(
                                  fontSize: ResponsiveUtils.getHeadingSize(context) * 1.5,
                                  color: Colors.white,
                                  fontWeight: FontWeight.bold,
                                ),
                              ),
                              SizedBox(height: 8),
                              Text(
                                'Risk Assessment Platform',
                                style: AppTheme.regularTextStyle.copyWith(
                                  fontSize: ResponsiveUtils.getBodySize(context),
                                  color: Colors.white70,
                                ),
                              ),
                            ],
                          ),
                        ),
                        
                        SizedBox(height: 60),
                        
                        // Status Section
                        Container(
                          width: double.infinity,
                          margin: EdgeInsets.symmetric(horizontal: 32),
                          padding: EdgeInsets.all(24),
                          decoration: BoxDecoration(
                            color: Colors.white.withOpacity(0.9),
                            borderRadius: BorderRadius.circular(16),
                            boxShadow: [
                              BoxShadow(
                                color: Colors.black.withOpacity(0.1),
                                blurRadius: 10,
                                offset: Offset(0, 5),
                              ),
                            ],
                          ),
                          child: Column(
                            children: [
                              // Loading indicator or error icon
                              if (!_isError && !_showRetry)
                                SizedBox(
                                  width: 24,
                                  height: 24,
                                  child: CircularProgressIndicator(
                                    strokeWidth: 3,
                                    valueColor: AlwaysStoppedAnimation<Color>(
                                      AppTheme.authAccentColor,
                                    ),
                                  ),
                                )
                              else if (_isError)
                                Icon(
                                  Icons.error_outline,
                                  size: 32,
                                  color: Colors.red[600],
                                ),
                              
                              SizedBox(height: 16),
                              
                              // Status message
                              Text(
                                _statusMessage,
                                textAlign: TextAlign.center,
                                style: AppTheme.regularTextStyle.copyWith(
                                  fontSize: ResponsiveUtils.getBodySize(context),
                                  color: _isError ? Colors.red[700] : Colors.black87,
                                  fontWeight: FontWeight.w500,
                                ),
                              ),
                              
                              // Retry button
                              if (_showRetry) ...[
                                SizedBox(height: 20),
                                ElevatedButton.icon(
                                  onPressed: _retryHealthCheck,
                                  icon: Icon(Icons.refresh, color: Colors.white),
                                  label: Text(
                                    'Retry Connection',
                                    style: TextStyle(
                                      color: Colors.white,
                                      fontWeight: FontWeight.w600,
                                    ),
                                  ),
                                  style: ElevatedButton.styleFrom(
                                    backgroundColor: AppTheme.authAccentColor,
                                    padding: EdgeInsets.symmetric(
                                      horizontal: 24,
                                      vertical: 12,
                                    ),
                                    shape: RoundedRectangleBorder(
                                      borderRadius: BorderRadius.circular(30),
                                    ),
                                  ),
                                ),
                                
                                SizedBox(height: 12),
                                
                                TextButton(
                                  onPressed: () {
                                    // Skip health check and go to login
                                    Navigator.pushReplacement(
                                      context,
                                      MaterialPageRoute(builder: (context) => SignInPage()),
                                    );
                                  },
                                  child: Text(
                                    'Continue Anyway',
                                    style: TextStyle(
                                      color: Colors.grey[600],
                                      fontSize: ResponsiveUtils.getSmallTextSize(context),
                                    ),
                                  ),
                                ),
                              ],
                            ],
                          ),
                        ),
                        
                        SizedBox(height: 40),
                        
                        // App version or additional info
                        Text(
                          'Powered by AI Risk Assessment',
                          style: TextStyle(
                            color: Colors.white60,
                            fontSize: ResponsiveUtils.getSmallTextSize(context),
                          ),
                        ),
                      ],
                    ),
                  ),
                );
              },
            ),
          ),
        ),
      ),
    );
  }
}

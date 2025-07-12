import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:risq/screens/login/signin_page.dart';
import 'package:risq/theme/theme.dart';

void main() async {
  WidgetsFlutterBinding.ensureInitialized();
  
  // Set preferred orientations to portrait only
  await SystemChrome.setPreferredOrientations([
    DeviceOrientation.portraitUp,
    DeviceOrientation.portraitDown,
  ]);
  
  // Set system UI overlay style for light theme
  SystemChrome.setSystemUIOverlayStyle(
    const SystemUiOverlayStyle(
      statusBarColor: Colors.transparent,
      statusBarIconBrightness: Brightness.dark, // For Android (dark icons)
      statusBarBrightness: Brightness.light, // For iOS (dark icons)
    ),
  );
  
  runApp(const RisQ());
}

class RisQ extends StatelessWidget {
  const RisQ({super.key});

  @override
  Widget build(BuildContext context) {
    return MaterialApp(
      title: 'RisQ',
      debugShowCheckedModeBanner: false,
      theme: AppTheme.lightTheme,
      themeMode: ThemeMode.system,
      builder: (context, child) {
        // Update status bar based on current theme brightness
        final brightness = MediaQuery.platformBrightnessOf(context);
        SystemChrome.setSystemUIOverlayStyle(
          SystemUiOverlayStyle(
            statusBarColor: Colors.transparent,
            statusBarIconBrightness: brightness == Brightness.dark 
                ? Brightness.light 
                : Brightness.dark,
          )
        );
        
        // Apply a maximum width constraint to the entire app
        return MediaQuery(
          // Preserve original media query data but apply our constraints
          data: MediaQuery.of(context).copyWith(
            // Optional: adjust text scaling to ensure consistent text size
            textScaleFactor: 1.0,
          ),
          child: child!,
        );
      },
      home: SplashScreen(),
    );
  }
}

class SplashScreen extends StatefulWidget {
  const SplashScreen({super.key});

  @override
  State<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends State<SplashScreen> {
  @override
  void initState() {
    super.initState();
    _checkLoginStatus();
  }

  Future<void> _checkLoginStatus() async {
    // Add a small delay for better UX
    await Future.delayed(Duration(milliseconds: 1500));
    
    // Always navigate to sign-in page (no auto-login)
    if (mounted) {
      Navigator.pushReplacement(
        context,
        MaterialPageRoute(builder: (context) => SignInPage()),
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Colors.white,
      body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            // App Logo or Name
            Text(
              'RisQ',
              style: TextStyle(
                fontSize: 48,
                fontWeight: FontWeight.bold,
                color: AppTheme.authAccentColor,
              ),
            ),
            SizedBox(height: 20),
            // Loading indicator
            CircularProgressIndicator(
              color: AppTheme.authAccentColor,
            ),
          ],
        ),
      ),
    );
  }
}
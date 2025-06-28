import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:risq/config/app_config.dart';

class AuthService {
  // Get base URL from configuration
  static String get baseUrl => AppConfig.backendUrl;
  
  // Login method
  static Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('${baseUrl}/api/v1/auth/login'),
        headers: {
          'Content-Type': 'application/json',
        },
        body: json.encode({
          'email': email,
          'password': password,
        }),
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        // Store token and user data for future requests
        if (responseData['token'] != null) {
          await storeAuthToken(responseData['token']);
        }
        if (responseData['user'] != null) {
          await storeUserData(responseData['user']);
        }
        
        return {
          'success': true,
          'data': responseData,
          'token': responseData['token'],
          'user': responseData['user'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Login failed',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Auth Service Error: $e');
      }
      
      if (e.toString().contains('SocketException') || 
          e.toString().contains('TimeoutException')) {
        return {
          'success': false,
          'message': 'Network error. Please check your connection.',
        };
      }
      
      return {
        'success': false,
        'message': 'An error occurred. Please try again.',
      };
    }
  }

  // Register method
  static Future<Map<String, dynamic>> register({
    required String name,
    required String email,
    required String password,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('${baseUrl}/api/v1/auth/signup'),
        headers: {
          'Content-Type': 'application/json',
        },
        body: json.encode({
          'name': name,
          'email': email,
          'password': password,
        }),
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 201 || response.statusCode == 200) {
        // Store token and user data for future requests
        if (responseData['token'] != null) {
          await storeAuthToken(responseData['token']);
        }
        if (responseData['user'] != null) {
          await storeUserData(responseData['user']);
        }
        
        return {
          'success': true,
          'data': responseData,
          'token': responseData['token'],
          'user': responseData['user'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Registration failed',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Auth Service Error: $e');
      }
      
      if (e.toString().contains('SocketException') || 
          e.toString().contains('TimeoutException')) {
        return {
          'success': false,
          'message': 'Network error. Please check your connection.',
        };
      }
      
      return {
        'success': false,
        'message': 'An error occurred. Please try again.',
      };
    }
  }

  // Logout method
  static Future<void> logout() async {
    // Clear stored tokens/user data
    await clearAuthToken();
    await clearUserData();
  }

  // Store authentication token
  static Future<void> storeAuthToken(String token) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('auth_token', token);
  }

  // Get stored authentication token
  static Future<String?> getAuthToken() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getString('auth_token');
  }

  // Clear authentication token
  static Future<void> clearAuthToken() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('auth_token');
  }

  // Store user data
  static Future<void> storeUserData(Map<String, dynamic> userData) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString('user_data', json.encode(userData));
  }

  // Get stored user data
  static Future<Map<String, dynamic>?> getUserData() async {
    final prefs = await SharedPreferences.getInstance();
    final userDataString = prefs.getString('user_data');
    if (userDataString != null) {
      return json.decode(userDataString);
    }
    return null;
  }

  // Clear user data
  static Future<void> clearUserData() async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.remove('user_data');
  }
}

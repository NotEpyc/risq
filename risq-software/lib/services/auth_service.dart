import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';

class AuthService {
  // Replace this with your actual Railways app URL
  static const String _baseUrl = 'http://resqbackend-production.up.railway.app';
  
  // Login method
  static Future<Map<String, dynamic>> login({
    required String email,
    required String password,
  }) async {
    try {
      final response = await http.post(
        Uri.parse('$_baseUrl/auth/login'),
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
        Uri.parse('$_baseUrl/auth/register'),
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
    // You might want to implement token clearing logic here
    // For example, using SharedPreferences or secure storage
  }

  // Optional: Store authentication token
  // static Future<void> storeAuthToken(String token) async {
  //   final prefs = await SharedPreferences.getInstance();
  //   await prefs.setString('auth_token', token);
  // }

  // Optional: Get stored authentication token
  // static Future<String?> getAuthToken() async {
  //   final prefs = await SharedPreferences.getInstance();
  //   return prefs.getString('auth_token');
  // }
}

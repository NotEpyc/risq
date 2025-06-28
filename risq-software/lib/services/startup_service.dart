import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:risq/services/auth_service.dart';

class StartupService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  
  // Submit startup onboarding data
  static Future<Map<String, dynamic>> onboardStartup({
    required Map<String, dynamic> startupData,
    String? authToken,
  }) async {
    try {
      // Get auth token if not provided
      final token = authToken ?? await AuthService.getAuthToken();
      
      if (token == null) {
        return {
          'success': false,
          'message': 'Authentication error. Please login again.',
        };
      }

      final response = await http.post(
        Uri.parse('$_baseUrl/api/v1/startup/onboard'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
        body: json.encode(startupData),
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return {
          'success': true,
          'data': responseData,
          'startup': responseData['startup'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to submit startup data',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Startup Service Error: $e');
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
        'message': 'Unable to submit startup information. Please try again.',
      };
    }
  }

  // Get startup profile
  static Future<Map<String, dynamic>> getStartupProfile({
    String? authToken,
  }) async {
    try {
      // Get auth token if not provided
      final token = authToken ?? await AuthService.getAuthToken();
      
      if (token == null) {
        return {
          'success': false,
          'message': 'Authentication error. Please login again.',
        };
      }

      final response = await http.get(
        Uri.parse('$_baseUrl/api/v1/startup/profile'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': responseData,
          'startup': responseData['startup'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to fetch startup profile',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Startup Service Error: $e');
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
        'message': 'Unable to fetch startup profile. Please try again.',
      };
    }
  }
}

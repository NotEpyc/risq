import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:risq/services/auth_service.dart';

class SpeculationService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  
  // Submit speculation for AI analysis using the correct decision route
  static Future<Map<String, dynamic>> submitSpeculation({
    required String description,
    required String category,
    required String context,
    required String timeline,
    required String budget,
    String? authToken,
  }) async {
    try {
      // Get auth token if not provided
      final token = authToken ?? await AuthService.getAuthToken();
      
      final headers = {
        'Content-Type': 'application/json',
      };
      
      // Add authorization header if token is available
      if (token != null) {
        headers['Authorization'] = 'Bearer $token';
      }

      final response = await http.post(
        Uri.parse('$_baseUrl/api/v1/decisions/speculate'),
        headers: headers,
        body: json.encode({
          'description': description,
          'category': category,
          'context': context,
          'timeline': timeline,
          'budget': budget,
        }),
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return {
          'success': true,
          'data': responseData,
          'analysis': responseData['analysis'],
          'riskScore': responseData['riskScore'],
          'confidence': responseData['confidence'],
          'recommendations': responseData['recommendations'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Analysis failed',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Speculation Service Error: $e');
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

  // Get speculation history from decisions endpoint
  static Future<Map<String, dynamic>> getSpeculationHistory({
    String? authToken,
  }) async {
    try {
      // Get auth token if not provided
      final token = authToken ?? await AuthService.getAuthToken();
      
      final headers = <String, String>{
        'Content-Type': 'application/json',
      };
      
      if (token != null) {
        headers['Authorization'] = 'Bearer $token';
      }

      final response = await http.get(
        Uri.parse('$_baseUrl/api/v1/decisions/'),
        headers: headers,
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': responseData,
          'decisions': responseData['decisions'] ?? [],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to fetch history',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Speculation Service Error: $e');
      }
      
      return {
        'success': false,
        'message': 'Failed to fetch speculation history.',
      };
    }
  }
}

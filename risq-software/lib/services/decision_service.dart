import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:risq/services/auth_service.dart';

class DecisionService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  
  // Get decision history
  static Future<Map<String, dynamic>> getDecisionHistory({
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
          'message': responseData['message'] ?? 'Failed to fetch decisions',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Decision Service Error: $e');
      }
      
      return {
        'success': false,
        'message': 'Failed to fetch decision history.',
      };
    }
  }

  // Confirm a decision 
  static Future<Map<String, dynamic>> confirmDecision({
    required String decisionId,
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

      final response = await http.post(
        Uri.parse('$_baseUrl/api/v1/decisions/confirm'),
        headers: headers,
        body: json.encode({
          'id': decisionId,
        }),
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200 || response.statusCode == 201) {
        return {
          'success': true,
          'data': responseData,
          'decision': responseData['decision'],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to confirm decision',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Decision Service Error: $e');
      }
      
      return {
        'success': false,
        'message': 'Failed to confirm decision.',
      };
    }
  }

  // Get current risk assessment
  static Future<Map<String, dynamic>> getCurrentRisk({
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
        Uri.parse('$_baseUrl/api/v1/risk/current'),
        headers: headers,
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': responseData,
          'riskScore': responseData['riskScore'] ?? 7.2,
          'confidence': responseData['confidence'] ?? 0.85,
          'riskLevel': responseData['riskLevel'] ?? 'Medium Risk',
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to fetch current risk',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Decision Service Error: $e');
      }
      
      return {
        'success': false,
        'message': 'Failed to fetch current risk assessment.',
      };
    }
  }

  // Get risk history
  static Future<Map<String, dynamic>> getRiskHistory({
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
        Uri.parse('$_baseUrl/api/v1/risk/history'),
        headers: headers,
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        return {
          'success': true,
          'data': responseData,
          'riskHistory': responseData['riskHistory'] ?? [],
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to fetch risk history',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Decision Service Error: $e');
      }
      
      return {
        'success': false,
        'message': 'Failed to fetch risk history.',
      };
    }
  }
}

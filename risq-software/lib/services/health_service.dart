import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:risq/services/auth_service.dart';

class HealthService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  
  // Health check endpoint
  static Future<Map<String, dynamic>> checkHealth() async {
    try {
      final response = await http.get(
        Uri.parse('$_baseUrl/health'),
        headers: {
          'Content-Type': 'application/json',
        },
      ).timeout(
        Duration(seconds: 10), // 10 second timeout
        onTimeout: () {
          throw Exception('Health check timeout - server may be down');
        },
      );

      if (response.statusCode == 200) {
        final responseData = json.decode(response.body);
        return {
          'success': true,
          'message': 'Server is healthy',
          'data': responseData,
          'statusCode': response.statusCode,
        };
      } else {
        return {
          'success': false,
          'message': 'Server returned status ${response.statusCode}',
          'statusCode': response.statusCode,
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Health Service Error: $e');
      }
      
      String errorMessage;
      if (e.toString().contains('SocketException')) {
        errorMessage = 'No internet connection or server is unreachable';
      } else if (e.toString().contains('TimeoutException') || e.toString().contains('timeout')) {
        errorMessage = 'Server is taking too long to respond';
      } else if (e.toString().contains('FormatException')) {
        errorMessage = 'Server returned invalid response format';
      } else {
        errorMessage = 'Unable to connect to server: ${e.toString()}';
      }
      
      return {
        'success': false,
        'message': errorMessage,
        'error': e.toString(),
      };
    }
  }

  // Quick ping to check if server is reachable
  static Future<bool> isServerReachable() async {
    try {
      final result = await checkHealth();
      return result['success'] == true;
    } catch (e) {
      return false;
    }
  }
}

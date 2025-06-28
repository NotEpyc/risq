import 'dart:convert';
import 'package:http/http.dart' as http;
import 'package:flutter/foundation.dart';
import 'package:risq/services/auth_service.dart';

class NotificationService {
  // Get base URL from AuthService
  static String get _baseUrl => AuthService.baseUrl;
  
  // Get risk-based notifications
  static Future<Map<String, dynamic>> getRiskNotifications({
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

      // Fetch current risk data to generate notifications
      final response = await http.get(
        Uri.parse('$_baseUrl/api/v1/risk/current'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      final responseData = json.decode(response.body);

      if (response.statusCode == 200) {
        final riskData = responseData['data'];
        final notifications = _generateRiskNotifications(riskData);
        
        return {
          'success': true,
          'data': notifications,
          'totalCount': notifications.length,
        };
      } else {
        return {
          'success': false,
          'message': responseData['message'] ?? 'Failed to fetch risk notifications',
        };
      }
    } catch (e) {
      if (kDebugMode) {
        print('Notification Service Error: $e');
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
        'message': 'Unable to fetch notifications. Please try again.',
      };
    }
  }

  // Generate notifications based on risk data
  static List<Map<String, dynamic>> _generateRiskNotifications(Map<String, dynamic> riskData) {
    final List<Map<String, dynamic>> notifications = [];
    final now = DateTime.now();
    
    // Extract risk information
    final double score = (riskData['score'] ?? 7.2).toDouble();
    final String level = riskData['level'] ?? 'medium';
    final double confidence = (riskData['confidence'] ?? 0.85).toDouble();
    final List<String> suggestions = List<String>.from(riskData['suggestions'] ?? []);
    final List<String> factors = List<String>.from(riskData['factors'] ?? []);
    
    // Risk Score Change Notification
    notifications.add({
      'id': 'risk_score_${now.millisecondsSinceEpoch}',
      'type': 'risk_update',
      'title': 'Risk Score Updated',
      'message': 'Your current risk score is $score/10 (${_formatRiskLevel(level)})',
      'priority': _getRiskPriority(level),
      'timestamp': now.toIso8601String(),
      'isRead': false,
      'data': {
        'score': score,
        'level': level,
        'confidence': confidence,
      },
    });

    // High Priority Risk Notification
    if (level == 'high' || level == 'critical') {
      notifications.add({
        'id': 'high_risk_alert_${now.millisecondsSinceEpoch}',
        'type': 'alert',
        'title': level == 'critical' ? '🚨 Critical Risk Alert' : '⚠️ High Risk Warning',
        'message': level == 'critical' 
            ? 'Your startup is in critical risk territory. Immediate action required!'
            : 'Your risk level is high. Review AI suggestions and take action.',
        'priority': 'high',
        'timestamp': now.toIso8601String(),
        'isRead': false,
        'data': {
          'score': score,
          'level': level,
        },
      });
    }

    // AI Suggestions Notifications
    if (suggestions.isNotEmpty) {
      notifications.add({
        'id': 'ai_suggestions_${now.millisecondsSinceEpoch}',
        'type': 'suggestion',
        'title': '💡 New AI Suggestions Available',
        'message': 'AI has generated ${suggestions.length} new recommendations to improve your startup.',
        'priority': 'medium',
        'timestamp': now.toIso8601String(),
        'isRead': false,
        'data': {
          'suggestions': suggestions,
          'suggestionsCount': suggestions.length,
        },
      });
    }

    // Confidence Level Notification
    if (confidence < 0.7) {
      notifications.add({
        'id': 'low_confidence_${now.millisecondsSinceEpoch}',
        'type': 'info',
        'title': 'ℹ️ Analysis Confidence',
        'message': 'AI confidence is ${(confidence * 100).round()}%. Consider providing more data for better analysis.',
        'priority': 'low',
        'timestamp': now.toIso8601String(),
        'isRead': false,
        'data': {
          'confidence': confidence,
        },
      });
    }

    // Positive Factors Notification
    final positiveFactors = factors.where((factor) => 
        factor.toLowerCase().contains('strong') || 
        factor.toLowerCase().contains('good') ||
        factor.toLowerCase().contains('experienced') ||
        factor.toLowerCase().contains('differentiation')).toList();
    
    if (positiveFactors.isNotEmpty) {
      notifications.add({
        'id': 'positive_factors_${now.millisecondsSinceEpoch}',
        'type': 'success',
        'title': '✅ Positive Factors Identified',
        'message': 'AI found ${positiveFactors.length} positive factors in your startup analysis.',
        'priority': 'low',
        'timestamp': now.toIso8601String(),
        'isRead': false,
        'data': {
          'positiveFactors': positiveFactors,
        },
      });
    }

    // Sort by priority and timestamp
    notifications.sort((a, b) {
      final priorityOrder = {'high': 3, 'medium': 2, 'low': 1};
      final aPriority = priorityOrder[a['priority']] ?? 1;
      final bPriority = priorityOrder[b['priority']] ?? 1;
      
      if (aPriority != bPriority) {
        return bPriority.compareTo(aPriority); // High priority first
      }
      
      // If same priority, sort by timestamp (newest first)
      return DateTime.parse(b['timestamp']).compareTo(DateTime.parse(a['timestamp']));
    });

    return notifications;
  }

  // Helper methods
  static String _formatRiskLevel(String level) {
    switch (level.toLowerCase()) {
      case 'low': return 'Low Risk';
      case 'medium': return 'Medium Risk';
      case 'high': return 'High Risk';
      case 'critical': return 'Critical Risk';
      default: return level;
    }
  }

  static String _getRiskPriority(String level) {
    switch (level.toLowerCase()) {
      case 'critical': return 'high';
      case 'high': return 'high';
      case 'medium': return 'medium';
      case 'low': return 'low';
      default: return 'medium';
    }
  }

  // Mark notification as read
  static Future<bool> markAsRead(String notificationId) async {
    // This would typically update the backend, but for now we'll just return success
    // In a real implementation, you'd make an API call to mark the notification as read
    try {
      await Future.delayed(Duration(milliseconds: 100)); // Simulate API call
      return true;
    } catch (e) {
      return false;
    }
  }

  // Mark all notifications as read
  static Future<bool> markAllAsRead() async {
    try {
      await Future.delayed(Duration(milliseconds: 200)); // Simulate API call
      return true;
    } catch (e) {
      return false;
    }
  }
}

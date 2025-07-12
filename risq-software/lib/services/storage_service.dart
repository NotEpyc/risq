import 'package:shared_preferences/shared_preferences.dart';

class StorageService {
  static const String _keyIsLoggedIn = 'is_logged_in';
  static const String _keyUserName = 'user_name';
  static const String _keyUserEmail = 'user_email';
  static const String _keyUserPassword = 'user_password';
  static const String _keyUserIndustry = 'user_industry';
  static const String _keyUserCompanySize = 'user_company_size';
  static const String _keyUserExperience = 'user_experience';
  static const String _keyUserEducation = 'user_education';

  // Save user data during signup
  static Future<bool> saveUserData({
    required String name,
    required String email,
    required String password,
    String? industry,
    String? companySize,
    List<String>? experience,
    List<String>? education,
  }) async {
    try {
      print('StorageService: Attempting to save user data...');
      print('Name: $name, Email: $email');
      
      if (name.trim().isEmpty || email.trim().isEmpty || password.isEmpty) {
        print('StorageService: Error - required fields are empty');
        return false;
      }
      
      final prefs = await SharedPreferences.getInstance();
      print('StorageService: SharedPreferences instance obtained');
      
      // Save required data and check each operation
      bool allSuccessful = true;
      
      final nameResult = await prefs.setString(_keyUserName, name.trim());
      print('StorageService: Name saved: $nameResult');
      if (!nameResult) allSuccessful = false;
      
      final emailResult = await prefs.setString(_keyUserEmail, email.trim());
      print('StorageService: Email saved: $emailResult');
      if (!emailResult) allSuccessful = false;
      
      final passwordResult = await prefs.setString(_keyUserPassword, password);
      print('StorageService: Password saved: $passwordResult');
      if (!passwordResult) allSuccessful = false;
      
      if (!allSuccessful) {
        print('StorageService: Error - one or more save operations failed');
        return false;
      }
      
      // Verify the data was saved by reading it back
      await Future.delayed(Duration(milliseconds: 100)); // Small delay to ensure write completion
      
      final savedName = prefs.getString(_keyUserName);
      final savedEmail = prefs.getString(_keyUserEmail);
      final savedPassword = prefs.getString(_keyUserPassword);
      
      print('StorageService: Verification - Name: $savedName, Email: $savedEmail, Password exists: ${savedPassword != null}');
      
      if (savedName == null || savedEmail == null || savedPassword == null) {
        print('StorageService: Error - data verification failed - some data is null');
        return false;
      }
      
      if (savedName != name.trim() || savedEmail != email.trim() || savedPassword != password) {
        print('StorageService: Error - data verification failed - data mismatch');
        print('Expected: name=$name, email=$email');
        print('Actual: name=$savedName, email=$savedEmail');
        return false;
      }
      
      // Save optional data
      if (industry != null && industry.isNotEmpty) {
        await prefs.setString(_keyUserIndustry, industry);
        print('StorageService: Industry saved');
      }
      if (companySize != null && companySize.isNotEmpty) {
        await prefs.setString(_keyUserCompanySize, companySize);
        print('StorageService: Company size saved');
      }
      if (experience != null && experience.isNotEmpty) {
        await prefs.setStringList(_keyUserExperience, experience);
        print('StorageService: Experience saved');
      }
      if (education != null && education.isNotEmpty) {
        await prefs.setStringList(_keyUserEducation, education);
        print('StorageService: Education saved');
      }
      
      print('StorageService: All data saved and verified successfully');
      return true;
    } catch (e) {
      print('StorageService: Exception saving user data: $e');
      print('StorageService: Error details: ${e.toString()}');
      if (e is Error) {
        print('StorageService: Stack trace: ${e.stackTrace}');
      }
      return false;
    }
  }

  // Verify login credentials
  static Future<Map<String, dynamic>?> verifyLogin({
    required String email,
    required String password,
  }) async {
    try {
      final prefs = await SharedPreferences.getInstance();
      
      final storedEmail = prefs.getString(_keyUserEmail);
      final storedPassword = prefs.getString(_keyUserPassword);
      
      if (storedEmail == email && storedPassword == password) {
        // Login successful, mark as logged in
        await prefs.setBool(_keyIsLoggedIn, true);
        
        return {
          'name': prefs.getString(_keyUserName) ?? '',
          'email': storedEmail,
          'industry': prefs.getString(_keyUserIndustry),
          'companySize': prefs.getString(_keyUserCompanySize),
          'experience': prefs.getStringList(_keyUserExperience),
          'education': prefs.getStringList(_keyUserEducation),
        };
      }
      
      return null; // Login failed
    } catch (e) {
      print('Error verifying login: $e');
      return null;
    }
  }

  // Check if user is logged in
  static Future<bool> isLoggedIn() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      return prefs.getBool(_keyIsLoggedIn) ?? false;
    } catch (e) {
      print('Error checking login status: $e');
      return false;
    }
  }

  // Get stored user data
  static Future<Map<String, dynamic>?> getUserData() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      
      if (prefs.getBool(_keyIsLoggedIn) ?? false) {
        return {
          'name': prefs.getString(_keyUserName) ?? '',
          'email': prefs.getString(_keyUserEmail) ?? '',
          'industry': prefs.getString(_keyUserIndustry),
          'companySize': prefs.getString(_keyUserCompanySize),
          'experience': prefs.getStringList(_keyUserExperience),
          'education': prefs.getStringList(_keyUserEducation),
        };
      }
      
      return null;
    } catch (e) {
      print('Error getting user data: $e');
      return null;
    }
  }

  // Check if user has an account
  static Future<bool> hasAccount() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      final email = prefs.getString(_keyUserEmail);
      return email != null && email.isNotEmpty;
    } catch (e) {
      print('Error checking account: $e');
      return false;
    }
  }

  // Logout user
  static Future<bool> logout() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.setBool(_keyIsLoggedIn, false);
      return true;
    } catch (e) {
      print('Error logging out: $e');
      return false;
    }
  }

  // Clear all user data (for testing or account deletion)
  static Future<bool> clearAllData() async {
    try {
      final prefs = await SharedPreferences.getInstance();
      await prefs.clear();
      return true;
    } catch (e) {
      print('Error clearing data: $e');
      return false;
    }
  }

  // Update additional user data (for onboarding)
  static Future<bool> updateUserData({
    String? industry,
    String? companySize,
    List<String>? experience,
    List<String>? education,
  }) async {
    try {
      print('StorageService: Updating additional user data...');
      
      final prefs = await SharedPreferences.getInstance();
      print('StorageService: SharedPreferences instance obtained for update');
      
      if (industry != null) {
        await prefs.setString(_keyUserIndustry, industry);
        print('StorageService: Industry updated');
      }
      if (companySize != null) {
        await prefs.setString(_keyUserCompanySize, companySize);
        print('StorageService: Company size updated');
      }
      if (experience != null) {
        await prefs.setStringList(_keyUserExperience, experience);
        print('StorageService: Experience updated');
      }
      if (education != null) {
        await prefs.setStringList(_keyUserEducation, education);
        print('StorageService: Education updated');
      }
      
      print('StorageService: Additional data updated successfully');
      return true;
    } catch (e) {
      print('Error updating user data: $e');
      print('Error details: ${e.toString()}');
      return false;
    }
  }
}
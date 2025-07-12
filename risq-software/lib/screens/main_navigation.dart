import 'package:flutter/material.dart';
import 'package:risq/theme/theme.dart';
import 'package:risq/screens/pages/home_page.dart';
import 'package:risq/screens/pages/speculation_page.dart';
import 'package:risq/screens/pages/decision_page.dart';

class MainNavigation extends StatefulWidget {
  final String userName;
  final String userEmail;

  const MainNavigation({
    super.key,
    required this.userName,
    required this.userEmail,
  });

  @override
  State<MainNavigation> createState() => _MainNavigationState();
}

class _MainNavigationState extends State<MainNavigation> {
  int _currentIndex = 0;

  late final List<Widget> _pages;

  @override
  void initState() {
    super.initState();
    _pages = [
      HomePage(
        userName: widget.userName, 
        userEmail: widget.userEmail,
        onNavigateToTab: _onTabTapped,
      ),
      SpeculationPage(userName: widget.userName, userEmail: widget.userEmail),
      DecisionPage(userName: widget.userName, userEmail: widget.userEmail),
    ];
  }

  void _onTabTapped(int index) {
    setState(() {
      _currentIndex = index;
    });
  }
  Widget _buildNavItem(int index, IconData icon, String label) {
    final isSelected = _currentIndex == index;
    
    return Expanded(
      child: GestureDetector(
        onTap: () => _onTabTapped(index),
        child: AnimatedContainer(
          duration: Duration(milliseconds: 300),
          padding: EdgeInsets.symmetric(horizontal: 8, vertical: 12),
          decoration: BoxDecoration(
            border: isSelected ? Border.all(color: AppTheme.authAccentColor, width: 2) : null,
            borderRadius: BorderRadius.circular(30),
          ),
          child: Row(
            mainAxisSize: MainAxisSize.min,
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Icon(
                icon,
                size: 24,
                color: isSelected ? AppTheme.authAccentColor : Colors.grey[600],
              ),
              if (isSelected) ...[
                SizedBox(width: 6),
                Flexible(
                  child: Text(
                    label,
                    style: TextStyle(
                      fontSize: 12,
                      fontWeight: FontWeight.w600,
                      color: AppTheme.authAccentColor,
                    ),
                    overflow: TextOverflow.ellipsis,
                    maxLines: 1,
                  ),
                ),
              ],
            ],
          ),
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: IndexedStack(
        index: _currentIndex,
        children: _pages,
      ),
      bottomNavigationBar: Container(
        decoration: BoxDecoration(
          color: Colors.white,
          boxShadow: [
            BoxShadow(
              color: Colors.black.withOpacity(0.1),
              blurRadius: 10,
              offset: Offset(0, -2),
            ),
          ],
        ),
        child: SafeArea(
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 15.0, vertical: 8),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceAround,
              children: [
                _buildNavItem(0, Icons.home_outlined, 'Home'),
                _buildNavItem(1, Icons.psychology_outlined, 'Speculate'),
                _buildNavItem(2, Icons.fact_check_outlined, 'Decide'),
              ],
            ),
          ),
        ),
      ),
    );
  }
}

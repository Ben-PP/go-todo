import 'package:flutter/material.dart';

final themeData = ThemeData(
  colorScheme: ColorScheme(
    brightness: Brightness.dark,
    primary: Colors.cyan.shade900,
    onPrimary: Colors.white,
    secondary: Colors.teal,
    onSecondary: Colors.black,
    error: Colors.red,
    onError: Colors.white,
    surface: Colors.grey[900]!,
    onSurface: Colors.white70,
    primaryContainer: Colors.grey[800],
    onPrimaryContainer: Colors.white,
  ),
  useMaterial3: true,
  appBarTheme: AppBarTheme(
    backgroundColor: Colors.cyan.shade900,
    titleTextStyle: TextStyle(
      fontSize: 24,
      fontWeight: FontWeight.bold,
      color: Colors.white.withAlpha(220),
    ),
  ),
  elevatedButtonTheme: ElevatedButtonThemeData(
    style: ElevatedButton.styleFrom(
      backgroundColor: Colors.cyan.shade900,
      foregroundColor: Colors.white.withAlpha(220),
      textStyle: const TextStyle(fontSize: 24),
      padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
      minimumSize: const Size(100, 60),
    ),
  ),
  textTheme: const TextTheme(
    headlineMedium: TextStyle(fontSize: 36, fontWeight: FontWeight.bold),
    headlineSmall: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
    bodyMedium: TextStyle(fontSize: 18),
    labelLarge: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
    labelMedium: TextStyle(fontSize: 16),
  ),
);

import 'package:flutter/material.dart';
import 'package:go_todo/globals.dart';

SnackBar getSnackBar(
    {required BuildContext context,
    required Widget content,
    bool isError = false}) {
  const radius = BorderRadius.all(Radius.circular(5));
  const borderWidth = 3.0;
  final screenWidth = MediaQuery.sizeOf(context).width;
  final isMobile = screenWidth <= ScreenSize.small.value;

  RoundedRectangleBorder getShape(Color color) {
    return RoundedRectangleBorder(
      borderRadius: radius,
      side: BorderSide(
        color: color,
        width: borderWidth,
      ),
    );
  }

  final color = isError ? Colors.red : Colors.teal;

  return SnackBar(
    content: content,
    duration: const Duration(seconds: 5),
    shape: getShape(color),
    margin: isMobile
        ? null
        : EdgeInsets.symmetric(
            horizontal: (screenWidth - 600) / 2, vertical: 10),
  );
}

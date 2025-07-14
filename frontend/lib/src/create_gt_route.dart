import 'package:flutter/material.dart';

Route createGtRoute(
  BuildContext context,
  Widget page, {
  emergeVertically = false,
}) {
  const animationDuration = 300;

  return PageRouteBuilder(
    pageBuilder: (context, animation, secondaryAnimation) => page,
    transitionDuration: const Duration(milliseconds: animationDuration),
    reverseTransitionDuration: const Duration(milliseconds: animationDuration),
    transitionsBuilder: (context, animation, secondaryAnimation, child) {
      final begin =
          emergeVertically ? const Offset(0.0, 1.0) : const Offset(-1.0, 0.0);
      const end = Offset.zero;
      const curve = Curves.easeInOutCubic;

      var tween = Tween(begin: begin, end: end).chain(CurveTween(curve: curve));
      var offsetAnimation = animation.drive(tween);

      return SlideTransition(
        position: offsetAnimation,
        child: child,
      );
    },
  );
}

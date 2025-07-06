import 'package:flutter/material.dart';
import 'package:go_todo/globals.dart';

class GtSmallWidthContainer extends StatelessWidget {
  const GtSmallWidthContainer({super.key, required this.child});
  final Widget child;

  @override
  Widget build(BuildContext context) {
    var screenWidth = MediaQuery.sizeOf(context).width;

    return Center(
      child: SizedBox(
        width: screenWidth < ScreenSize.small.value
            ? double.infinity
            : ScreenSize.small.value.toDouble(),
        child: child,
      ),
    );
  }
}

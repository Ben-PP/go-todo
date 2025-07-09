import 'package:flutter/material.dart';

const circleSize = 200.0;

class GtLoadingPage extends StatelessWidget {
  const GtLoadingPage({super.key});

  @override
  Widget build(BuildContext context) {
    return const Center(
      child: SizedBox(
        width: circleSize,
        height: circleSize,
        child: CircularProgressIndicator(),
      ),
    );
  }
}

import 'package:flutter/material.dart';

class GtFadingScrollView extends StatelessWidget {
  const GtFadingScrollView({super.key, required this.children, this.title});
  final List<Widget> children;
  final String? title;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final screenHeight = MediaQuery.sizeOf(context).height;
    const pixelOffset = 20.0;
    final topStop = pixelOffset / screenHeight;
    final bottomStop = 1.0 - topStop;
    const topBottomWhiteSpace = SizedBox(height: 15);

    return Column(
      children: [
        if (title != null)
          Padding(
            padding: const EdgeInsets.only(bottom: 0.0),
            child: Text(
              title!,
              style: Theme.of(context).textTheme.headlineSmall,
            ),
          ),
        Expanded(
          child: ShaderMask(
            shaderCallback: (Rect rect) {
              return LinearGradient(
                begin: Alignment.topCenter,
                end: Alignment.bottomCenter,
                colors: [
                  colorScheme.primaryContainer,
                  Colors.transparent,
                  Colors.transparent,
                  colorScheme.surface,
                ],
                stops: [0.0, topStop, bottomStop, 1.0],
              ).createShader(rect);
            },
            blendMode: BlendMode.dstOut,
            child: SingleChildScrollView(
              child: Column(
                children: [
                  topBottomWhiteSpace,
                  ...children,
                  topBottomWhiteSpace,
                ],
              ),
            ),
          ),
        ),
      ],
    );
  }
}

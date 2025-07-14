import 'package:flutter/material.dart';

class GtFadingScrollView extends StatelessWidget {
  const GtFadingScrollView({
    super.key,
    required this.children,
    this.title,
    this.subtitle,
    this.actions,
    this.crossAxisAlignment = CrossAxisAlignment.center,
  });
  final List<Widget> children;
  final String? title;
  final String? subtitle;
  final CrossAxisAlignment crossAxisAlignment;
  final List<Widget>? actions;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;
    final screenHeight = MediaQuery.sizeOf(context).height;
    const pixelOffset = 20.0;
    final topStop = pixelOffset / screenHeight;
    final bottomStop = 1.0 - topStop;
    const topBottomWhiteSpace = SizedBox(height: 15);

    return Column(
      crossAxisAlignment: crossAxisAlignment,
      children: [
        Padding(
          padding: const EdgeInsets.only(left: 4.0),
          child: Row(
            mainAxisAlignment: actions == null
                ? MainAxisAlignment.center
                : MainAxisAlignment.spaceBetween,
            children: [
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (title != null)
                    Text(title!, style: textTheme.headlineSmall),
                  if (subtitle != null)
                    Text(subtitle!, style: textTheme.bodySmall),
                ],
              ),
              if (actions != null)
                FittedBox(
                  fit: BoxFit.scaleDown,
                  child: Row(
                    children: actions!,
                  ),
                )
            ],
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
                crossAxisAlignment: CrossAxisAlignment.center,
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

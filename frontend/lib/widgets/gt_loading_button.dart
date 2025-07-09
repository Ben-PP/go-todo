import 'package:flutter/material.dart';

class GtLoadingButton extends StatelessWidget {
  const GtLoadingButton({
    super.key,
    required this.text,
    required this.onPressed,
    this.isLoading = false,
    this.isDisabled = false,
  });
  final String text;
  final bool isLoading;
  final VoidCallback onPressed;
  final bool isDisabled;

  @override
  Widget build(BuildContext context) {
    return ElevatedButton(
      onPressed: !isDisabled ? onPressed : null,
      style: isDisabled
          ? Theme.of(context).elevatedButtonTheme.style!.copyWith(
              backgroundColor:
                  WidgetStateProperty.all(Colors.cyan.shade900.withAlpha(80)))
          : null,
      child: isLoading ? const CircularProgressIndicator() : Text(text),
    );
  }
}

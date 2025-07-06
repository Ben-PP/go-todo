import 'package:flutter/material.dart';

class GtTextField extends StatefulWidget {
  const GtTextField({
    super.key,
    required this.controller,
    this.onChanged,
    this.keyboardType = TextInputType.text,
    this.filled = false,
    this.label,
    this.hint,
    this.isSecret = false,
    this.leading,
    this.trailing,
  });
  final TextEditingController controller;
  final ValueChanged<String>? onChanged;
  final TextInputType keyboardType;
  final bool filled;
  final String? label;
  final String? hint;
  final bool isSecret;
  final Widget? trailing;
  final Widget? leading;

  @override
  State<GtTextField> createState() => _GtTextFieldState();
}

class _GtTextFieldState extends State<GtTextField> {
  var showText = false;
  @override
  Widget build(BuildContext context) {
    return TextField(
      controller: widget.controller,
      decoration: InputDecoration(
        border: !widget.filled ? const OutlineInputBorder() : null,
        filled: widget.filled,
        fillColor: widget.filled
            ? Theme.of(context).colorScheme.primaryContainer
            : null,
        labelText: widget.label,
        labelStyle: Theme.of(context).textTheme.labelMedium,
        floatingLabelStyle: Theme.of(context).textTheme.labelMedium,
        hintText: widget.hint,
        hintStyle: Theme.of(context).textTheme.labelMedium?.copyWith(
              color: Theme.of(context).colorScheme.onSurface.withAlpha(120),
            ),
        suffixIcon: widget.isSecret
            ? IconButton(
                onPressed: () => setState(() => showText = !showText),
                icon: showText
                    ? const Icon(Icons.visibility_off)
                    : const Icon(Icons.visibility),
              )
            : widget.trailing,
        prefixIcon: widget.leading,
      ),
      obscureText: widget.isSecret ? !showText : false,
      autocorrect: !widget.isSecret,
      enableSuggestions: !widget.isSecret,
      onChanged: widget.onChanged,
      keyboardType:
          widget.isSecret ? TextInputType.visiblePassword : widget.keyboardType,
    );
  }
}

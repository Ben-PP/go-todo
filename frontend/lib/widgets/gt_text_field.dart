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
    this.autovalidateMode = AutovalidateMode.disabled,
    this.validator,
    this.textInputAction,
    this.maxLength,
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
  final AutovalidateMode autovalidateMode;
  final String? Function(String?)? validator;
  final TextInputAction? textInputAction;
  final int? maxLength;

  @override
  State<GtTextField> createState() => _GtTextFieldState();
}

class _GtTextFieldState extends State<GtTextField> {
  final focusNode = FocusNode(skipTraversal: true);
  var showText = false;

  @override
  void dispose() {
    super.dispose();
    focusNode.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return TextFormField(
      maxLength: widget.maxLength,
      onTapOutside: (event) => FocusScope.of(context).unfocus(),
      autovalidateMode: widget.autovalidateMode,
      validator: widget.validator,
      controller: widget.controller,
      textInputAction: widget.textInputAction,
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
                focusNode: focusNode,
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

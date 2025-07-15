import 'package:flutter/material.dart';

class GtCard extends StatelessWidget {
  const GtCard({
    super.key,
    required this.title,
    this.subtitle,
    this.trailing,
    this.leading,
    this.isSelected = false,
    this.onTap,
  });
  final Widget? trailing;
  final Widget? leading;
  final String title;
  final String? subtitle;
  final bool isSelected;
  final VoidCallback? onTap;

  @override
  Widget build(BuildContext context) {
    final colorScheme = Theme.of(context).colorScheme;
    final textTheme = Theme.of(context).textTheme;

    return Card(
      clipBehavior: Clip.antiAlias,
      color: colorScheme.primaryContainer,
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(12),
        side: isSelected
            ? BorderSide(color: colorScheme.secondary, width: 2)
            : BorderSide.none,
      ),
      child: Material(
        type: MaterialType.transparency,
        child: InkWell(
          onTap: onTap,
          hoverColor: colorScheme.secondary.withAlpha(50),
          splashColor: colorScheme.secondary,
          child: Padding(
            padding: const EdgeInsets.all(12.0),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                if (leading != null) leading!,
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Padding(
                        padding: const EdgeInsets.only(bottom: 8.0),
                        child: Text(
                          title,
                          style: textTheme.labelLarge,
                        ),
                      ),
                      Text(
                        subtitle ?? '',
                        style: textTheme.bodySmall,
                      ),
                    ],
                  ),
                ),
                if (trailing != null) trailing!,
              ],
            ),
          ),
        ),
      ),
    );
  }
}

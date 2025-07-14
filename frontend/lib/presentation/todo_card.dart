import 'package:flutter/material.dart';
import '../domain/todo.dart';

class TodoCard extends StatefulWidget {
  const TodoCard({super.key, required this.todo});
  final Todo todo;

  @override
  State<TodoCard> createState() => _TodoCardState();
}

class _TodoCardState extends State<TodoCard> {
  var confirmDismiss = false;

  @override
  Widget build(BuildContext context) {
    final textTheme = Theme.of(context).textTheme;
    final colorScheme = Theme.of(context).colorScheme;
    return SizedBox(
      child: Card(
        clipBehavior: Clip.antiAlias,
        color: colorScheme.primaryContainer,
        child: ListTile(
          onLongPress: () {
            // TODO Open edit options
            // Delete or pin
          },
          onTap: () {
            // TODO Open todo list route
            // On mobile this should open a new route
            // On desktop this could be opened in side panel
          },
          hoverColor: colorScheme.secondary.withAlpha(50),
          splashColor: colorScheme.secondary,
          title: Text(
            widget.todo.title,
            style: textTheme.labelLarge,
          ),
          subtitle: Text(
            widget.todo.description,
            style: textTheme.bodySmall,
          ),
        ),
      ),
    );
  }
}

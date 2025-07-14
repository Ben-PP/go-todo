import 'package:flutter/material.dart';
import '../domain/todo_list.dart';

class TodoListCard extends StatefulWidget {
  const TodoListCard(
      {super.key, required this.list, this.onTap, this.isSelected = false});
  final TodoList list;
  final VoidCallback? onTap;
  final bool isSelected;

  @override
  State<TodoListCard> createState() => _TodoListCardState();
}

class _TodoListCardState extends State<TodoListCard> {
  var confirmDismiss = false;

  @override
  Widget build(BuildContext context) {
    final doneCount =
        widget.list.todos.where((todo) => todo.isCompleted).length;
    final totalCount = widget.list.todos.length;
    final textTheme = Theme.of(context).textTheme;
    final colorScheme = Theme.of(context).colorScheme;
    return SizedBox(
      child: Card(
        clipBehavior: Clip.antiAlias,
        color: colorScheme.primaryContainer,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: widget.isSelected
              ? BorderSide(color: colorScheme.secondary, width: 2)
              : BorderSide.none,
        ),
        child: ListTile(
          onLongPress: () {
            // TODO Open edit options
            // Delete or pin
          },
          onTap: widget.onTap,
          hoverColor: colorScheme.secondary.withAlpha(50),
          splashColor: colorScheme.secondary,
          title: Text(
            widget.list.title,
            style: textTheme.labelLarge,
          ),
          subtitle: Text(
            widget.list.description ?? '',
            style: textTheme.bodySmall,
          ),
          trailing: Text(
            '$doneCount/$totalCount',
            style: textTheme.labelLarge,
          ),
        ),
      ),
    );
  }
}

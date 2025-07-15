import 'package:flutter/material.dart';
import 'package:go_todo/widgets/gt_card.dart';
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
    return SizedBox(
      child: GtCard(
        title: widget.list.title,
        subtitle: widget.list.description,
        isSelected: widget.isSelected,
        onTap: widget.onTap,
        trailing: Text(
          '$doneCount/$totalCount',
          style: textTheme.labelLarge,
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import '../domain/todo_list.dart';
import '../presentation/todo_card.dart';
import '../widgets/gt_fading_scroll_view.dart';

class TodoView extends StatelessWidget {
  const TodoView({super.key, required this.todoList});
  final TodoList todoList;

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: double.infinity,
      child: GtFadingScrollView(
        title: todoList.title,
        children: [...todoList.todos.map((todo) => TodoCard(todo: todo))],
      ),
    );
  }
}

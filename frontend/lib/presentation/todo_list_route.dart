import 'package:flutter/material.dart';

import '../domain/todo_list.dart';
import '../presentation/route_scaffold.dart';
import '../presentation/todo_view.dart';

class TodoListRoute extends StatelessWidget {
  const TodoListRoute({super.key, required this.todoList});
  final TodoList todoList;

  @override
  Widget build(BuildContext context) {
    return RouteScaffold(
      implyLeading: true,
      showDrawer: false,
      body: TodoView(
        todoList: todoList,
        afterDelete: () => Navigator.of(context).pop(),
      ),
    );
  }
}

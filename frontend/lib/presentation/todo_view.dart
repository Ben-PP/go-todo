import 'package:flutter/material.dart';
import '../domain/todo_list.dart';
import '../widgets/gt_card.dart';
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
        subtitle: todoList.description,
        actions: [
          MenuAnchor(
            builder: (context, controller, child) {
              return IconButton(
                  onPressed: () {
                    if (controller.isOpen) {
                      controller.close();
                    } else {
                      controller.open();
                    }
                  },
                  icon: const Icon(Icons.more_vert));
            },
            menuChildren: [
              MenuItemButton(
                onPressed: () {
                  // Handle edit action
                },
                child: const Text('Edit'),
              ),
              MenuItemButton(
                onPressed: () {
                  // Handle delete action
                },
                child: const Text('Delete'),
              ),
            ],
          ),
        ],
        children: [
          ...todoList.todos.map((todo) => GtCard(
                title: todo.title,
                subtitle: todo.description,
                isSelected: false,
                trailing: IconButton(
                  onPressed: () {},
                  icon: Icon(
                    todo.isCompleted
                        ? Icons.check_box_outlined
                        : Icons.check_box_outline_blank,
                  ),
                ),
              ))
        ],
      ),
    );
  }
}

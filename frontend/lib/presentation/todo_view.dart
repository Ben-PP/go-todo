import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/todo_list.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import '../domain/todo_list.dart' as todo_list_domain;
import '../widgets/gt_card.dart';
import '../widgets/gt_fading_scroll_view.dart';

class TodoView extends ConsumerWidget {
  const TodoView({
    super.key,
    required this.todoList,
    required this.afterDelete,
  });
  final todo_list_domain.TodoList todoList;
  final VoidCallback afterDelete;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
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
              // Delete action
              MenuItemButton(
                onPressed: () async {
                  bool success = await showDialog(
                        context: context,
                        builder: (context) {
                          return AlertDialog(
                            title: const Text('Delete List'),
                            content: const Text(
                                'Are you sure you want to delete this list?'),
                            actions: [
                              TextButton(
                                onPressed: () =>
                                    Navigator.of(context).pop(false),
                                child: const Text('Cancel'),
                              ),
                              TextButton(
                                onPressed: () async {
                                  try {
                                    await ref
                                        .read(todoListProvider.notifier)
                                        .deleteList(todoList.id);
                                    if (context.mounted) {
                                      final snackBar = getSnackBar(
                                        context: context,
                                        content: const Text('List deleted.'),
                                      );
                                      ScaffoldMessenger.of(context)
                                          .clearSnackBars();
                                      ScaffoldMessenger.of(context)
                                          .showSnackBar(snackBar);
                                      Navigator.of(context).pop(true);
                                    }
                                  } on GtApiException catch (_) {
                                    if (context.mounted) {
                                      final snackBar = getSnackBar(
                                        context: context,
                                        content:
                                            const Text('Failed to delete list'),
                                        isError: true,
                                      );
                                      ScaffoldMessenger.of(context)
                                          .clearSnackBars();
                                      ScaffoldMessenger.of(context)
                                          .showSnackBar(snackBar);
                                    }
                                  }
                                },
                                child: const Text('Delete'),
                              ),
                            ],
                          );
                        },
                      ) ??
                      false;
                  if (success) afterDelete();
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

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/presentation/create_list_route.dart';
import 'package:go_todo/widgets/gt_fading_scroll_view.dart';

import '../application/todo_list.dart';
import '../domain/todo_list.dart' as todo_list_domain;
import '../globals.dart';
import '../presentation/todo_list_card.dart';
import '../presentation/todo_list_route.dart';
import '../presentation/todo_view.dart';
import '../src/create_gt_route.dart';
import '../widgets/gt_loading_button.dart';
import '../widgets/gt_loading_page.dart';

class TodoListsView extends ConsumerStatefulWidget {
  const TodoListsView({super.key});

  @override
  ConsumerState<TodoListsView> createState() => _TodoListsViewState();
}

class _TodoListsViewState extends ConsumerState<TodoListsView> {
  var isRefreshing = false;
  String? selectedListId;

  @override
  Widget build(BuildContext context) {
    final AsyncValue<List<todo_list_domain.TodoList>> todoLists =
        ref.watch(todoListProvider);
    final colorScheme = Theme.of(context).colorScheme;
    final isDesktop = MediaQuery.sizeOf(context).width > ScreenSize.large.value;
    Widget content = const GtLoadingPage();

    if (isRefreshing || todoLists.isRefreshing || todoLists is AsyncLoading) {
      return content;
    }

    switch (todoLists) {
      case AsyncLoading():
        content = const GtLoadingPage();
        break;
      case AsyncData(:final value):
        if (isDesktop) {
          setState(() {
            selectedListId ??= value.isNotEmpty ? value.first.id : null;
          });
        }
        content = Center(
          child: SizedBox(
            width: MediaQuery.sizeOf(context).width > ScreenSize.large.value
                ? ScreenSize.large.value.toDouble()
                : double.infinity,
            child: Row(
              mainAxisAlignment: MainAxisAlignment.center,
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Flexible(
                  flex: 4,
                  child: GtFadingScrollView(
                    title: isDesktop ? 'Todo Lists' : null,
                    subtitle: isDesktop ? 'Select a list to view' : null,
                    children: [
                      ...value.map((list) => TodoListCard(
                            list: list,
                            onTap: () {
                              if (isDesktop) {
                                setState(() {
                                  selectedListId = list.id;
                                });
                              } else {
                                Navigator.of(context).push(createGtRoute(
                                  context,
                                  TodoListRoute(todoList: list),
                                  emergeVertically: true,
                                ));
                              }
                            },
                            isSelected: selectedListId == list.id && isDesktop,
                          )),
                      SizedBox(
                        width: double.infinity,
                        child: Card(
                          shape: RoundedRectangleBorder(
                            borderRadius: BorderRadius.circular(12),
                            side: BorderSide(
                              color: colorScheme.secondary,
                              width: 4,
                            ),
                          ),
                          clipBehavior: Clip.antiAlias,
                          child: InkWell(
                            hoverColor: colorScheme.secondary.withAlpha(50),
                            splashColor: colorScheme.secondary,
                            onTap: () {
                              Navigator.of(context).push(createGtRoute(
                                context,
                                const CreateListRoute(),
                                emergeVertically: true,
                              ));
                            },
                            child: Padding(
                              padding: const EdgeInsets.all(16.0),
                              child: Row(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Text(
                                    'Add New List',
                                    style:
                                        Theme.of(context).textTheme.labelLarge,
                                  ),
                                  Icon(
                                    Icons.add,
                                    size: 32,
                                    color: colorScheme.onSurface,
                                  ),
                                ],
                              ),
                            ),
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
                if (isDesktop)
                  VerticalDivider(
                    width: 16,
                    thickness: 2,
                    color: Theme.of(context).colorScheme.primaryContainer,
                  ),
                if (isDesktop)
                  Flexible(
                    flex: 8,
                    child: TodoView(
                      todoList: value.firstWhere((l) {
                        if (selectedListId != null) {
                          return l.id == selectedListId;
                        }
                        return l.todos.isNotEmpty;
                      }),
                    ),
                  ),
              ],
            ),
          ),
        );
        break;
      case AsyncError():
        content = SizedBox(
          width: double.infinity,
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              Padding(
                padding: const EdgeInsets.symmetric(vertical: 20.0),
                child: Text(
                  'Error loading todo lists.',
                  style: Theme.of(context).textTheme.headlineSmall,
                ),
              ),
              GtLoadingButton(
                  text: 'Try Again',
                  onPressed: () {
                    var _ = ref.refresh(todoListProvider);
                  }),
            ],
          ),
        );
        break;
    }
    return content;
  }
}

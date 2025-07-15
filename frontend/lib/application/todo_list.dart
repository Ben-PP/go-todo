import 'package:go_todo/data/gt_api.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../domain/todo_list.dart' as todo_list_domain;

part 'todo_list.g.dart';

@riverpod
class TodoList extends _$TodoList {
  @override
  Future<List<todo_list_domain.TodoList>> build() async {
    await Future.delayed(const Duration(milliseconds: 200));
    var todoLists = await GtApi().getLists();
    return todoLists;
  }

  Future<void> createList({
    required String title,
    String? description,
  }) async {
    final newList = await GtApi().createList(
      title: title,
      description: description,
    );
    state = AsyncData([...state.value ?? [], newList]);
  }

  Future<void> deleteList(String listId) async {
    await GtApi().deleteList(listId);
    state = AsyncData(
      (state.value ?? []).where((list) => list.id != listId).toList(),
    );
  }
}

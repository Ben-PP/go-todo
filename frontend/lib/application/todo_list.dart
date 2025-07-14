import 'package:go_todo/data/gt_api.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

import '../domain/todo_list.dart';

part 'todo_list.g.dart';

@Riverpod()
Future<List<TodoList>> todoLists(Ref ref) async {
  await Future.delayed(const Duration(milliseconds: 200));
  var todoLists = await GtApi().getLists();
  return todoLists;
}

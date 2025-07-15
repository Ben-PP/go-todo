// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'todo_list.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

@ProviderFor(TodoList)
const todoListProvider = TodoListProvider._();

final class TodoListProvider
    extends $AsyncNotifierProvider<TodoList, List<todo_list_domain.TodoList>> {
  const TodoListProvider._()
      : super(
          from: null,
          argument: null,
          retry: null,
          name: r'todoListProvider',
          isAutoDispose: true,
          dependencies: null,
          $allTransitiveDependencies: null,
        );

  @override
  String debugGetCreateSourceHash() => _$todoListHash();

  @$internal
  @override
  TodoList create() => TodoList();
}

String _$todoListHash() => r'd272ab1d86a4cbe8d7a5b954575b13135cb7edef';

abstract class _$TodoList
    extends $AsyncNotifier<List<todo_list_domain.TodoList>> {
  FutureOr<List<todo_list_domain.TodoList>> build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<AsyncValue<List<todo_list_domain.TodoList>>,
        List<todo_list_domain.TodoList>>;
    final element = ref.element as $ClassProviderElement<
        AnyNotifier<AsyncValue<List<todo_list_domain.TodoList>>,
            List<todo_list_domain.TodoList>>,
        AsyncValue<List<todo_list_domain.TodoList>>,
        Object?,
        Object?>;
    element.handleValue(ref, created);
  }
}

// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package

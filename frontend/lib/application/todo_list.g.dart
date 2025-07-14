// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'todo_list.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

@ProviderFor(todoLists)
const todoListsProvider = TodoListsProvider._();

final class TodoListsProvider extends $FunctionalProvider<
        AsyncValue<List<TodoList>>, List<TodoList>, FutureOr<List<TodoList>>>
    with $FutureModifier<List<TodoList>>, $FutureProvider<List<TodoList>> {
  const TodoListsProvider._()
      : super(
          from: null,
          argument: null,
          retry: null,
          name: r'todoListsProvider',
          isAutoDispose: true,
          dependencies: null,
          $allTransitiveDependencies: null,
        );

  @override
  String debugGetCreateSourceHash() => _$todoListsHash();

  @$internal
  @override
  $FutureProviderElement<List<TodoList>> $createElement(
          $ProviderPointer pointer) =>
      $FutureProviderElement(pointer);

  @override
  FutureOr<List<TodoList>> create(Ref ref) {
    return todoLists(ref);
  }
}

String _$todoListsHash() => r'747af123309affaa27e4941b7d5dc1b2738b5db0';

// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package

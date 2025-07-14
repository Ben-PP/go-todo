import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/todo_list.dart';

import '../application/authentication_provider.dart';
import '../widgets/gt_loading_page.dart';
import './auth_view/login_view.dart';
import './auth_view/retry_refresh_view.dart';
import 'todo_lists_view.dart';
import './route_scaffold.dart';

class HomeRoute extends ConsumerWidget {
  const HomeRoute({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    var authState = ref.watch(authenticationProvider);
    var hasAuthError = authState == AuthState.error;
    if (hasAuthError) {
      return const RetryRefreshRoute();
    } else if (authState == AuthState.pending) {
      return const RouteScaffold(body: GtLoadingPage());
    }
    return RouteScaffold(
      appBarActions: [
        if (authState == AuthState.authenticated)
          IconButton(
              onPressed: () {
                ref.invalidate(todoListsProvider);
              },
              icon: const Icon(Icons.refresh)),
        if (authState == AuthState.authenticated)
          IconButton(
              onPressed: () {
                // TODO open create list route
              },
              icon: const Icon(Icons.add)),
      ],
      body: authState != AuthState.unauthenticated
          ? const TodoListsView()
          : const LoginView(),
    );
  }
}

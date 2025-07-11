import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/presentation/auth_view/login_view.dart';
import 'package:go_todo/presentation/auth_view/retry_refresh_view.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/widgets/gt_loading_page.dart';

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
      body: authState != AuthState.unauthenticated
          ? const Center(
              child: Text(
                'Home Route',
                style: TextStyle(fontSize: 24),
              ),
            )
          : const LoginView(),
    );
  }
}

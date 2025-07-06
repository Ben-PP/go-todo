import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/presentation/login_view.dart';
import 'package:go_todo/presentation/route_scaffold.dart';

class HomeRoute extends ConsumerWidget {
  const HomeRoute({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return RouteScaffold(
      appBarActions:
          ref.watch(authenticationProvider) == AuthenticationState.initial
              ? [
                  IconButton(
                      onPressed: () {
                        // TODO Allow to change API url
                      },
                      icon: Icon(Icons.settings))
                ]
              : null,
      body: ref.watch(authenticationProvider) == AuthenticationState.success
          ? const Center(
              child: Text(
                'Home Route',
                style: TextStyle(fontSize: 24),
              ),
            )
          : LoginView(),
    );
  }
}

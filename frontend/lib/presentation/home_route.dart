import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/presentation/auth_view/login_view.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/widgets/gt_loading_button.dart';
import 'package:go_todo/widgets/gt_loading_page.dart';
import 'package:go_todo/widgets/gt_small_width_container.dart';

class HomeRoute extends ConsumerStatefulWidget {
  const HomeRoute({super.key});

  @override
  ConsumerState<HomeRoute> createState() => _HomeRouteState();
}

class _HomeRouteState extends ConsumerState<HomeRoute> {
  @override
  Widget build(BuildContext context) {
    var authState = ref.watch(authenticationProvider);
    var hasAuthError = authState == AuthState.error;
    var isRetrying = false; // TODO Move this
    if (hasAuthError) {
      return RouteScaffold(
          body: Center(
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            const Padding(
              padding: EdgeInsets.only(bottom: 20.0),
              child: Text('There was error authenticating...'),
            ),
            GtSmallWidthContainer(
              child: SizedBox(
                width: double.infinity,
                child: GtLoadingButton(
                  isLoading: isRetrying,
                  onPressed: () async {
                    setState(() {
                      isRetrying = true;
                    });
                    await ref.read(authenticationProvider.notifier).refresh();
                    setState(() {
                      isRetrying = false;
                    });
                  },
                  text: 'Retry',
                ),
              ),
            )
          ],
        ),
      ));
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

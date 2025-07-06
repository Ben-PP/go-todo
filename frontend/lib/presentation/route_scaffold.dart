import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';

class RouteScaffold extends ConsumerWidget {
  const RouteScaffold(
      {super.key,
      this.title,
      required this.body,
      this.bottomNavigationBar,
      this.appBarActions});
  final Widget? title;
  final Widget body;
  final Widget? bottomNavigationBar;
  final List<Widget>? appBarActions;

  Widget buildDrawerEntry(
      BuildContext context, String text, VoidCallback onPressed) {
    return TextButton(
      onPressed: onPressed,
      child: SizedBox(
        width: double.infinity,
        child: Text(
          text,
          style: Theme.of(context).textTheme.labelLarge,
          textAlign: TextAlign.start,
        ),
      ),
    );
  }

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: title ?? const Text('Go Todo'),
        actions: appBarActions,
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(10.0),
          child: body,
        ),
      ),
      drawer: ref.watch(authenticationProvider) == AuthenticationState.success
          ? Drawer(
              child: SafeArea(
                child: Padding(
                  padding: const EdgeInsets.all(10.0),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            'GO-TODO',
                            style: Theme.of(context).textTheme.headlineSmall,
                          ),
                          IconButton(
                            onPressed: () => Navigator.pop(context),
                            icon: const Icon(Icons.arrow_back),
                          ),
                        ],
                      ),
                      const Padding(
                        padding: EdgeInsets.fromLTRB(0, 8, 80, 8),
                        child: Divider(),
                      ),
                      buildDrawerEntry(context, 'Admin', () {}),
                      const Padding(
                        padding: EdgeInsets.fromLTRB(0, 8, 80, 8),
                        child: Divider(),
                      ),
                      buildDrawerEntry(context, 'Logout', () async {
                        await ref
                            .read(authenticationProvider.notifier)
                            .logout();
                        if (context.mounted) {
                          Navigator.pop(context);
                        }
                      }),
                    ],
                  ),
                ),
              ),
            )
          : null,
      bottomNavigationBar: bottomNavigationBar,
    );
  }
}

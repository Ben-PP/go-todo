import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/presentation/gt_drawer.dart';

class RouteScaffold extends ConsumerWidget {
  const RouteScaffold({
    super.key,
    this.title,
    required this.body,
    this.bottomNavigationBar,
    this.appBarActions,
    this.showDrawer = true,
    this.implyLeading = true,
  });
  final Widget? title;
  final Widget body;
  final Widget? bottomNavigationBar;
  final List<Widget>? appBarActions;
  final bool showDrawer;
  final bool implyLeading;

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return Scaffold(
      appBar: AppBar(
        title: title ?? const Text('Go Todo'),
        actions: appBarActions,
        automaticallyImplyLeading: implyLeading,
      ),
      body: SafeArea(
        child: Padding(
          padding: const EdgeInsets.all(15.0),
          child: body,
        ),
      ),
      drawer: showDrawer ? const GtDrawer() : null,
      bottomNavigationBar: bottomNavigationBar,
    );
  }
}

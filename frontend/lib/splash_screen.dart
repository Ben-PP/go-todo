import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/api_url_route.dart';
import 'package:go_todo/presentation/home_route.dart';
import 'package:go_todo/presentation/route_scaffold.dart';
import 'package:go_todo/src/create_gt_route.dart';
import 'package:go_todo/widgets/gt_loading_page.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  late final Future<void> initFuture;

  @override
  void initState() {
    super.initState();
    initFuture = _redirect();

    // Set up periodic JWT refresh
    Timer.periodic(const Duration(minutes: 25), (timer) {
      if (GtApi().refreshJWT != null) {
        ref.read(authenticationProvider.notifier).refresh();
      }
    });
  }

  Future<void> _redirect() async {
    await Future.delayed(const Duration(seconds: 0));
    await GtApi().init();

    WidgetsBinding.instance.addPostFrameCallback((_) async {
      if (GtApi().baseUrl == null) {
        Navigator.of(context).push(createGtRoute(context, const ApiUrlRoute()));
      } else {
        await ref.read(authenticationProvider.notifier).refresh();
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder(
        future: initFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting &&
              ref.watch(authenticationProvider) == AuthState.pending) {
            return const RouteScaffold(
              body: GtLoadingPage(),
            );
          }
          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }
          return const HomeRoute();
        });
  }
}

import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/application/authentication_provider.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/api_url_route.dart';
import 'package:go_todo/presentation/home_route.dart';
import 'package:go_todo/presentation/route_scaffold.dart';

class SplashScreen extends ConsumerStatefulWidget {
  const SplashScreen({super.key});

  @override
  ConsumerState<SplashScreen> createState() => _SplashScreenState();
}

class _SplashScreenState extends ConsumerState<SplashScreen> {
  var hasBaseUrl = GtApi().baseUrl != null;

  late final Future<void> initFuture;

  @override
  void initState() {
    super.initState();
    initFuture = _redirect();
  }

  Future<void> _redirect() async {
    await Future.delayed(const Duration(seconds: 0));
    await GtApi().init();
    var baseUrl = GtApi().baseUrl;

    if (baseUrl != null) {
      setState(() {
        hasBaseUrl = true;
      });
    }
    if (GtApi().refreshJWT != null) {
      await ref.read(authenticationProvider.notifier).refresh();
    }
    // Set up periodic JWT refresh
    Timer.periodic(const Duration(minutes: 25), (timer) async {
      if (GtApi().refreshJWT != null) {
        ref.read(authenticationProvider.notifier).refresh();
        //await GtApi().refresh();
      }
    });
  }

  toggleHasBaseUrl() {
    setState(() {
      hasBaseUrl = !hasBaseUrl;
    });
  }

  Widget buildHomeRoute() {
    if (!hasBaseUrl) {
      return ApiUrlRoute(onSuccess: toggleHasBaseUrl);
    }
    return const HomeRoute();
  }

  @override
  Widget build(BuildContext context) {
    const circleSize = 200.0;

    return FutureBuilder(
        future: initFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.waiting) {
            return const RouteScaffold(
              body: Center(
                child: SizedBox(
                  width: circleSize,
                  height: circleSize,
                  child: CircularProgressIndicator(),
                ),
              ),
            );
          }
          if (snapshot.hasError) {
            return Center(child: Text('Error: ${snapshot.error}'));
          }
          return buildHomeRoute();
        });
  }
}

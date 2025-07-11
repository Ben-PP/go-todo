import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/presentation/route_scaffold.dart';

import '../../application/authentication_provider.dart';
import '../../widgets/gt_loading_button.dart';
import '../../widgets/gt_small_width_container.dart';

class RetryRefreshRoute extends ConsumerStatefulWidget {
  const RetryRefreshRoute({super.key});

  @override
  ConsumerState<RetryRefreshRoute> createState() => _RetryRefreshViewState();
}

class _RetryRefreshViewState extends ConsumerState<RetryRefreshRoute> {
  var isRetrying = false;
  @override
  Widget build(BuildContext context) {
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
  }
}

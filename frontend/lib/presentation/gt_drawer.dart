import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/presentation/api_url_route.dart';
import 'package:go_todo/src/create_gt_route.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import 'package:logging/logging.dart';

import '../application/authentication_provider.dart';

class GtDrawer extends ConsumerStatefulWidget {
  const GtDrawer({super.key});

  @override
  ConsumerState<GtDrawer> createState() => _GtDrawerState();
}

class _GtDrawerState extends ConsumerState<GtDrawer> {
  var isLoading = false;

  Widget divider() {
    return const Padding(
      padding: EdgeInsets.fromLTRB(0, 8, 80, 8),
      child: Divider(),
    );
  }

  Future<void> logout(BuildContext context) async {
    setState(() {
      isLoading = true;
    });
    var snackMessage = '';
    var isError = false;
    try {
      await ref.read(authenticationProvider.notifier).logout();
      snackMessage = 'Successfully logged out!';
    } on GtApiException catch (error) {
      isError = true;
      switch (error.type) {
        case GtApiExceptionType.malformedBody:
          snackMessage = 'Logout requests body was malformed :D';
          break;
        case GtApiExceptionType.unauthorized:
          snackMessage = 'Failed to logout with invalid credentials.';
          break;
        case GtApiExceptionType.forbidden:
          snackMessage = 'It is forbidden to log other users out!';
          break;
        case GtApiExceptionType.serverError:
          snackMessage = 'Server was not happy with the request...';
          break;
        case GtApiExceptionType.unknown:
          // I know Copilot, the spelling mistake is intentional...
          snackMessage = 'Server responded with unhandled statsu code or smth.';
        case GtApiExceptionType.hostNotResponding:
          snackMessage = 'Server ghosted us. No response.';
        default:
          snackMessage = 'Server response was not handled well.';
      }
    } catch (error) {
      snackMessage = 'There was mysterious exception not handled.';
      isError = true;
      log(
        'Unhandled error during logout',
        error: error,
        level: Level.SEVERE.value,
      );
    } finally {
      setState(() {
        isLoading = false;
      });
    }
    if (!context.mounted) return;
    ScaffoldMessenger.of(context).showSnackBar(getSnackBar(
      context: context,
      content: Text(snackMessage),
      isError: isError,
    ));
    if (!isError) Navigator.pop(context);
  }

  @override
  Widget build(BuildContext context) {
    var isAuthenticated =
        ref.watch(authenticationProvider) != AuthState.unauthenticated;
    return Drawer(
      child: SafeArea(
        child: Padding(
          padding: const EdgeInsets.fromLTRB(15, 15, 15, 30),
          child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
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
              divider(),
              Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  if (isAuthenticated)
                    GtDrawerEntry(text: 'Admin', onPressed: () {}),
                ],
              ),
              Column(
                mainAxisAlignment: MainAxisAlignment.end,
                children: [
                  const GtDrawerTitle(text: 'Settings'),
                  GtDrawerEntry(
                      text: 'Change API Url',
                      onPressed: () {
                        Navigator.of(context).push(
                          createGtRoute(
                            context,
                            const ApiUrlRoute(canCancel: true),
                          ),
                        );
                      }),
                  if (isAuthenticated)
                    GtDrawerEntry(
                      text: 'Logout',
                      onPressed: () async => logout(context),
                      isLoading: isLoading,
                    ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class GtDrawerTitle extends StatelessWidget {
  const GtDrawerTitle({super.key, required this.text});
  final String text;

  Widget divider() {
    return const Padding(
      padding: EdgeInsets.fromLTRB(0, 8, 80, 8),
      child: Divider(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        SizedBox(
          width: double.infinity,
          child: Text(
            text,
            style: Theme.of(context)
                .textTheme
                .titleMedium!
                .copyWith(color: Colors.white54),
            textAlign: TextAlign.start,
          ),
        ),
        divider(),
      ],
    );
  }
}

class GtDrawerEntry extends StatelessWidget {
  const GtDrawerEntry({
    super.key,
    required this.text,
    required this.onPressed,
    this.isLoading = false,
  });
  final String text;
  final VoidCallback onPressed;
  final bool isLoading;

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onPressed,
      child: SizedBox(
        width: double.infinity,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.spaceBetween,
          children: [
            Text(
              text,
              style: Theme.of(context)
                  .textTheme
                  .labelLarge!
                  .copyWith(color: Theme.of(context).colorScheme.secondary),
              textAlign: TextAlign.start,
            ),
            if (isLoading)
              CircularProgressIndicator(
                constraints: BoxConstraints.tight(const Size(15, 15)),
              ),
          ],
        ),
      ),
    );
  }
}

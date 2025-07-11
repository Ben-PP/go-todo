import 'dart:developer';

import 'package:flutter/material.dart';
import 'package:go_todo/data/gt_api.dart';
import 'package:go_todo/src/get_snack_bar.dart';
import 'package:logging/logging.dart';

showErrorSnack(BuildContext context, GtApiException error,
    {Map<GtApiExceptionType, String>? map}) {
  var snackMessage = '';
  switch (error.type) {
    case GtApiExceptionType.forbidden:
      snackMessage = map?[GtApiExceptionType.forbidden] ??
          'You are not allowed to perform this action.';
      break;
    case GtApiExceptionType.hostUnknown:
      snackMessage = map?[GtApiExceptionType.hostUnknown] ?? 'Unknown host.';
      break;
    case GtApiExceptionType.hostNotResponding:
      snackMessage =
          map?[GtApiExceptionType.hostNotResponding] ?? 'Host did not respond';
      break;
    case GtApiExceptionType.conflict:
      snackMessage = map?[GtApiExceptionType.conflict] ??
          'Conflict occurred. This usually means you are trying to create a resource that already exists.';
    case GtApiExceptionType.urlNull:
      snackMessage =
          map?[GtApiExceptionType.urlNull] ?? 'Url has not been set.';
      break;
    case GtApiExceptionType.refreshJWTNull:
      snackMessage = map?[GtApiExceptionType.refreshJWTNull] ??
          'Refresh JWT has not been set. Please login again.';
      break;
    case GtApiExceptionType.unauthorized:
      snackMessage = map?[GtApiExceptionType.unauthorized] ??
          'You are not authorized to perform this action.';
      break;
    case GtApiExceptionType.malformedBody:
      snackMessage =
          map?[GtApiExceptionType.malformedBody] ?? 'Malformed request body.';
      break;
    case GtApiExceptionType.serverError:
      snackMessage = map?[GtApiExceptionType.serverError] ??
          'Server error occurred. Try again later.';
      break;
    case GtApiExceptionType.unknownResponse:
      snackMessage = map?[GtApiExceptionType.unknownResponse] ??
          'Received an unknown response from the server.';
      break;
    case GtApiExceptionType.unknown:
      snackMessage = map?[GtApiExceptionType.unknown] ??
          'An unknown error occurred. Heart breaking, I know.';
      break;
  }
  final snackBar = getSnackBar(
    context: context,
    content: Text(snackMessage),
    isError: true,
  );
  ScaffoldMessenger.of(context).clearSnackBars();
  ScaffoldMessenger.of(context).showSnackBar(snackBar);
  log(
    'Showed snackbar for error.',
    error: error,
    stackTrace: StackTrace.current,
    level: Level.INFO.value,
  );
}

import 'dart:developer';

import 'package:go_todo/data/gt_api.dart';
import 'package:logging/logging.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'authentication_provider.g.dart';

enum AuthState {
  unauthenticated,
  authenticated,
  pending,
  error,
}

@riverpod
class Authentication extends _$Authentication {
  @override
  AuthState build() {
    return AuthState.pending;
  }

  /// Logs user in using [username] and [password].
  ///
  /// Tries to log in with [username] and [password]. Rethrows any errors thrown
  /// from GtApi.login().
  Future<void> login(String username, String password) async {
    try {
      await GtApi().login(username, password);
      state = AuthState.authenticated;
    } on GtApiException catch (_) {
      state = AuthState.unauthenticated;
      rethrow;
    } catch (error) {
      log(
        'Unknown exception in AuthProvider login.',
        error: error,
        level: Level.SEVERE.value,
      );
      rethrow;
    }
  }

  /// Logs the user auth.
  Future<void> logout() async {
    try {
      await GtApi().logout();
      state = AuthState.unauthenticated;
    } catch (error) {
      log('Something went wrong while logging out.',
          error: error, level: Level.SEVERE.value);
      rethrow;
    }
  }

  Future<void> refresh() async {
    try {
      if (GtApi().refreshJWT == null) {
        state = AuthState.unauthenticated;
        return;
      }
      await GtApi().refresh();
      if (GtApi().refreshJWT != null) {
        state = AuthState.authenticated;
      } else {
        state = AuthState.unauthenticated;
      }
    } catch (error) {
      state = AuthState.error;
      log('Failed to refresh tokens', error: error, level: Level.SEVERE.value);
    }
  }
}

import 'package:go_todo/data/gt_api.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';

part 'authentication_provider.g.dart';

enum AuthenticationState {
  initial,
  loading,
  success,
  error,
}

@riverpod
class Authentication extends _$Authentication {
  @override
  AuthenticationState build() {
    return AuthenticationState.initial;
  }

  Future<void> login(String username, String password) async {
    state = AuthenticationState.loading;

    // Call login service
    try {
      await GtApi().login(username, password);
      state = AuthenticationState.success;
    } catch (error) {
      state = AuthenticationState.error;
    }
  }

  Future<void> logout() async {
    try {
      await GtApi().logout();
      state = AuthenticationState.initial;
    } catch (error) {
      // TODO Handle this error in the UI
      state = AuthenticationState.error;
      rethrow;
    }
  }

  Future<void> refresh() async {
    try {
      state = AuthenticationState.success;
    } catch (error) {
      // TODO Handle this error in the UI some how
      state = AuthenticationState.error;
    }
  }
}

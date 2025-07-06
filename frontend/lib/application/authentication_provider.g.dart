// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'authentication_provider.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

@ProviderFor(Authentication)
const authenticationProvider = AuthenticationProvider._();

final class AuthenticationProvider
    extends $NotifierProvider<Authentication, AuthenticationState> {
  const AuthenticationProvider._()
      : super(
          from: null,
          argument: null,
          retry: null,
          name: r'authenticationProvider',
          isAutoDispose: true,
          dependencies: null,
          $allTransitiveDependencies: null,
        );

  @override
  String debugGetCreateSourceHash() => _$authenticationHash();

  @$internal
  @override
  Authentication create() => Authentication();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(AuthenticationState value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<AuthenticationState>(value),
    );
  }
}

String _$authenticationHash() => r'db7a4c9b917b3274866e6ba33222824d484eb3d4';

abstract class _$Authentication extends $Notifier<AuthenticationState> {
  AuthenticationState build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<AuthenticationState, AuthenticationState>;
    final element = ref.element as $ClassProviderElement<
        AnyNotifier<AuthenticationState, AuthenticationState>,
        AuthenticationState,
        Object?,
        Object?>;
    element.handleValue(ref, created);
  }
}

// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package

// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'has_api_url.dart';

// **************************************************************************
// RiverpodGenerator
// **************************************************************************

@ProviderFor(HasApiUrl)
const hasApiUrlProvider = HasApiUrlProvider._();

final class HasApiUrlProvider extends $NotifierProvider<HasApiUrl, bool> {
  const HasApiUrlProvider._()
      : super(
          from: null,
          argument: null,
          retry: null,
          name: r'hasApiUrlProvider',
          isAutoDispose: true,
          dependencies: null,
          $allTransitiveDependencies: null,
        );

  @override
  String debugGetCreateSourceHash() => _$hasApiUrlHash();

  @$internal
  @override
  HasApiUrl create() => HasApiUrl();

  /// {@macro riverpod.override_with_value}
  Override overrideWithValue(bool value) {
    return $ProviderOverride(
      origin: this,
      providerOverride: $SyncValueProvider<bool>(value),
    );
  }
}

String _$hasApiUrlHash() => r'b8e3f36cf4eeecc278851c5649f31757bb931451';

abstract class _$HasApiUrl extends $Notifier<bool> {
  bool build();
  @$mustCallSuper
  @override
  void runBuild() {
    final created = build();
    final ref = this.ref as $Ref<bool, bool>;
    final element = ref.element as $ClassProviderElement<
        AnyNotifier<bool, bool>, bool, Object?, Object?>;
    element.handleValue(ref, created);
  }
}

// ignore_for_file: type=lint
// ignore_for_file: subtype_of_sealed_class, invalid_use_of_internal_member, invalid_use_of_visible_for_testing_member, deprecated_member_use_from_same_package

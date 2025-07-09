import 'package:riverpod_annotation/riverpod_annotation.dart';
//import '../data/gt_api.dart';

part 'has_api_url.g.dart';

@riverpod
class HasApiUrl extends _$HasApiUrl {
  @override
  bool build() {
    return false;
  }

  setState(bool has) {
    state = has;
  }

  /*Future<void> setUrl(String url) async {
    try {
      await GtApi().setBaseUrl(url);
      state = true;
    } catch (error) {
      state = false;
    }
  }

  Future<void> loadUrl() async {
    if (GtApi().baseUrl != null) {
      state = true;
    } else {
      state = false;
    }
  }*/
}

import 'dart:convert';

import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:http/http.dart' as http;
import 'package:shared_preferences/shared_preferences.dart';

class GtApi {
  static final GtApi _instance = GtApi._internal();
  static const String defaultPath = '/api/v1';

  String? baseUrl;
  String? accessJWT;
  String? refreshJWT;

  factory GtApi() {
    return _instance;
  }

  GtApi._internal();

  /// Initialized the GtApi singleton.
  ///
  /// Tries to get the baseUrl and the jwts from shared preferences. Removes the
  /// refresh jwt if it has expired.
  Future<void> init() async {
    final prefs = SharedPreferencesAsync();
    baseUrl = await prefs.getString('baseUrl');
    accessJWT = await prefs.getString('accessJWT');
    refreshJWT = await prefs.getString('refreshJWT');
    if (refreshJWT != null) {
      var test = JWT.decode(refreshJWT!);
      var exp = DateTime.fromMillisecondsSinceEpoch(test.payload['exp'] * 1000);
      if (exp.isBefore(DateTime.now())) {
        await prefs.remove('refreshJWT');
        await prefs.remove('accessJWT');
        refreshJWT = null;
        accessJWT = null;
      } else {
        await refresh();
      }
    }
  }

  /// Check that baseUrl is not null
  _checkBaseUrl() {
    if (baseUrl == null) {
      throw Exception('Base URL is not set');
    }
  }

  /// Set the base url for the applications backend calls
  Future<void> setBaseUrl(String url) async {
    var fullUrl = url.endsWith('/')
        ? url.substring(0, url.length - 1) + defaultPath
        : url + defaultPath;

    var response = await http.get(Uri.parse('$fullUrl/status'));
    if (response.statusCode != 200) {
      throw Exception('Invalid API URL: $url');
    }

    baseUrl = fullUrl;
    final prefs = SharedPreferencesAsync();
    await prefs.setString('baseUrl', fullUrl);
  }

  /// Refresh the tokens
  Future<void> refresh() async {
    _checkBaseUrl();

    var response = await http.post(
      Uri.parse('$baseUrl/auth/refresh'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({'refresh_token': refreshJWT}),
    );

    if (response.statusCode == 200) {
      var data = jsonDecode(response.body);
      accessJWT = data['access_token'];
      refreshJWT = data['refresh_token'];

      final prefs = SharedPreferencesAsync();
      await prefs.setString('accessJWT', accessJWT!);
      await prefs.setString('refreshJWT', refreshJWT!);
    } else {
      throw Exception('Failed to refresh tokens: ${response.body}');
    }
  }

  /// Login with credentials
  ///
  /// Uses [username] and [password] to get tokens form the API.
  Future<void> login(String username, String password) async {
    _checkBaseUrl();
    var response = await http.post(
      Uri.parse('$baseUrl/auth/login'),
      headers: {'Content-Type': 'application/json'},
      body: jsonEncode({
        'username': username,
        'password': password,
      }), //'{"username": "$username", "password": "$password"}',
    );

    if (response.statusCode == 200) {
      var data = jsonDecode(response.body);

      accessJWT = data['access_token'];
      refreshJWT = data['refresh_token'];

      final prefs = SharedPreferencesAsync();
      await prefs.setString('accessJWT', accessJWT!);
      await prefs.setString('refreshJWT', refreshJWT!);
    } else {
      throw Exception('Failed to login: ${response.body}');
    }
  }

  Future<void> logout() async {
    _checkBaseUrl();

    var response = await http.post(Uri.parse('$baseUrl/auth/logout'),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $accessJWT'
        },
        body: jsonEncode({
          'refresh_token': refreshJWT,
        }));

    if (response.statusCode != 204) {
      throw Exception('Failed to logout');
    }

    var prefs = SharedPreferencesAsync();
    await prefs.remove('accessJWT');
    await prefs.remove('refreshJWT');
  }
}

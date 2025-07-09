import 'dart:convert';
import 'dart:developer';
import 'dart:io';

import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:http/http.dart' as http;
import 'package:logging/logging.dart';
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
        if (baseUrl != null) {
          //await refresh();
        }
      }
    }
    //prefs.remove('baseUrl');
    //prefs.remove('refreshJWT');
  }

  /// Check that baseUrl is not null
  bool _hasBaseUrl() {
    if (baseUrl == null) {
      return false;
    }
    return true;
  }

  /// Set the base url for the applications backend calls
  Future<void> setBaseUrl(String url) async {
    var fullUrl = url.endsWith('/')
        ? url.substring(0, url.length - 1) + defaultPath
        : url + defaultPath;
    try {
      var response = await http.get(Uri.parse('$fullUrl/status'));
      if (response.statusCode != 200) {
        throw Exception('Invalid API URL: $url');
      }

      baseUrl = fullUrl;
      final prefs = SharedPreferencesAsync();
      await prefs.setString('baseUrl', fullUrl);
    } on http.ClientException catch (error) {
      log('Failed to connect API.', error: error, level: Level.SEVERE.value);
      if (error.message.contains('Connection refused')) {
        throw GtApiException(
          cause: 'Host refused connection.',
          type: GtApiExceptionType.hostNotResponding,
        );
      }
      throw GtApiException(
        cause: error.toString(),
        type: GtApiExceptionType.hostUnknown,
      );
    } catch (error) {
      log('Failed to connect API.', error: error, level: Level.SEVERE.value);
      throw GtApiException(
        cause: error.toString(),
        type: GtApiExceptionType.unknown,
      );
    }
  }

  /// Refresh the tokens
  Future<void> refresh() async {
    if (!_hasBaseUrl()) {
      final error = GtApiException(
          cause: 'BaseUrl not set.', type: GtApiExceptionType.urlNull);
      log('App has no baseUrl', level: Level.SEVERE.value, error: error);
      throw error;
    }
    if (refreshJWT == null) {
      final error = GtApiException(
        cause: 'No refresh token saved',
        type: GtApiExceptionType.refreshJWTNull,
      );
      log(
        'App has no saved refresh JWT.',
        error: error,
        level: Level.SEVERE.value,
      );
    }
    try {
      var response = await http.post(
        Uri.parse('$baseUrl/auth/refresh'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({'refresh_token': refreshJWT}),
      );
      switch (response.statusCode) {
        case 200:
          var data = jsonDecode(response.body);
          accessJWT = data['access_token'];
          refreshJWT = data['refresh_token'];
          final prefs = SharedPreferencesAsync();
          await prefs.setString('accessJWT', accessJWT!);
          await prefs.setString('refreshJWT', refreshJWT!);
          log('Tokens refreshed', level: Level.INFO.value);
          break;
        case 400:
          var error = GtApiException(
            cause: 'Malformed request: ${response.body}',
            type: GtApiExceptionType.malformedBody,
          );
          log(error.cause, error: error, level: Level.SEVERE.value);
          throw error;
        case 401:
          final prefs = SharedPreferencesAsync();
          await prefs.remove('accessJWT');
          await prefs.remove('refreshJWT');
          final error = GtApiException(
            cause: 'Token was unauthorized for refresh.',
            type: GtApiExceptionType.refreshJWTUnauthorized,
          );
          log(error.cause, error: error, level: Level.INFO.value);
          throw error;
        case 500:
          final error = GtApiException(
            cause: 'Internal server error: ${response.body}',
            type: GtApiExceptionType.serverError,
          );
          log(error.cause, error: error, level: Level.SEVERE.value);
          throw error;
        default:
          final error = GtApiException(
            cause: 'Unknown error: ${response.body}',
            type: GtApiExceptionType.unknownResponse,
          );
          log(error.cause, error: error, level: Level.SEVERE.value);
          throw Exception('Failed to refresh tokens: ${response.body}');
      }
    } on http.ClientException catch (error) {
      final gtError = GtApiException(
        cause: 'Could not connect to $baseUrl',
        type: GtApiExceptionType.hostNotResponding,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    } catch (error) {
      final gtError = GtApiException(
        cause: 'Unknown error happened during JWT refresh',
        type: GtApiExceptionType.unknown,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    }
  }

  /// Login with credentials
  ///
  /// Uses [username] and [password] to get tokens form the API.
  Future<void> login(String username, String password) async {
    if (!_hasBaseUrl()) {
      final error = Exception('BaseUrl not set');
      log('App has no baseUrl', level: Level.SEVERE.value, error: error);
      throw error;
    }

    try {
      var response = await http.post(
        Uri.parse('$baseUrl/auth/login'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      );

      if (response.statusCode == 200) {
        var data = jsonDecode(response.body);

        accessJWT = data['access_token'];
        refreshJWT = data['refresh_token'];

        final prefs = SharedPreferencesAsync();
        await prefs.setString('accessJWT', accessJWT!);
        await prefs.setString('refreshJWT', refreshJWT!);
        return;
      }

      late final String cause;
      late final GtApiExceptionType type;
      switch (response.statusCode) {
        case 400:
          cause = 'Malformed request: ${response.body}';
          type = GtApiExceptionType.malformedBody;
          break;
        case 401:
          cause = 'Invalid credentials: ${response.body}';
          type = GtApiExceptionType.invalidCredentials;
          break;
        case 500:
          cause = 'Server error: ${response.body}';
          type = GtApiExceptionType.serverError;
          break;
        default:
          cause = 'Unknown error: ${response.body}';
          type = GtApiExceptionType.unknownResponse;
      }
      final error = GtApiException(
        cause: cause,
        type: type,
      );
      log(error.cause, error: error, level: Level.SEVERE.value);
      throw error;
    } on GtApiException catch (_) {
      rethrow;
    } on http.ClientException catch (error) {
      final gtError = GtApiException(
        cause: 'Could not connect to $baseUrl',
        type: GtApiExceptionType.hostNotResponding,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    } catch (error) {
      final gtError = GtApiException(
        cause: 'Unknown error happened during login',
        type: GtApiExceptionType.unknown,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    }
  }

  Future<void> logout() async {
    if (!_hasBaseUrl()) {
      final error = Exception('BaseUrl not set');
      log('App has no baseUrl', level: Level.SEVERE.value, error: error);
      throw error;
    }

    try {
      var response = await http.post(Uri.parse('$baseUrl/auth/logout'),
          headers: {
            'Content-Type': 'application/json',
            'Authorization': 'Bearer $accessJWT'
          },
          body: jsonEncode({
            'refresh_token': refreshJWT,
          }));

      if (response.statusCode == 204) {
        var prefs = SharedPreferencesAsync();
        await prefs.remove('accessJWT');
        await prefs.remove('refreshJWT');
        return;
      }

      late final String cause;
      late final GtApiExceptionType type;
      switch (response.statusCode) {
        case 400:
          cause = 'Malformed request: ${response.body}';
          type = GtApiExceptionType.malformedBody;
          break;
        case 401:
          cause = 'Invalid credentials: ${response.body}';
          type = GtApiExceptionType.invalidCredentials;
          break;
        case 403:
          cause = 'Forbidden: ${response.body}';
          type = GtApiExceptionType.forbidden;
          break;
        case 500:
          cause = 'Server error: ${response.body}';
          type = GtApiExceptionType.serverError;
          break;
        default:
          cause = 'Unknown error: ${response.body}';
          type = GtApiExceptionType.unknownResponse;
          break;
      }
      final error = GtApiException(
        cause: cause,
        type: type,
      );
      log(error.cause, error: error, level: Level.SEVERE.value);
      throw error;
    } on GtApiException catch (_) {
      rethrow;
    } on SocketException catch (error) {
      final gtError = GtApiException(
        cause: 'Could not connect to $baseUrl',
        type: GtApiExceptionType.hostNotResponding,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    } catch (error) {
      final gtError = GtApiException(
        cause: 'Unknown error happened during logout',
        type: GtApiExceptionType.unknown,
      );
      log(gtError.cause, error: error, level: Level.SEVERE.value);
      throw gtError;
    }
  }
}

enum GtApiExceptionType {
  unknown,
  forbidden,
  hostUnknown,
  hostNotResponding,
  invalidCredentials,
  urlNull,
  refreshJWTNull,
  refreshJWTUnauthorized,
  malformedBody,
  serverError,
  unknownResponse,
}

class GtApiException implements Exception {
  String cause;
  GtApiExceptionType type;
  GtApiException({required this.cause, required this.type});

  @override
  String toString() {
    return 'GtApiException ($type): $cause';
  }
}

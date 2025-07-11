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
      }
    }
  }

  /// Check that baseUrl is not null
  bool _hasBaseUrl() {
    if (baseUrl == null) {
      return false;
    }
    return true;
  }

  /// Remove jwts from shared preferences and from memory.
  Future<void> _removeTokens() async {
    final prefs = SharedPreferencesAsync();
    await prefs.remove('accessJWT');
    await prefs.remove('refreshJWT');
    accessJWT = null;
    refreshJWT = null;
  }

  /// Throws proper GtApiException for error responses.
  ///
  /// To modify error cause messages, you can provide a [map] with status code
  /// as a key and message as a value. If not provided, default messages are used.
  _handleErrorStatus(http.Response response, {Map<int, String>? map}) {
    logError(GtApiException error, {Level level = Level.SEVERE}) => log(
          error.cause,
          error: error,
          level: level.value,
          stackTrace: StackTrace.current,
        );
    switch (response.statusCode) {
      case 400:
        final error = GtApiException(
          cause: map?[400] ?? 'Malformed body: ${response.body}',
          type: GtApiExceptionType.malformedBody,
        );
        logError(error);
        throw error;
      case 401:
        final error = GtApiException(
          cause: map?[401] ?? 'Unauthorized: ${response.body}',
          type: GtApiExceptionType.unauthorized,
        );
        logError(error, level: Level.INFO);
        throw error;
      case 403:
        final error = GtApiException(
          cause: map?[403] ?? 'Forbidden: ${response.body}',
          type: GtApiExceptionType.forbidden,
        );
        logError(error);
        throw error;
      case 409:
        final error = GtApiException(
          cause: map?[409] ?? 'Conflict: ${response.body}',
          type: GtApiExceptionType.conflict,
        );
        logError(error);
        throw error;
      case 500:
        final error = GtApiException(
          cause: map?[500] ?? 'Internal server error: ${response.body}',
          type: GtApiExceptionType.serverError,
        );
        logError(error);
        throw error;
      default:
        final error = GtApiException(
          cause: 'No handler defined for status ${response.statusCode}',
          type: GtApiExceptionType.unknownResponse,
        );
        logError(error);
        throw error;
    }
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
      if (response.statusCode != 200) {
        // If status is unauthorized, we need to remove the saved tokens
        if (response.statusCode == 401) {
          await _removeTokens();
        }

        _handleErrorStatus(
          response,
          map: {401: 'Token was unauthorized for refresh.'},
        );
      }
      var data = jsonDecode(response.body);
      accessJWT = data['access_token'];
      refreshJWT = data['refresh_token'];
      final prefs = SharedPreferencesAsync();
      await prefs.setString('accessJWT', accessJWT!);
      await prefs.setString('refreshJWT', refreshJWT!);
      log('Tokens refreshed', level: Level.INFO.value);
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

      if (response.statusCode != 200) {
        _handleErrorStatus(
          response,
          map: {401: 'Invalid credentials: ${response.body}'},
        );
      }
      var data = jsonDecode(response.body);

      accessJWT = data['access_token'];
      refreshJWT = data['refresh_token'];

      final prefs = SharedPreferencesAsync();
      await prefs.setString('accessJWT', accessJWT!);
      await prefs.setString('refreshJWT', refreshJWT!);
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

      if (response.statusCode != 204) {
        if (response.statusCode == 401) {
          await _removeTokens();
        }
        _handleErrorStatus(
          response,
          map: {401: 'Invalid credentials: ${response.body}'},
        );
      }

      await _removeTokens();
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

  /// Creates a user.
  ///
  /// Tries to create a new user with [username] and [password].
  Future<void> createUser(String username, String password) async {
    if (!_hasBaseUrl()) {
      final error = Exception('BaseUrl not set');
      log('App has no baseUrl', level: Level.SEVERE.value, error: error);
      throw error;
    }

    try {
      var response = await http.post(
        Uri.parse('$baseUrl/user/'),
        headers: {'Content-Type': 'application/json'},
        body: jsonEncode({
          'username': username,
          'password': password,
        }),
      );

      if (response.statusCode != 201) {
        _handleErrorStatus(response, map: {
          409: 'User with username $username already exists.',
        });
      }
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
        cause: 'Unknown error happened during user creation.',
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
  conflict,
  hostUnknown,
  hostNotResponding,
  urlNull,
  refreshJWTNull,
  unauthorized,
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

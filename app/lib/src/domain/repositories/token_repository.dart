import 'dart:async';
import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/env/env.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';

class TokenRepository {
  final SharedPreferences sharedPreferences;
  final Dio httpClient;

  TokenRepository({
    required this.sharedPreferences,
    required this.httpClient,
  });

  Future<Token> getToken(String email, String password) async {
    final response = await httpClient.post(
      '${Env.apiUrl}/sign/in',
      data: {
        'email': email,
        'password': password,
      },
    );
    if (response.statusCode == 200 && response.data != null) {
      final token = Token.fromJson(json.decode(response.data!) as Map<String, Object?>);
      return token;
    } else {
      throw Exception(
        'getToken ${response.statusCode}, data=${response.data}',
      );
    }
  }

  Future<Token> renewToken(Token token) async {
    final response = await httpClient.post(
      '${Env.apiUrl}/sign/renew',
      data: {
        'token': token,
      },
    );
    return Token.fromJson(response.data);
  }

  Token getLocalToken() {
    Map<String, dynamic> json = jsonDecode(sharedPreferences.getString('token') ?? '{}');
    return Token.fromJson(json);
  }

  void saveLocalToken(Token value) {
    sharedPreferences.setString("token", value.toString());
  }
}

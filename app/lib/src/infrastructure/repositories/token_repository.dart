import 'dart:convert';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository_interface.dart';

class TokenRepository implements TokenRepositoryInterface {
  final SharedPreferences _sharedPreferences;

  TokenRepository({
    required SharedPreferences sharedPreferences,
  }) : _sharedPreferences = sharedPreferences;

  @override
  Token getToken() {
    Map<String, dynamic> json = jsonDecode(_sharedPreferences.getString('token') ?? '{}');
    return Token.fromJson(json);
  }

  @override
  void saveToken(Token value) {
    _sharedPreferences.setString("token", value.toString());
  }
}

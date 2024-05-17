import 'dart:async';
import 'dart:convert';
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/env/env.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';
import 'package:thetiptop_client/src/infrastructure/services/dio/dio_service_provider.dart';
import 'package:thetiptop_client/src/infrastructure/services/exception/viewer_exception_provider.dart';

part 'generated/token_repository.g.dart';

class TokenRepository {
  final Dio dioService;
  final ViewerExceptionService viewerExceptionService;
  final SharedPreferences sharedPreferences;

  TokenRepository({
    required this.dioService,
    required this.viewerExceptionService,
    required this.sharedPreferences,
  });

  Future<Token?> getToken(String email, String password) async {
    try {
      FormData formData = FormData();
      formData.fields.add(MapEntry("email", email));
      formData.fields.add(MapEntry("password", password));
      final response = await dioService.post('${Env.apiUrl}/sign/in', data: formData);
      final token = Token.fromJson(json.decode(response.toString()));
      return token;
    } on DioException catch (error) {
      viewerExceptionService.showDioError(error);
      throw Exception(error);
    }
  }

  Future<Token> renewToken(Token token) async {
    final response = await dioService.post(
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

  void removeLocalToken() {
    sharedPreferences.remove("token");
  }
}

@riverpod
TokenRepository tokenRepository(TokenRepositoryRef ref) {
  return TokenRepository(
    dioService: ref.read(dioServiceProvider),
    viewerExceptionService: ref.read(viewerExceptionServiceProvider.notifier),
    sharedPreferences: ref.watch(sharedPreferencesProvider),
  );
}

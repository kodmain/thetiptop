import 'package:dio/dio.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository.dart';
import 'package:thetiptop_client/src/infrastructure/services/dio/dio_service_provider.dart';

class AuthInterceptor extends Interceptor {
  final DioServiceRef dioServiceRef;

  AuthInterceptor({required this.dioServiceRef});

  @override
  void onRequest(RequestOptions options, RequestInterceptorHandler handler) {
    final repo = dioServiceRef.read(tokenRepositoryProvider);
    final Token token = repo.getLocalToken();
    if (token.jwt.isNotEmpty) {
      options.headers['Authorization'] = 'Bearer ${token.jwt}';
    }
    return handler.next(options);
  }

  @override
  void onResponse(Response response, ResponseInterceptorHandler handler) {
    return handler.next(response);
  }

  @override
  void onError(DioException err, ErrorInterceptorHandler handler) {
    return handler.next(err);
  }
}

import 'package:dio/dio.dart';
import 'package:flutter/foundation.dart';
import 'package:pretty_dio_logger/pretty_dio_logger.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:thetiptop_client/src/infrastructure/services/dio/auth_interceptor.dart';

part 'generated/dio_service_provider.g.dart';

/// Provider est idéal pour accéder aux dépendances qui ne changent pas,
/// DioProvider pour au client Dio par le reste de l'application [dio_service_provider.g.dart]
///
@riverpod
Dio dioService(DioServiceRef ref) {
  final Dio dio = Dio(
    BaseOptions(
      connectTimeout: const Duration(seconds: 5),
      receiveTimeout: const Duration(seconds: 5),
      headers: {
        Headers.acceptHeader: 'application/json',
        Headers.contentTypeHeader: 'multipart/form-data',
      },
      followRedirects: false,
    ),
  );

  /// Arrête le client dio, juste avant la destruction du fournisseur.
  ///
  ref.onDispose(dio.close);

  /// Ajoute l'intercepteur d'authentication à l'instance.
  ///
  dio.interceptors.add(AuthInterceptor(dioServiceRef: ref));

  /// Ajoute l'intercepteur de log à l'instance.
  /// Uniquement en mode debug et profil.
  ///
  if (!kReleaseMode) {
    dio.interceptors.add(
      PrettyDioLogger(
        request: true,
        requestHeader: true,
        requestBody: true,
        responseBody: true,
        responseHeader: true,
        error: true,
      ),
    );
  }

  return dio;
}

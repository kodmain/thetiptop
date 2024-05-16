import 'dart:async';
import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:thetiptop_client/src/domain/env/env.dart';
import 'package:thetiptop_client/src/domain/models/client.dart';
import 'package:thetiptop_client/src/infrastructure/services/dio/dio_service_provider.dart';
import 'package:thetiptop_client/src/infrastructure/services/exception/viewer_exception_provider.dart';

part 'generated/client_repository.g.dart';

class ClientRepository {
  final Dio dioService;
  final ViewerExceptionService viewerExceptionService;

  ClientRepository({
    required this.dioService,
    required this.viewerExceptionService,
  });

  Future<Client?> signUp(String email, String password) async {
    try {
      FormData formData = FormData();
      formData.fields.add(MapEntry("email", email));
      formData.fields.add(MapEntry("password", password));
      await dioService.post('${Env.apiUrl}/sign/up', data: formData);
      return null;
    } on DioException catch (error) {
      viewerExceptionService.showDioError(error);
      throw Exception(error);
    }
  }
}

@riverpod
ClientRepository clientRepository(ClientRepositoryRef ref) {
  return ClientRepository(
    dioService: ref.read(dioServiceProvider),
    viewerExceptionService: ref.read(viewerExceptionServiceProvider.notifier),
  );
}

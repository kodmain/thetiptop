import 'dart:async';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:thetiptop_client/src/domain/models/client.dart';
import 'package:thetiptop_client/src/domain/repositories/client_repository.dart';

part 'generated/client_controller.g.dart';

@riverpod
class ClientController extends _$ClientController {
  @override
  FutureOr<Client?> build() async {
    state = const AsyncData(null);
    return state.value;
  }

  /// Création d'un compte client
  ///
  Future<void> signUp(Map<String, dynamic> formData) async {
    final repo = ref.read(clientRepositoryProvider);
    state = const AsyncLoading<Client?>();
    state = await AsyncValue.guard(() async {
      return repo.signUp(formData['email'], formData['password']);
    });
  }
}

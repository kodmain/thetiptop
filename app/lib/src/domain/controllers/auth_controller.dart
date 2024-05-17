import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/enums/auth_status.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository.dart';

part 'generated/auth_controller.g.dart';

@riverpod
class AuthController extends _$AuthController {
  @override
  FutureOr<AuthStatus?> build() async {
    state = const AsyncData(AuthStatus.disconnected);
    return state.value;
  }

  FutureOr<AuthStatus?> signin(Map<String, dynamic> formData) async {
    final repo = ref.read(tokenRepositoryProvider);
    state = const AsyncLoading<AuthStatus?>();
    state = await AsyncValue.guard(() async {
      Token? token = await repo.getToken(formData['email'], formData['password']);
      if (token == null) {
        repo.removeLocalToken();
        return AuthStatus.disconnected;
      }
      repo.saveLocalToken(token);
      return AuthStatus.connected;
    });

    return state.value;
  }

  FutureOr<AuthStatus?> signout() async {
    final repo = ref.read(tokenRepositoryProvider);
    state = const AsyncLoading<AuthStatus?>();
    state = await AsyncValue.guard(() async {
      repo.removeLocalToken();
      return AuthStatus.disconnected;
    });
    return state.value;
  }
}

/*import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/enums/auth_status.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository.dart';



class AuthController extends Notifier<AuthStatus> {
  @override
  AuthStatus build() {
    final SharedPreferences sharedPreferences = ref.watch(sharedPreferencesProvider);
    final tokenRepository = TokenRepository(sharedPreferences: sharedPreferences);
    return AuthStatus.disconnected;
  }

  void signin() {
    state = AuthStatus.connected;
  }

  void signout() {
    state = AuthStatus.disconnected;
  }
}
*/
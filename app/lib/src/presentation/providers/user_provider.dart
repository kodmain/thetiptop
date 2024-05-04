import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/controllers/auth_controller.dart';
import 'package:thetiptop_client/src/domain/controllers/user_controller.dart';
import 'package:thetiptop_client/src/domain/models/response.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';
import 'package:thetiptop_client/src/infrastructure/repositories/token_repository.dart';

/*final authProvider = StateNotifierProvider<AuthController>((ref) {
  return AuthController();
});*/

final userProvider = AsyncNotifierProvider<UserController, Response>(UserController.new);

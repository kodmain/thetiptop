import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/domain/controllers/auth_controller.dart';
import 'package:thetiptop_client/src/domain/enums/auth_status.dart';

final authProvider = NotifierProvider<AuthController, AuthStatus>(AuthController.new);

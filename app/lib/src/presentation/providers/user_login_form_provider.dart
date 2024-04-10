import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/domain/models/user_login.dart';
import 'package:thetiptop_client/src/domain/models/user_login_form_state.dart';

class UserLoginFormProvider extends StateNotifier<UserLoginFormState> {
  UserLoginFormProvider() : super(UserLoginFormState(UserLogin.empty()));
}

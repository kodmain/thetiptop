import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:thetiptop_client/src/domain/models/user_login.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/user_login_form_state.freezed.dart';
part 'generated/user_login_form_state.g.dart';

@freezed
class UserLoginFormState with _$UserLoginFormState {
  const factory UserLoginFormState(UserLogin form) = _UserLoginFormState;

  factory UserLoginFormState.fromJson(Map<String, dynamic> json) => _$UserLoginFormStateFromJson(json);
}

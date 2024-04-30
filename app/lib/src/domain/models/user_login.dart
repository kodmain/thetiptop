import 'package:freezed_annotation/freezed_annotation.dart';
import 'package:thetiptop_client/src/domain/models/input.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/user_login.freezed.dart';
part 'generated/user_login.g.dart';

@freezed
class UserLogin with _$UserLogin {
  const UserLogin._();
  const factory UserLogin({
    required Input email,
    required Input password,
  }) = _UserLogin;

  factory UserLogin.empty() => const UserLogin(
        email: Input(value: ''),
        password: Input(value: ''),
      );

  factory UserLogin.fromJson(Map<String, dynamic> json) => _$UserLoginFromJson(json);
}

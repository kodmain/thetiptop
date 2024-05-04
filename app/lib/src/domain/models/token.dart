import 'package:freezed_annotation/freezed_annotation.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/token.freezed.dart';
part 'generated/token.g.dart';

@freezed
class Token with _$Token {
  const factory Token({
    @Default("") String jwt,
  }) = _Token;

  factory Token.fromJson(Map<String, dynamic> json) => _$TokenFromJson(json);
}

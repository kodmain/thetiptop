import 'package:freezed_annotation/freezed_annotation.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/input.freezed.dart';
part 'generated/input.g.dart';

@freezed
class Input with _$Input {
  const factory Input({
    required String value,
    @Default("") String error,
    @Default(false) bool isValid,
  }) = _Input;

  factory Input.fromJson(Map<String, dynamic> json) => _$InputFromJson(json);
}

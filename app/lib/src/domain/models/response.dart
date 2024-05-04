import 'package:freezed_annotation/freezed_annotation.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/response.freezed.dart';
part 'generated/response.g.dart';

@freezed
class Response with _$Response {
  const factory Response({
    @Default("") String body,
  }) = _Response;

  factory Response.fromJson(Map<String, dynamic> json) => _$ResponseFromJson(json);
}

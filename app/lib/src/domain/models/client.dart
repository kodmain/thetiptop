import 'package:freezed_annotation/freezed_annotation.dart';

/// flutter pub run build_runner build
/// flutter pub run build_runner watch
///
part 'generated/client.freezed.dart';
part 'generated/client.g.dart';

@freezed
class Client with _$Client {
  const factory Client({
    required String id,
    required String email,
  }) = _Client;

  factory Client.fromJson(Map<String, dynamic> json) => _$ClientFromJson(json);
}

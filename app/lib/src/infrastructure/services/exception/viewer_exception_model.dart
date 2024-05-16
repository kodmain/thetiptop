import 'package:freezed_annotation/freezed_annotation.dart';

part 'generated/viewer_exception_model.freezed.dart';
part 'generated/viewer_exception_model.g.dart';

@freezed
class ViewerException with _$ViewerException {
  const factory ViewerException({
    required String error,
  }) = _ViewerException;

  factory ViewerException.fromJson(Map<String, dynamic> json) => _$ViewerExceptionFromJson(json);
}

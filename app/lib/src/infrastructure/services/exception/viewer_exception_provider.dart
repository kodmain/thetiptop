import 'dart:convert';

import 'package:dio/dio.dart';
import 'package:riverpod_annotation/riverpod_annotation.dart';
import 'package:thetiptop_client/src/infrastructure/services/exception/viewer_exception_model.dart';

part 'generated/viewer_exception_provider.g.dart';

@riverpod
class ViewerExceptionService extends _$ViewerExceptionService {
  @override
  ViewerException? build() {
    state = null;
    return state;
  }

  ViewerException? showDioError(DioException err) {
    state = ViewerException.fromJson(json.decode(err.response.toString()));
    return state;
  }
}

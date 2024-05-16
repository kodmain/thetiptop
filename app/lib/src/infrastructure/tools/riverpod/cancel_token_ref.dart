import 'dart:developer';
import 'package:dio/dio.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';

extension CancelTokenRef on Ref {
  CancelToken cancelTokenOnDispose() {
    log('cancelTokenOnDispose | add $hashCode');
    final cancelToken = CancelToken();
    onDispose(() {
      log('cancelTokenOnDispose | dispose $hashCode');
      cancelToken.cancel;
    });
    return cancelToken;
  }
}

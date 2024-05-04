import 'dart:async';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/domain/models/response.dart';

class UserController extends AsyncNotifier<Response> {
  @override
  Future<Response> build() async {
    return Response.fromJson({});
  }

  Future<void> create_account(Map<String, String?> formData) async {
    print("create_account");
    print(formData.toString());
  }
}

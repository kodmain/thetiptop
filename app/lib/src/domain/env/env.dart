import 'package:envied/envied.dart';

part 'generated/env.g.dart';

@Envied(path: '.env')
final class Env {
  @EnviedField(varName: 'API_URL', obfuscate: false)
  static const String apiUrl = _Env.apiUrl;
}

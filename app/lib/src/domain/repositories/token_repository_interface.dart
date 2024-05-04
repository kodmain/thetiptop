import '../models/token.dart';

abstract class TokenRepositoryInterface {
  Token getToken();
  void saveToken(Token value);
}

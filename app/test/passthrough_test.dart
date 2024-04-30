import 'package:flutter_test/flutter_test.dart';

void main() {
  group("passthrough", () {
    test("always true", () {
      expect(true, isTrue);
    });
  });
}

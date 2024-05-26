import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';

class MockVoidCallback extends Mock {
  void call();
}

void main() {
  group('BtnLinkWidget - only required', () {
    late MockVoidCallback onPressed;
    late Widget btnLinkWidget;

    setUp(() {
      onPressed = MockVoidCallback();
      btnLinkWidget = MaterialApp(
        home: Material(
          child: BtnLinkWidget(
            onPressed: onPressed,
            text: 'Submit',
          ),
        ),
      );
    });

    testWidgets('renders text', (WidgetTester tester) async {
      await tester.pumpWidget(btnLinkWidget);
      expect(find.textContaining('Submit', findRichText: true), findsOneWidget);
    });

    testWidgets('triggers onPressed callback', (WidgetTester tester) async {
      await tester.pumpWidget(btnLinkWidget);
      await tester.tap(find.byType(TextButton));
      verify(() => onPressed()).called(1);
    });

    testWidgets('renders auto style', (WidgetTester tester) async {
      await tester.pumpWidget(btnLinkWidget);
      final richText = tester.widget<RichText>(find.byType(RichText));
      final textSpan = richText.text as TextSpan;
      final textStyle = textSpan.style!;
      expect(textSpan.text, 'Submit');
      expect(textStyle.color, const Color.fromARGB(255, 51, 0, 255));
      expect(textStyle.decoration, TextDecoration.underline);
      expect(textStyle.decorationStyle, TextDecorationStyle.solid);
      expect(textStyle.decorationColor, const Color.fromARGB(255, 51, 0, 255));
    });
  });

  group('BtnLinkWidget - optionals', () {
    late MockVoidCallback onPressed;
    late Widget btnLinkWidget;

    setUp(() {
      onPressed = MockVoidCallback();
      btnLinkWidget = MaterialApp(
        home: Material(
          child: BtnLinkWidget(
            onPressed: onPressed,
            text: 'Submit',
            fontSize: 20.0,
            fontFamily: 'TestFont',
          ),
        ),
      );
    });

    testWidgets('renders style', (WidgetTester tester) async {
      await tester.pumpWidget(btnLinkWidget);
      final richText = tester.widget<RichText>(find.byType(RichText));
      final textSpan = richText.text as TextSpan;
      final textStyle = textSpan.style!;
      expect(textSpan.text, 'Submit');
      expect(textStyle.fontSize, 20.0);
      expect(textStyle.fontFamily, 'TestFont');
    });
  });
}

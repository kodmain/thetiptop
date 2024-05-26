import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';

class MockVoidCallback extends Mock {
  void call();
}

void main() {
  group('BtnActionWidget - only required', () {
    late MockVoidCallback onPressed;
    late Widget btnActionWidget;

    setUp(() {
      onPressed = MockVoidCallback();
      btnActionWidget = MaterialApp(
        home: Material(
          child: Row(
            children: [
              BtnActionWidget(
                onPressed: onPressed,
                text: 'Submit',
              ),
            ],
          ),
        ),
      );
    });

    testWidgets('renders text', (WidgetTester tester) async {
      await tester.pumpWidget(btnActionWidget);
      expect(find.text('Submit'), findsOneWidget);
    });

    testWidgets('triggers onPressed callback', (WidgetTester tester) async {
      await tester.pumpWidget(btnActionWidget);
      await tester.tap(find.text('Submit'));
      verify(() => onPressed()).called(1);
    });
  });

  group('BtnActionWidget - optionals', () {
    late MockVoidCallback onPressed;
    late Widget btnActionWidget;

    setUp(() {
      onPressed = MockVoidCallback();
      btnActionWidget = MaterialApp(
        home: Material(
          child: Row(
            children: [
              BtnActionWidget(
                onPressed: onPressed,
                text: 'Submit',
              ),
            ],
          ),
        ),
      );
    });

    testWidgets('disabledLoading default (false) - isLoading default (false)', (WidgetTester tester) async {
      await tester.pumpWidget(btnActionWidget);
      expect(find.byType(SizedBox), findsExactly(2));
      expect(find.text('Submit'), findsOneWidget);
    });

    testWidgets('disabledLoading (true) - isLoading default (false)', (WidgetTester tester) async {
      await tester.pumpWidget(MaterialApp(
        home: Material(
          child: Row(
            children: [
              BtnActionWidget(
                disableLoading: true,
                onPressed: onPressed,
                text: 'Submit',
              ),
            ],
          ),
        ),
      ));
      expect(find.byType(SizedBox), findsExactly(0));
      expect(find.text('Submit'), findsOneWidget);
    });

    testWidgets('disabledLoading (true) - isLoading (true)', (WidgetTester tester) async {
      await tester.pumpWidget(MaterialApp(
        home: Material(
          child: Row(
            children: [
              BtnActionWidget(
                disableLoading: true,
                isLoading: true,
                onPressed: onPressed,
                text: 'Submit',
              ),
            ],
          ),
        ),
      ));
      expect(find.byType(SizedBox), findsExactly(0));
      expect(find.text('Submit'), findsOneWidget);
    });

    testWidgets('disabledLoading default (false) - isLoading (true)', (WidgetTester tester) async {
      await tester.pumpWidget(MaterialApp(
        home: Material(
          child: Row(
            children: [
              BtnActionWidget(
                isLoading: true,
                onPressed: onPressed,
                text: 'Submit',
              ),
            ],
          ),
        ),
      ));
      expect(find.byType(SizedBox), findsExactly(3));
      expect(find.byType(CircularProgressIndicator), findsExactly(1));
      expect(find.text('Submit'), findsOneWidget);
    });
  });
}

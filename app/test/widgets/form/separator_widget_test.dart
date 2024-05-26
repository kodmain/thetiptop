import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/separator_widget.dart';

void main() {
  testWidgets('SeparatorWidget', (WidgetTester tester) async {
    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: SeparatorWidget(
            text: 'OR',
          ),
        ),
      ),
    );

    // Vérifie que le texte est affiché correctement
    expect(find.text('OR'), findsOneWidget);

    // Vérifie que les conteneurs de séparation sont présents
    final containers = find.byType(Container);
    expect(containers, findsNWidgets(2));

    // Vérifie que les conteneurs de séparation ont la bonne couleur
    expect((tester.widget<Container>(containers.at(0)).decoration as BoxDecoration).color, const Color.fromARGB(163, 51, 70, 70));
    expect((tester.widget<Container>(containers.at(1)).decoration as BoxDecoration).color, const Color.fromARGB(163, 51, 70, 70));

    // Vérifie que le style du texte est correct
    final textWidget = tester.widget<Text>(find.text('OR'));
    final textStyle = textWidget.style!;
    expect(textStyle.color, const Color.fromARGB(163, 51, 70, 70));
    expect(textStyle.fontFamily, "Raleway-SemiBold");
    expect(textStyle.fontSize, 18);

    // Vérifie que le padding du Row est correct
    final paddingWidget = tester.widget<Padding>(find.byType(Padding).first);
    final padding = paddingWidget.padding as EdgeInsets;
    expect(padding.top, 30);
    expect(padding.bottom, 30);
  });
}

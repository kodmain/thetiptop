import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:thetiptop_client/src/domain/enums/gain_status.dart';
import 'package:thetiptop_client/src/domain/enums/gain_type.dart';
import 'package:thetiptop_client/src/presentation/widgets/item/item_gain_widget.dart';

void main() {
  testWidgets('ItemGainWidget displays correct information', (WidgetTester tester) async {
    const gainType = GainType.signature;
    const gainStatus = GainStatus.wait;
    const infoNumberTicket = 'Ticket 123';
    const infoDateWin = '01/01/2024';

    await tester.pumpWidget(
      const MaterialApp(
        home: Scaffold(
          body: ItemGainWidget(
            type: gainType,
            status: gainStatus,
            infoNumberTicket: infoNumberTicket,
            infoDateWin: infoDateWin,
          ),
        ),
      ),
    );

    // Vérifie que l'icône SVG est affichée
    expect(find.byType(SvgPicture), findsOneWidget);

    // Vérifie que le label du type est affiché
    expect(find.text(gainType.getLabel()), findsOneWidget);

    // Vérifie que le numéro du ticket et la date sont affichés
    expect(find.textContaining(infoNumberTicket, findRichText: true), findsOneWidget);
    expect(find.textContaining(infoDateWin, findRichText: true), findsOneWidget);

    // Vérifie que le label du statut est affiché
    expect(find.text(gainStatus.getLabel()), findsOneWidget);

    // Vérifie que la couleur du conteneur de statut est correcte
    final statusContainer = tester.widget<Container>(
      find
          .descendant(
            of: find.byType(ItemGainWidget),
            matching: find.byWidgetPredicate((widget) => widget is Container && widget.decoration != null),
          )
          .last,
    );
    final boxDecoration = statusContainer.decoration as BoxDecoration;
    expect(boxDecoration.color, gainStatus.getColor());
  });
}

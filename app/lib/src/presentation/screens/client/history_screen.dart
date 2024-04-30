import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/domain/enums/gain_status.dart';
import 'package:thetiptop_client/src/domain/enums/gain_type.dart';
import 'package:thetiptop_client/src/presentation/widgets/item/item_gain_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_client_widget.dart';

class HistoryScreen extends HookConsumerWidget {
  const HistoryScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final List<String> entries = <String>['A', 'B', 'C'];

    return Semantics(
      header: true,
      label: "page historique des gains",
      child: LayoutClientWidget(
        child: Column(
          children: [
            const MenuClientWidget(),
            const SizedBox(
              height: 40,
            ),
            const Padding(
              padding: EdgeInsets.fromLTRB(0, 20, 0, 20),
              child: Text(
                "Historique des gains",
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: "Raleway",
                  fontSize: 22,
                ),
              ),
            ),
            ListView.separated(
              shrinkWrap: true,
              separatorBuilder: (BuildContext context, int index) => const SizedBox(height: 5),
              padding: const EdgeInsets.all(8),
              itemCount: entries.length,
              itemBuilder: (BuildContext context, int index) {
                return const ItemGainWidget(
                  type: GainType.bigCoffret,
                  status: GainStatus.wait,
                  infoNumberTicket: "1254 6458 10",
                  infoDateWin: "02/04/2024",
                );
              },
            ),
          ],
        ),
      ),
    );
  }
}

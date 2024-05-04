import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/btn_icon_big_widget.dart';

class MenuAdminWidget extends StatelessWidget {
  const MenuAdminWidget({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        const SizedBox(
          height: DefaultTheme.bigSpacer,
        ),
        BtnIconBigWidget(
          onPressed: () {},
          label: "Recherche utilisateur",
          asset: "assets/images/parts_icon-search-user.svg",
        ),
        BtnIconBigWidget(
          onPressed: () {},
          label: "Administrateur",
          asset: "assets/images/parts_icon-administrator.svg",
        ),
        BtnIconBigWidget(
          onPressed: () {},
          label: "Devices",
          asset: "assets/images/parts_icon-device.svg",
        ),
        BtnIconBigWidget(
          onPressed: () {},
          label: "Statistiques",
          asset: "assets/images/parts_icon-stats.svg",
        ),
        BtnIconBigWidget(
          onPressed: () {},
          label: "Paramètres",
          asset: "assets/images/parts_icon-param.svg",
        ),
      ],
    );
  }
}

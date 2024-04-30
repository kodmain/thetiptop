import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:go_router/go_router.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/btn_icon_big_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/btn_icon_widget.dart';

class MenuAdminWidget extends StatelessWidget {
  const MenuAdminWidget({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Column(
      mainAxisAlignment: MainAxisAlignment.center,
      crossAxisAlignment: CrossAxisAlignment.center,
      children: [
        SizedBox(
          height: 100,
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

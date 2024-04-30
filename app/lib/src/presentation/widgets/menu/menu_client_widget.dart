import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/btn_icon_widget.dart';

class MenuClientWidget extends StatelessWidget {
  const MenuClientWidget({
    super.key,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      mainAxisAlignment: MainAxisAlignment.center,
      children: [
        BtnIconWidget(
          onPressed: () {
            context.go(AppRouter.gameRouteName);
          },
          semanticLabel: "Accueil du jeu-concours ThéTipTop",
          asset: "assets/images/parts_icon-home.svg",
        ),
        BtnIconWidget(
          onPressed: () {
            context.go(AppRouter.profilRouteName);
          },
          semanticLabel: "Profil utilisateur",
          asset: "assets/images/parts_icon-user.svg",
        ),
        BtnIconWidget(
          onPressed: () {
            context.go(AppRouter.historyRouteName);
          },
          semanticLabel: "Historique des gains",
          asset: "assets/images/parts_icon-history.svg",
        ),
        BtnIconWidget(
          onPressed: () {
            context.go(AppRouter.signinRouteName);
          },
          semanticLabel: "Déconnexion",
          asset: "assets/images/parts_icon-exit.svg",
        ),
      ],
    );
  }
}

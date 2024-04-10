import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form_validator.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/checkbox_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/input_text_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/separator_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_client_widget.dart';

class ProfilScreen extends HookConsumerWidget {
  const ProfilScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = GlobalKey<FormState>();
    double screenWidth = MediaQuery.of(context).size.width;

    return Semantics(
      header: true,
      label: "page éditer votre profil",
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
                "Éditer votre profil",
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: "Raleway",
                  fontSize: 22,
                ),
              ),
            ),
            Form(
              key: formKey,
              child: Column(
                children: [
                  InputText(
                    labelText: "Adresse email",
                    initialValue: "steeve.dupond@gmail.com",
                    validator: (value) => FormValidator().isEmail(value: value),
                  ),
                  CheckboxWidget(
                    value: true,
                    onChanged: (value) {
                      print("ok$value");
                    },
                    text: "J’accepte de recevoir les newsletters de ThéTipTop.",
                  ),
                  InputText(
                    obscureText: true,
                    labelText: "Mot de passe",
                    validator: (value) => FormValidator().notEmpty(value: value),
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          context.go(AppRouter.signinRouteName);
                        },
                        backgroundColor: AppColor.greyCancel.color,
                        foregroundColor: AppColor.whiteCream.color,
                        text: "Annuler",
                      ),
                      const SizedBox(
                        width: 10,
                      ),
                      BtnActionWidget(
                        onPressed: () {
                          context.go(AppRouter.passwordResetRouteName);
                        },
                        backgroundColor: AppColor.red.color,
                        foregroundColor: AppColor.whiteCream.color,
                        text: "Enregistrer",
                      ),
                    ],
                  ),
                  const SeparatorWidget(
                    text: "Autres Actions",
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          context.go(AppRouter.passwordResetRouteName);
                        },
                        padding: const EdgeInsets.fromLTRB(0, 15, 0, 0),
                        backgroundColor: AppColor.blackGreen.color,
                        foregroundColor: AppColor.whiteCream.color,
                        text: "Éditer mon mot de passe",
                      ),
                    ],
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          print("get user data");
                        },
                        padding: const EdgeInsets.fromLTRB(0, 5, 0, 5),
                        backgroundColor: AppColor.blackGreen.color,
                        foregroundColor: AppColor.whiteCream.color,
                        text: "Récupérer mes données",
                      ),
                    ],
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          print("delete account");
                        },
                        padding: EdgeInsets.zero,
                        backgroundColor: AppColor.red.color,
                        foregroundColor: AppColor.whiteCream.color,
                        text: "Supprimer mon compte",
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

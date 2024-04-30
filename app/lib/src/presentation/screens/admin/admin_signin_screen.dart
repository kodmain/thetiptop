import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form_validator.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/input_text_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class AdminSigninScreen extends HookConsumerWidget {
  const AdminSigninScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = GlobalKey<FormState>();

    return LayoutClientWidget(
      child: Form(
        key: formKey,
        child: Column(
          children: [
            const SizedBox(
              height: 75,
            ),
            InputText(
              labelText: "Adresse email",
              validator: (value) => FormValidator().isEmail(value: value),
            ),
            InputText(
              obscureText: true,
              labelText: "Mot de passe",
              validator: (value) => FormValidator().notEmpty(value: value),
            ),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                BtnLinkWidget(
                  onPressed: () {
                    context.go(AppRouter.passwordRenewRouteName);
                  },
                  text: "Mot de passe oublié",
                ),
              ],
            ),
            Row(
              children: [
                BtnActionWidget(
                  onPressed: () {
                    if (formKey.currentState != null && formKey.currentState!.validate()) {
                      context.go(AppRouter.clientRouteName);
                    }
                  },
                  backgroundColor: AppColor.red.color,
                  foregroundColor: AppColor.whiteCream.color,
                  text: "Connexion",
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

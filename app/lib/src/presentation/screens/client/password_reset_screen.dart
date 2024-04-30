import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form_validator.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/input_text_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class PasswordResetScreen extends HookConsumerWidget {
  const PasswordResetScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = GlobalKey<FormState>();
    String? password;

    return LayoutClientWidget(
      child: Form(
        key: formKey,
        child: Column(
          children: [
            const Padding(
              padding: EdgeInsets.fromLTRB(0, 20, 0, 20),
              child: Text(
                "Un code de validation vient de vous \r être envoyé par email",
                textAlign: TextAlign.center,
                style: TextStyle(
                  fontFamily: "Raleway-Bold",
                  fontSize: 16,
                ),
              ),
            ),
            InputText(
              labelText: "Code de Validation",
              validator: (value) => FormValidator().notEmpty(value: value),
            ),
            InputText(
              obscureText: true,
              labelText: "Nouveau mot de passe",
              validator: (value) {
                password = value;
                return FormValidator().isComplex(value: value);
              },
            ),
            InputText(
              obscureText: true,
              labelText: "Confirmation du nouveau mot de passe",
              validator: (value) => FormValidator().isEqual(firstValue: value, secondValue: password, message: "Erreur de confirmation de mot de passe"),
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
                    context.go(AppRouter.signinRouteName);
                  },
                  backgroundColor: AppColor.red.color,
                  foregroundColor: AppColor.whiteCream.color,
                  text: "Enregistrer",
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form_validator.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class AdminSigninScreen extends ConsumerStatefulWidget {
  const AdminSigninScreen({super.key});

  @override
  AdminSigninScreenState createState() => AdminSigninScreenState();
}

class AdminSigninScreenState extends ConsumerState<AdminSigninScreen> {
  // Clé globale pour le widget Form
  final _formKey = GlobalKey<FormState>();

  // État pour le champ texte
  String _email = '';
  String _password = '';

  @override
  Widget build(BuildContext context) {
    final formKey = GlobalKey<FormState>();

    return LayoutClientWidget(
      child: Form(
        key: formKey,
        child: Column(
          children: [
            const SizedBox(
              height: DefaultTheme.bigSpacer,
            ),
            TextFormField(
              decoration: const InputDecoration(labelText: 'Adresse email'),
              validator: (value) => FormValidator().isEmail(value: value),
              onSaved: (value) {
                _email = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              decoration: const InputDecoration(labelText: 'Mot de passe'),
              validator: (value) => FormValidator().notEmpty(value: value),
              onSaved: (value) {
                _password = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.smallSpacer,
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
                  style: Theme.of(context).outlinedButtonTheme.style,
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

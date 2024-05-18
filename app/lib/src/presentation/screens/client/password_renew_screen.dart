import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class PasswordRenewScreen extends ConsumerStatefulWidget {
  const PasswordRenewScreen({super.key});

  @override
  PasswordRenewScreenState createState() => PasswordRenewScreenState();
}

class PasswordRenewScreenState extends ConsumerState<PasswordRenewScreen> {
  // Clé globale pour le widget Form
  final _formKey = GlobalKey<FormState>();

  // État pour le champ texte
  String _email = '';

  @override
  Widget build(BuildContext context) {
    final formKey = GlobalKey<FormState>();

    return LayoutClientWidget(
      child: Form(
        key: _formKey,
        child: Column(
          children: [
            const SizedBox(
              height: 75,
            ),
            TextFormField(
              decoration: InputDecoration(labelText: context.loc.labelEmail),
              validator: (value) => Validator().isEmail(value: value),
              onSaved: (value) {
                _email = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            Row(
              children: [
                BtnActionWidget(
                  onPressed: () {
                    context.go(AppRouter.signinRouteName);
                  },
                  style: Theme.of(context).filledButtonTheme.style,
                  text: context.loc.btnAnnuler,
                ),
                const SizedBox(
                  width: 10,
                ),
                BtnActionWidget(
                  onPressed: () {
                    context.go(AppRouter.passwordResetRouteName);
                  },
                  style: Theme.of(context).outlinedButtonTheme.style,
                  text: context.loc.btnSuivant,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

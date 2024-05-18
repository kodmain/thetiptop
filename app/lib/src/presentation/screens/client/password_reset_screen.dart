import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class PasswordResetScreen extends ConsumerStatefulWidget {
  const PasswordResetScreen({super.key});

  @override
  PasswordResetScreenState createState() => PasswordResetScreenState();
}

class PasswordResetScreenState extends ConsumerState<PasswordResetScreen> {
  @override
  Widget build(BuildContext context) {
    // Clé globale pour le widget Form
    final _formKey = GlobalKey<FormState>();

    // État pour le champ texte
    String _email = '';
    String _password = '';
    String _codeValidation = '';

    return LayoutClientWidget(
      child: Form(
        key: _formKey,
        child: Column(
          children: [
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 50, 0, 50),
              child: Text(
                context.loc.infoValidationCode,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.bodyLarge,
              ),
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(labelText: context.loc.labelValidationCode),
              validator: (value) => Validator().notEmpty(value: value),
              onSaved: (value) {
                _codeValidation = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(labelText: context.loc.labelNewPassword),
              validator: (value) {
                _password = value!;
                return Validator().isComplex(value: value);
              },
              onSaved: (value) {
                _password = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(labelText: context.loc.labelNewPasswordConfirm),
              validator: (value) => Validator().isEqual(
                firstValue: value,
                secondValue: _password,
                message: context.loc.errorPasswordConfirm,
              ),
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
                  width: DefaultTheme.smallSpacer,
                ),
                BtnActionWidget(
                  onPressed: () {
                    context.go(AppRouter.signinRouteName);
                  },
                  style: Theme.of(context).outlinedButtonTheme.style,
                  text: context.loc.btnEnregistrer,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

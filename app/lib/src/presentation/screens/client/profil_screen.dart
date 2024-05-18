import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/checkbox_form_field.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/separator_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_client_widget.dart';

class ProfilScreen extends ConsumerStatefulWidget {
  const ProfilScreen({super.key});
  @override
  ProfilScreenState createState() => ProfilScreenState();
}

class ProfilScreenState extends ConsumerState<ProfilScreen> {
  @override
  Widget build(BuildContext context) {
    // Clé globale pour le widget Form
    final _formKey = GlobalKey<FormState>();

    // État pour le champ texte
    String _email = '';
    String _password = '';
    bool _isAgreeNewsletter = true;

    return Semantics(
      header: true,
      label: "page ${context.loc.pageProfilTitle}",
      child: LayoutClientWidget(
        child: Column(
          children: [
            const MenuClientWidget(),
            const SizedBox(
              height: 40,
            ),
            Padding(
              padding: const EdgeInsets.fromLTRB(0, 20, 0, 20),
              child: Text(
                context.loc.pageProfilTitle,
                textAlign: TextAlign.center,
                style: Theme.of(context).textTheme.titleLarge,
              ),
            ),
            Form(
              key: _formKey,
              child: Column(
                children: [
                  TextFormField(
                    decoration: InputDecoration(labelText: context.loc.labelEmail),
                    initialValue: "steeve.dupond@gmail.com",
                    validator: (value) => Validator().isEmail(value: value),
                    onSaved: (value) {
                      _email = value!;
                    },
                  ),
                  const SizedBox(
                    height: DefaultTheme.smallSpacer,
                  ),
                  CheckboxFormField(
                    textLabel: context.loc.checkBoxAgreeNewsletter,
                    textStyle: Theme.of(context).textTheme.bodyMedium,
                    initialValue: _isAgreeNewsletter,
                    onChanged: (value) {
                      setState(() {
                        _isAgreeNewsletter = value;
                      });
                    },
                  ),
                  const SizedBox(
                    height: DefaultTheme.smallSpacer,
                  ),
                  TextFormField(
                    decoration: InputDecoration(labelText: context.loc.labelPassword),
                    validator: (value) => Validator().notEmpty(value: value),
                    onSaved: (value) {
                      _password = value!;
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
                        width: DefaultTheme.smallSpacer,
                      ),
                      BtnActionWidget(
                        onPressed: () {
                          context.go(AppRouter.passwordResetRouteName);
                        },
                        style: Theme.of(context).outlinedButtonTheme.style,
                        text: context.loc.btnEnregistrer,
                      ),
                    ],
                  ),
                  SeparatorWidget(
                    text: context.loc.separatorOtherAction,
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          context.go(AppRouter.passwordResetRouteName);
                        },
                        style: Theme.of(context).elevatedButtonTheme.style,
                        text: context.loc.btnChangePassword,
                      ),
                    ],
                  ),
                  const SizedBox(
                    height: DefaultTheme.smallSpacer,
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          print("get user data");
                        },
                        style: Theme.of(context).elevatedButtonTheme.style,
                        text: context.loc.btnGetMyData,
                      ),
                    ],
                  ),
                  const SizedBox(
                    height: DefaultTheme.smallSpacer,
                  ),
                  Row(
                    children: [
                      BtnActionWidget(
                        onPressed: () {
                          print("delete account");
                        },
                        style: Theme.of(context).outlinedButtonTheme.style,
                        text: context.loc.btnDeleteMyAccount,
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

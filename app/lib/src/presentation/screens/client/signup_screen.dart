import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/models/response.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form_validator.dart';
import 'package:thetiptop_client/src/presentation/providers/user_provider.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/checkbox_form_field.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class SignupScreen extends ConsumerStatefulWidget {
  const SignupScreen({super.key});

  @override
  SignupScreenState createState() => SignupScreenState();
}

class SignupScreenState extends ConsumerState<SignupScreen> {
  // Clé globale pour le widget Form
  final _formKey = GlobalKey<FormState>();

  // État pour le champ texte
  String _email = '';
  String _password = '';
  String _passwordConf = '';

  // État pour la case à cocher
  bool _isAgreeNewsletter = false;
  bool _isAgreeTerms = false;

  @override
  Widget build(BuildContext context) {
    return LayoutClientWidget(
      child: Form(
        key: _formKey,
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
              obscureText: true,
              decoration: const InputDecoration(labelText: "Mot de passe"),
              validator: (value) {
                _password = value!;
                return FormValidator().isComplex(value: value);
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
              decoration: const InputDecoration(labelText: "Confirmation du mot de passe"),
              validator: (value) => FormValidator().isEqual(
                firstValue: value,
                secondValue: _password,
                message: "Erreur de confirmation de mot de passe",
              ),
            ),
            const SizedBox(
              height: 15,
            ),
            CheckboxFormField(
              textLabel: "J’accepte les conditions générales d’utilisation. ",
              linkLabel: "CGU",
              linkUrl: "https://concours.thetiptop.fr/cgu",
              textStyle: Theme.of(context).textTheme.bodyMedium,
              validator: (value) => FormValidator().isTrue(value: value),
              initialValue: false,
              onChanged: (value) {
                setState(() {
                  _isAgreeTerms = value;
                });
              },
            ),
            CheckboxFormField(
              textLabel: "J’accepte de recevoir les newsletters de ThéTipTop.",
              textStyle: Theme.of(context).textTheme.bodyMedium,
              initialValue: false,
              onChanged: (value) {
                setState(() {
                  _isAgreeNewsletter = value;
                });
              },
            ),
            const SizedBox(
              height: DefaultTheme.smallSpacer,
            ),
            Row(
              children: [
                BtnActionWidget(
                  onPressed: () {
                    context.go(AppRouter.signinRouteName);
                  },
                  style: Theme.of(context).filledButtonTheme.style,
                  text: "Annuler",
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                BtnActionWidget(
                  onPressed: () {
                    if (_formKey.currentState != null && _formKey.currentState!.validate()) {
                      _formKey.currentState!.save();
                    }
                  },
                  style: Theme.of(context).outlinedButtonTheme.style,
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

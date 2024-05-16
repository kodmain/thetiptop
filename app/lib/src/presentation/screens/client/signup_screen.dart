import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/controllers/client_controller.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';
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

  // État du formulaire
  Map<String, dynamic> formData = {
    'email': '',
    'password': '',
    'isAgreeTerms': false,
    'isAgreeNewsletter': false,
  };

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(clientControllerProvider);

    return LayoutClientWidget(
      child: Form(
        key: _formKey,
        child: Column(
          children: [
            const SizedBox(
              height: DefaultTheme.bigSpacer,
            ),
            TextFormField(
              decoration: InputDecoration(
                labelText: 'Adresse email'.hardcoded,
              ),
              validator: (value) => Validator().isEmail(value: value),
              onSaved: (value) {
                formData['email'] = value;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(
                labelText: "Mot de passe".hardcoded,
              ),
              validator: (value) {
                formData['password'] = value;
                return Validator().isComplex(value: value);
              },
              onSaved: (value) {
                formData['password'] = value;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(
                labelText: "Confirmation du mot de passe".hardcoded,
              ),
              validator: (value) => Validator().isEqual(
                firstValue: value,
                secondValue: formData['password'],
                message: "Erreur de confirmation de mot de passe".hardcoded,
              ),
            ),
            const SizedBox(
              height: 15,
            ),
            CheckboxFormField(
              textLabel: "J’accepte les conditions générales d’utilisation. ".hardcoded,
              linkLabel: "CGU",
              linkUrl: "https://concours.thetiptop.fr/cgu",
              textStyle: Theme.of(context).textTheme.bodyMedium,
              validator: (value) => Validator().isTrue(value: value),
              initialValue: false,
              onChanged: (value) {
                setState(() {
                  formData['isAgreeTerms'] = value;
                });
              },
              onSaved: (value) {
                formData['isAgreeTerms'] = value;
              },
            ),
            CheckboxFormField(
              textLabel: "J’accepte de recevoir les newsletters de ThéTipTop.".hardcoded,
              textStyle: Theme.of(context).textTheme.bodyMedium,
              initialValue: false,
              onChanged: (value) {
                setState(() {
                  formData['isAgreeNewsletter'] = value;
                });
              },
              onSaved: (value) {
                formData['isAgreeNewsletter'] = value;
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
                  text: "Annuler".hardcoded,
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                BtnActionWidget(
                  isLoading: state.isLoading,
                  onPressed: state.isLoading
                      ? null
                      : () {
                          if (_formKey.currentState!.validate()) {
                            _formKey.currentState!.save();
                            ref.read(clientControllerProvider.notifier).signUp(formData);
                          }
                        },
                  style: Theme.of(context).outlinedButtonTheme.style,
                  text: "Enregistrer".hardcoded,
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

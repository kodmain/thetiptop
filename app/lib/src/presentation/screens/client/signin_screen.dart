import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';
import 'package:thetiptop_client/src/domain/controllers/auth_controller.dart';
import 'package:thetiptop_client/src/domain/controllers/client_controller.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_action_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/form/separator_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';

class SigninScreen extends ConsumerStatefulWidget {
  const SigninScreen({super.key});

  @override
  SigninScreenState createState() => SigninScreenState();
}

class SigninScreenState extends ConsumerState<SigninScreen> {
  // Clé globale pour le widget Form
  final _formKey = GlobalKey<FormState>();

  // État du formulaire
  Map<String, dynamic> formData = {
    'email': '',
    'password': '',
  };

  @override
  Widget build(BuildContext context) {
    final state = ref.watch(clientControllerProvider);
    double screenWidth = MediaQuery.of(context).size.width;

    return LayoutClientWidget(
      child: Form(
        key: _formKey,
        child: Column(
          children: [
            const SizedBox(
              height: DefaultTheme.bigSpacer,
            ),
            TextFormField(
              decoration: InputDecoration(labelText: context.loc.labelEmail),
              validator: (value) => Validator().isEmail(value: value),
              onSaved: (value) {
                formData['email'] = value!;
              },
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            TextFormField(
              obscureText: true,
              decoration: InputDecoration(labelText: context.loc.labelPassword),
              validator: (value) => Validator().notEmpty(value: value),
              onSaved: (value) {
                formData['password'] = value!;
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
                  text: context.loc.linkLostPassword,
                ),
              ],
            ),
            const SizedBox(
              height: DefaultTheme.smallSpacer,
            ),
            Row(
              children: [
                BtnActionWidget(
                  isLoading: state.isLoading,
                  onPressed: state.isLoading
                      ? null
                      : () {
                          if (_formKey.currentState!.validate()) {
                            _formKey.currentState!.save();
                            ref.read(authControllerProvider.notifier).signin(formData);
                          }
                        },
                  style: Theme.of(context).outlinedButtonTheme.style,
                  text: context.loc.btnConnexion,
                ),
              ],
            ),
            SeparatorWidget(
              text: context.loc.separatorOr,
            ),
            Row(
              children: [
                BtnActionWidget(
                  disableLoading: true,
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).elevatedButtonTheme.style?.copyWith(
                        backgroundColor: const MaterialStatePropertyAll(
                          DefaultTheme.blueFC,
                        ),
                      ),
                  text: screenWidth > 950 ? context.loc.btnFacebookConnect('') : context.loc.btnFacebookConnect('rl'),
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                BtnActionWidget(
                  disableLoading: true,
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).elevatedButtonTheme.style?.copyWith(
                        backgroundColor: const MaterialStatePropertyAll(
                          DefaultTheme.blueGC,
                        ),
                      ),
                  text: screenWidth > 950 ? context.loc.btnGoogleConnect('') : context.loc.btnGoogleConnect('rl'),
                ),
              ],
            ),
            const SizedBox(
              height: DefaultTheme.spacer,
            ),
            MergeSemantics(
              child: Row(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text(
                    context.loc.textInscription,
                  ),
                  BtnLinkWidget(
                    onPressed: () {
                      context.go(AppRouter.signupRouteName);
                    },
                    text: context.loc.linkInscription,
                    fontFamily: DefaultTheme.fontFamilyBold,
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

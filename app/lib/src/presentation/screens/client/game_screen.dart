import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter/scheduler.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:rive/rive.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/layout_client_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_client_widget.dart';
import 'dart:html' as html;

class GameScreen extends HookConsumerWidget {
  const GameScreen({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    final formKey = GlobalKey<FormState>();
    SMIInput<double>? lotInput;

    return Semantics(
      header: true,
      label: "page jeu",
      child: LayoutClientWidget(
        child: Column(
          children: [
            const MenuClientWidget(),
            const SizedBox(
              height: 40,
            ),
            Semantics(
              label: "Roue de la chance",
              child: Form(
                key: formKey,
                autovalidateMode: AutovalidateMode.always,
                child: ConstrainedBox(
                  constraints: BoxConstraints.loose(
                    const Size(475, 350),
                  ),
                  child: Stack(
                    children: [
                      RiveAnimation.asset(
                        "assets/rive/roue.riv",
                        fit: BoxFit.contain,
                        onInit: (Artboard artboard) {
                          /// init controller
                          ///
                          final controller = StateMachineController.fromArtboard(artboard, 'rotation');
                          artboard.addController(controller!);

                          /// init input
                          ///
                          lotInput = controller.findInput<double>('lot') as SMINumber;

                          /// Ajoute un écouteur d'événement Rive
                          ///
                          controller.addEventListener((RiveEvent event) {
                            SchedulerBinding.instance.addPostFrameCallback((_) {
                              if (formKey.currentState != null && event.name == "click" && formKey.currentState!.validate()) {
                                lotInput?.change(0);
                              }

                              if (event.name == "launch") {
                                Future.delayed(const Duration(milliseconds: 3000), () {
                                  lotInput?.change(1);
                                });
                              }
                            });
                          });
                        },
                      ),
                      Container(
                        height: 140,
                        decoration: BoxDecoration(
                          color: DefaultTheme.blackGreen,
                          borderRadius: BorderRadius.circular(15),
                        ),
                        child: Padding(
                          padding: const EdgeInsets.fromLTRB(15, 25, 15, 15),
                          child: TextFormField(
                            onFieldSubmitted: (value) {
                              lotInput?.change(0);
                            },
                            keyboardType: TextInputType.number,
                            expands: false,
                            maxLines: 1,
                            minLines: 1,
                            maxLength: 10,
                            textAlign: TextAlign.center,
                            textAlignVertical: TextAlignVertical.center,
                            cursorColor: DefaultTheme.white,
                            style: TextStyle(
                              color: DefaultTheme.white,
                              fontWeight: FontWeight.bold,
                              fontSize: 25,
                              textBaseline: TextBaseline.ideographic,
                            ),
                            decoration: InputDecoration(
                                contentPadding: const EdgeInsets.fromLTRB(0, 25, 0, 25),
                                floatingLabelAlignment: FloatingLabelAlignment.center,
                                floatingLabelBehavior: FloatingLabelBehavior.always,
                                alignLabelWithHint: true,
                                filled: true,
                                fillColor: const Color.fromARGB(20, 217, 217, 217),
                                focusColor: DefaultTheme.inputFill,
                                focusedBorder: OutlineInputBorder(
                                  gapPadding: 25,
                                  borderRadius: BorderRadius.circular(10.0),
                                  borderSide: BorderSide(
                                    color: DefaultTheme.whiteCream,
                                    width: 2,
                                  ),
                                ),
                                enabledBorder: OutlineInputBorder(
                                  gapPadding: 25,
                                  borderRadius: BorderRadius.circular(10.0),
                                  borderSide: BorderSide(
                                    color: DefaultTheme.whiteCream,
                                    width: 2,
                                  ),
                                ),
                                errorBorder: OutlineInputBorder(
                                  gapPadding: 25,
                                  borderRadius: BorderRadius.circular(10.0),
                                  borderSide: BorderSide(
                                    color: DefaultTheme.whiteCream,
                                    width: 2,
                                  ),
                                ),
                                focusedErrorBorder: OutlineInputBorder(
                                  gapPadding: 25,
                                  borderRadius: BorderRadius.circular(10.0),
                                  borderSide: BorderSide(
                                    color: DefaultTheme.whiteCream,
                                    width: 2,
                                  ),
                                ),
                                hintText: "",
                                labelText: "Entrez votre numéro de ticket",
                                labelStyle: TextStyle(
                                  color: DefaultTheme.white,
                                  fontWeight: FontWeight.bold,
                                  fontFamily: DefaultTheme.fontFamily,
                                  fontSize: 20,
                                ),
                                floatingLabelStyle: TextStyle(
                                  color: DefaultTheme.white,
                                  fontWeight: FontWeight.bold,
                                  fontFamily: DefaultTheme.fontFamily,
                                  fontSize: 20,
                                ),
                                errorStyle: TextStyle(
                                  color: DefaultTheme.white,
                                  fontWeight: FontWeight.bold,
                                  fontFamily: DefaultTheme.fontFamily,
                                  fontSize: 13,
                                )),
                            validator: (value) {
                              if (!RegExp(r"^[0-9]{10}$").hasMatch(value ?? "")) {
                                return "Veuillez entrer un numéro de ticket valide";
                              }
                              return null;
                            },
                          ),
                        ),
                      ),
                    ],
                  ),
                ),
              ),
            ),
            const SizedBox(
              height: 75,
            ),
            MergeSemantics(
              child: Column(
                children: [
                  RichText(
                    textAlign: TextAlign.center,
                    text: const TextSpan(
                      text: "Votre n’avez pas de numéro de ticket ?",
                      style: TextStyle(
                        fontFamily: DefaultTheme.fontFamily,
                        fontSize: 18,
                      ),
                    ),
                  ),
                  const SizedBox(
                    height: 20,
                  ),
                  BtnLinkWidget(
                    onPressed: () {
                      html.window.open("https://url-boutique-the-tip-top.fr", 'new tab');
                    },
                    text: "Obtenez votre numéro de ticket pour tout achat supérieur à 49€ sur notre boutique ThéTipTop",
                    fontSize: 18,
                    fontFamily: DefaultTheme.fontFamilyBold,
                  ),
                ],
              ),
            )
          ],
        ),
      ),
    );
  }
}

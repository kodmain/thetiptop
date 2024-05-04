import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class BtnLinkWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String text;
  final double? fontSize;
  final String? fontFamily;

  const BtnLinkWidget({
    super.key,
    required this.onPressed,
    required this.text,
    this.fontSize = DefaultTheme.fontSize,
    this.fontFamily = DefaultTheme.fontFamily,
  });

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onPressed,
      child: RichText(
        textAlign: TextAlign.center,
        text: TextSpan(
          style: TextStyle(
            color: DefaultTheme.textLink,
            fontSize: fontSize,
            fontFamily: fontFamily,
            decoration: TextDecoration.underline,
            decorationStyle: TextDecorationStyle.solid,
            decorationColor: DefaultTheme.textLink,
          ),
          text: text,
        ),
      ),
    );
  }
}

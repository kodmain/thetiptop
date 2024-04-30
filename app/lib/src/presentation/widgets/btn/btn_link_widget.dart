import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class BtnLinkWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String text;
  final double? fontSize;
  final String? fontFamily;

  const BtnLinkWidget({
    super.key,
    required this.onPressed,
    required this.text,
    this.fontSize = 16,
    this.fontFamily = "Raleway",
  });

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onPressed,
      style: const ButtonStyle(
        overlayColor: MaterialStatePropertyAll(
          Colors.transparent,
        ),
        tapTargetSize: MaterialTapTargetSize.shrinkWrap,
        visualDensity: VisualDensity.compact,
        minimumSize: MaterialStatePropertyAll<Size>(
          Size(0, 0),
        ),
        padding: MaterialStatePropertyAll<EdgeInsetsGeometry>(
          EdgeInsets.zero,
        ),
        shape: MaterialStatePropertyAll<OutlinedBorder>(
          RoundedRectangleBorder(
            borderRadius: BorderRadius.zero,
          ),
        ),
      ),
      child: RichText(
        textAlign: TextAlign.center,
        text: TextSpan(
          style: TextStyle(
            color: AppColor.textLink.color,
            fontSize: fontSize,
            fontFamily: fontFamily,
            decoration: TextDecoration.underline,
            decorationStyle: TextDecorationStyle.solid,
            decorationColor: AppColor.textLink.color,
          ),
          text: text,
        ),
      ),
    );
  }
}

import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class BtnSimpleWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String text;

  const BtnSimpleWidget({
    super.key,
    required this.onPressed,
    required this.text,
  });

  @override
  Widget build(BuildContext context) {
    return TextButton(
      onPressed: onPressed,
      style: ButtonStyle(
        foregroundColor: MaterialStateProperty.all<Color>(
          AppColor.whiteCream.color,
        ),
        textStyle: MaterialStatePropertyAll<TextStyle>(
          TextStyle(
            fontFamily: "Raleway",
            decoration: TextDecoration.underline,
            decorationStyle: TextDecorationStyle.solid,
            decorationColor: AppColor.whiteCream.color,
          ),
        ),
        overlayColor: MaterialStateProperty.resolveWith<Color?>((states) {
          return states.contains(MaterialState.hovered) ? const Color.fromRGBO(55, 55, 55, 1) : null;
        }),
      ),
      child: Text(
        text,
      ),
    );
  }
}

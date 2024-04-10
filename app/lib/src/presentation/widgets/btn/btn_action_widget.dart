import 'package:flutter/material.dart';

class BtnActionWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String text;
  final double? fontSize;
  final Color? backgroundColor;
  final Color? foregroundColor;
  final EdgeInsetsGeometry? padding;

  const BtnActionWidget({
    super.key,
    required this.onPressed,
    required this.text,
    this.fontSize = 18,
    this.backgroundColor,
    this.foregroundColor,
    this.padding,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: Padding(
        padding: padding ?? const EdgeInsets.fromLTRB(0, 15, 0, 15),
        child: ElevatedButton(
          onPressed: onPressed,
          style: ButtonStyle(
            backgroundColor: MaterialStatePropertyAll(backgroundColor),
            foregroundColor: MaterialStatePropertyAll(foregroundColor),
            minimumSize: const MaterialStatePropertyAll(Size(0, 65)),
            shape: MaterialStateProperty.all<RoundedRectangleBorder>(
              RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(10),
              ),
            ),
          ),
          child: Padding(
            padding: const EdgeInsets.fromLTRB(0, 10, 0, 10),
            child: Text(
              text,
              textAlign: TextAlign.center,
              style: TextStyle(
                fontSize: fontSize,
                fontFamily: "Raleway",
              ),
            ),
          ),
        ),
      ),
    );
  }
}

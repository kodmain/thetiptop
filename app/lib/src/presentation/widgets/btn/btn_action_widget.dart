import 'package:flutter/material.dart';

class BtnActionWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String text;
  final ButtonStyle? style;

  const BtnActionWidget({
    super.key,
    required this.onPressed,
    required this.text,
    this.style,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: ElevatedButton(
        onPressed: onPressed,
        style: style,
        child: Text(
          text,
          textAlign: TextAlign.center,
        ),
      ),
    );
  }
}

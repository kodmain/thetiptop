import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class SeparatorWidget extends StatelessWidget {
  final String text;

  const SeparatorWidget({
    super.key,
    required this.text,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, DefaultTheme.spacer, 0, DefaultTheme.spacer),
      child: Row(
        children: [
          Expanded(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(0, 0, DefaultTheme.spacer, 0),
              child: Container(
                height: 2,
                decoration: BoxDecoration(
                  color: DefaultTheme.blackGreenTrans,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
          ),
          Text(
            text,
            style: const TextStyle(
              color: DefaultTheme.blackGreenTrans,
              fontFamily: DefaultTheme.fontFamilySemiBold,
              fontSize: 18,
            ),
          ),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.fromLTRB(DefaultTheme.spacer, 0, 0, 0),
              child: Container(
                height: 2,
                decoration: BoxDecoration(
                  color: DefaultTheme.blackGreenTrans,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
            ),
          ),
        ],
      ),
    );
  }
}

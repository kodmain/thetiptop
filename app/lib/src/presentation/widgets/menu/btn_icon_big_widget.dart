import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class BtnIconBigWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String label;
  final String asset;
  final double? iconWidth;
  final double? iconHeight;

  const BtnIconBigWidget({
    super.key,
    required this.onPressed,
    required this.label,
    required this.asset,
    this.iconWidth = 40,
    this.iconHeight = 40,
  });

  @override
  Widget build(BuildContext context) {
    final ThemeData themeData = Theme.of(context);

    return ElevatedButton(
      onPressed: onPressed,
      style: const ButtonStyle(
          padding: MaterialStatePropertyAll(
            EdgeInsets.fromLTRB(0, DefaultTheme.spacer, 0, DefaultTheme.spacer),
          ),
          shape: MaterialStatePropertyAll(
            RoundedRectangleBorder(
              borderRadius: BorderRadius.zero,
            ),
          ),
          elevation: MaterialStatePropertyAll(0)),
      child: Column(
        children: [
          SvgPicture.asset(
            asset,
            width: iconWidth,
            height: iconHeight,
            colorFilter: const ColorFilter.mode(
              DefaultTheme.whiteCream,
              BlendMode.srcIn,
            ),
          ), // Icône
          const SizedBox(height: DefaultTheme.smallSpacer),
          Text(
            label,
            style: themeData.textTheme.displaySmall,
          ),
        ],
      ),
    );
  }
}

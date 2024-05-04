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

    return IconButton(
      onPressed: onPressed,

      /*style: ButtonStyle(
        // minimumSize: const MaterialStatePropertyAll(Size(0, 65)),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(10),
          ),
        ),
      ),*/

      icon: SvgPicture.asset(
        asset,
        width: iconWidth,
        height: iconHeight,
        colorFilter: const ColorFilter.mode(
          DefaultTheme.white,
          BlendMode.srcIn,
        ),
      ),

      /*
      icon: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          SvgPicture.asset(
            asset,
            width: iconWidth,
            height: iconHeight,
            fit: BoxFit.cover,
            colorFilter: const ColorFilter.mode(
              DefaultTheme.white,
              BlendMode.srcIn,
            ),
          ),
          const SizedBox(
            height: DefaultTheme.smallSpacer,
          ),
          Text(
            label,
            style: themeData.textTheme.bodySmall!.copyWith(
              color: themeData.colorScheme.onBackground,
            ),
          ),
        ],
      ),*/
    );
  }
}

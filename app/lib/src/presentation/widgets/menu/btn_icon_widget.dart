import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class BtnIconWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String semanticLabel;
  final String asset;
  final double? width;
  final double? height;

  const BtnIconWidget({
    super.key,
    required this.onPressed,
    required this.semanticLabel,
    required this.asset,
    this.width = 30,
    this.height = 30,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.all(DefaultTheme.smallSpacer),
      child: IconButton(
        tooltip: semanticLabel,
        onPressed: onPressed,
        style: const ButtonStyle(
          minimumSize: MaterialStatePropertyAll(Size(0, 65)),
        ),
        icon: SvgPicture.asset(
          asset,
          width: width,
          height: height,
          colorFilter: const ColorFilter.mode(
            DefaultTheme.blackGreen,
            BlendMode.srcIn,
          ),
        ),
      ),
    );
  }
}

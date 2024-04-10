import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class BtnIconWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String semanticLabel;
  final String asset;
  final double? width;
  final double? height;
  //final Color? backgroundColor;
  //final Color? foregroundColor;

  const BtnIconWidget({
    super.key,
    required this.onPressed,
    required this.semanticLabel,
    required this.asset,
    this.width = 30,
    this.height = 30,
    //this.backgroundColor,
    // this.foregroundColor,
  });

  @override
  Widget build(BuildContext context) {
    return IconButton(
      tooltip: semanticLabel,
      onPressed: onPressed,
      style: ButtonStyle(
        // backgroundColor: MaterialStatePropertyAll(backgroundColor),
        // foregroundColor: MaterialStatePropertyAll(foregroundColor),
        minimumSize: const MaterialStatePropertyAll(Size(0, 65)),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(10),
          ),
        ),
      ),
      icon: Padding(
        padding: const EdgeInsets.all(10),
        child: SvgPicture.asset(
          asset,
          width: width,
          height: height,
          colorFilter: ColorFilter.mode(
            AppColor.blackGreen.color,
            BlendMode.srcIn,
          ),
        ),
      ),
    );
  }
}

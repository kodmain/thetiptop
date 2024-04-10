import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class BtnIconBigWidget extends StatelessWidget {
  final VoidCallback onPressed;
  final String label;
  final String asset;
  final double? iconWidth;
  final double? iconHeight;
  // final Color? backgroundColor;
  //final Color? foregroundColor;

  const BtnIconBigWidget({
    super.key,
    required this.onPressed,
    required this.label,
    required this.asset,
    this.iconWidth = 40,
    this.iconHeight = 40,
    // this.backgroundColor,
    // this.foregroundColor,
  });

  @override
  Widget build(BuildContext context) {
    return IconButton(
      onPressed: onPressed,
      style: ButtonStyle(
        //  backgroundColor: MaterialStatePropertyAll(AppColor.blackGreen.color),
        //  foregroundColor: MaterialStatePropertyAll(AppColor.white.color),
        minimumSize: const MaterialStatePropertyAll(Size(0, 65)),
        shape: MaterialStateProperty.all<RoundedRectangleBorder>(
          RoundedRectangleBorder(
            borderRadius: BorderRadius.circular(10),
          ),
        ),
      ),
      icon: Padding(
        padding: const EdgeInsets.all(10),
        child: Column(
          children: [
            SvgPicture.asset(
              asset,
              width: iconWidth,
              height: iconHeight,
              colorFilter: ColorFilter.mode(
                AppColor.white.color,
                BlendMode.srcIn,
              ),
            ),
            const SizedBox(
              height: 10,
            ),
            Text(
              label,
              style: TextStyle(
                color: AppColor.white.color,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

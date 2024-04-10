import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class SeparatorWidget extends StatelessWidget {
  final String text;

  const SeparatorWidget({
    super.key,
    required this.text,
  });

  @override
  Widget build(BuildContext context) {
    return Row(
      children: [
        Expanded(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(0, 20, 20, 20),
            child: Container(
              height: 2,
              decoration: BoxDecoration(
                color: AppColor.blackGreen65.color,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
        ),
        Text(
          text,
          style: TextStyle(
            color: AppColor.blackGreen65.color,
            fontFamily: "Raleway-SemiBold",
            fontSize: 18,
          ),
        ),
        Expanded(
          child: Padding(
            padding: const EdgeInsets.fromLTRB(20, 20, 0, 20),
            child: Container(
              height: 2,
              decoration: BoxDecoration(
                color: AppColor.blackGreen65.color,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
          ),
        ),
      ],
    );
  }
}

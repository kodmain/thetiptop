import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/domain/enums/gain_status.dart';
import 'package:thetiptop_client/src/domain/enums/gain_type.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class ItemGainWidget extends StatelessWidget {
  final GainType type;
  final GainStatus status;
  final String infoNumberTicket;
  final String infoDateWin;

  const ItemGainWidget({
    super.key,
    required this.type,
    required this.status,
    required this.infoNumberTicket,
    required this.infoDateWin,
  });

  @override
  Widget build(BuildContext context) {
    return Container(
      //color: AppColor.greyLight.color,
      child: Padding(
        padding: const EdgeInsets.all(20.0),
        child: Wrap(
          alignment: WrapAlignment.spaceBetween,
          crossAxisAlignment: WrapCrossAlignment.center,
          runAlignment: WrapAlignment.center,
          spacing: 20,
          runSpacing: 20,
          children: [
            Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                SvgPicture.asset(
                  type.getAssets(),
                  colorFilter: const ColorFilter.mode(
                    DefaultTheme.greyCancel,
                    BlendMode.srcIn,
                  ),
                ),
                const SizedBox(
                  width: 10,
                ),
                Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      type.getLabel(),
                      style: const TextStyle(
                        fontFamily: "Raleway",
                        fontSize: 20,
                      ),
                    ),
                    RichText(
                      text: TextSpan(
                        children: [
                          TextSpan(
                            text: infoNumberTicket,
                            style: const TextStyle(
                              fontSize: 12,
                            ),
                          ),
                          const TextSpan(
                            text: " - ",
                            style: TextStyle(
                              fontSize: 12,
                            ),
                          ),
                          TextSpan(
                            text: infoDateWin,
                            style: const TextStyle(
                              fontSize: 12,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
              ],
            ),
            Container(
              decoration: BoxDecoration(
                color: status.getColor(),
                shape: BoxShape.rectangle,
                borderRadius: BorderRadius.circular(10),
              ),
              child: Padding(
                padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
                child: Text(
                  status.getLabel(),
                  style: const TextStyle(
                    color: DefaultTheme.whiteCream,
                    fontFamily: "Raleway",
                    fontSize: 14,
                  ),
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }
}

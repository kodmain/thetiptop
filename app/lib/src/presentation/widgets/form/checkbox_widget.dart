import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_link_widget.dart';
import 'dart:html' as html;

class CheckboxWidget extends StatelessWidget {
  final ValueChanged<bool?>? onChanged;
  final String text;
  final double? fontSize;
  final String? linkUrl;
  final String? linkLabel;
  final bool? value;

  const CheckboxWidget({
    super.key,
    required this.onChanged,
    required this.text,
    this.fontSize = 18,
    this.linkUrl,
    this.linkLabel,
    this.value,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 5, 0, 5),
      child: Row(
        children: [
          Checkbox(
            onChanged: onChanged,
            materialTapTargetSize: MaterialTapTargetSize.padded,
            value: value ?? false,
          ),
          Expanded(
            child: Wrap(
              children: [
                RichText(
                  textAlign: TextAlign.start,
                  softWrap: true,
                  text: TextSpan(
                    style: TextStyle(
                      fontSize: fontSize,
                      fontFamily: "Raleway",
                    ),
                    text: text,
                  ),
                ),
                if (linkLabel != null && linkUrl != null) ...[
                  BtnLinkWidget(
                    onPressed: () {
                      html.window.open(linkUrl!, 'new tab');
                    },
                    text: linkLabel ?? "",
                    fontSize: fontSize,
                    fontFamily: "Raleway-Bold",
                  ),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }
}

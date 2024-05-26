import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';

class BtnActionWidget extends StatelessWidget {
  final VoidCallback? onPressed;
  final String text;
  final bool isLoading;
  final bool disableLoading;

  final ButtonStyle? style;

  const BtnActionWidget({
    super.key,
    required this.onPressed,
    required this.text,
    this.isLoading = false,
    this.disableLoading = false,
    this.style,
  });

  @override
  Widget build(BuildContext context) {
    return Expanded(
      child: ElevatedButton(
        onPressed: onPressed,
        style: style,
        child: Row(
          mainAxisAlignment: MainAxisAlignment.center,
          crossAxisAlignment: CrossAxisAlignment.center,
          children: [
            if (disableLoading == false) ...[
              const SizedBox(
                width: DefaultTheme.smallSpacer + DefaultTheme.smallSpacer,
              ),
            ],
            Text(
              text,
              textAlign: TextAlign.center,
            ),
            if (disableLoading == false && isLoading == true) ...[
              const SizedBox(
                width: DefaultTheme.smallSpacer,
              ),
              SizedBox.square(
                dimension: DefaultTheme.smallSpacer,
                child: CircularProgressIndicator(
                  strokeWidth: 2,
                  strokeCap: StrokeCap.round,
                  color: Theme.of(context).colorScheme.onPrimary,
                ),
              ),
            ] else ...[
              if (disableLoading == false) ...[
                const SizedBox(
                  width: DefaultTheme.smallSpacer + DefaultTheme.smallSpacer,
                ),
              ]
            ]
          ],
        ),
      ),
    );
  }
}

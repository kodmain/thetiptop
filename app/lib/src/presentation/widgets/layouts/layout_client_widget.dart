import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/snack_bar_exception.dart';

class LayoutClientWidget extends StatelessWidget {
  final Widget child;

  const LayoutClientWidget({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    double screenWidth = MediaQuery.of(context).size.width;

    return Scaffold(
      body: Column(
        children: [
          Expanded(
            child: Row(
              crossAxisAlignment: CrossAxisAlignment.stretch,
              children: [
                if (screenWidth > 655) ...[
                  Expanded(
                    flex: 4,
                    child: Semantics(
                      excludeSemantics: true,
                      child: SizedBox(
                        width: double.infinity,
                        child: Image.asset(
                          'assets/images/parts_decor.png',
                          repeat: ImageRepeat.repeat,
                          scale: 1.75,
                        ),
                      ),
                    ),
                  ),
                ],
                Expanded(
                  flex: 10,
                  child: SizedBox(
                    width: double.infinity,
                    child: SingleChildScrollView(
                      child: Padding(
                        padding: const EdgeInsets.fromLTRB(25, 25, 25, 25),
                        child: Column(
                          children: [
                            const SnackBarException(),
                            SvgPicture.asset(
                              'assets/images/parts_logo.svg',
                              height: 150,
                              semanticsLabel: "Logo ThéTipTop",
                              colorFilter: const ColorFilter.mode(
                                DefaultTheme.blackGreen,
                                BlendMode.srcIn,
                              ),
                            ),
                            ConstrainedBox(
                              constraints: BoxConstraints.loose(
                                const Size.fromWidth(475),
                              ),
                              child: child,
                            ),
                          ],
                        ),
                      ),
                    ),
                  ),
                ),
                if (screenWidth > 655) ...[
                  Expanded(
                    flex: 4,
                    child: Semantics(
                      excludeSemantics: true,
                      child: SizedBox(
                        width: double.infinity,
                        child: Image.asset(
                          'assets/images/parts_decor.png',
                          repeat: ImageRepeat.repeat,
                          scale: 1.75,
                        ),
                      ),
                    ),
                  ),
                ],
              ],
            ),
          ),
          Container(
            height: 50,
            decoration: const BoxDecoration(
              color: DefaultTheme.greyDark,
              shape: BoxShape.rectangle,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                TextButton(
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).menuButtonTheme.style,
                  child: const Text(
                    "boutique thétiptop",
                  ),
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                TextButton(
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).menuButtonTheme.style,
                  child: const Text(
                    "règlement",
                  ),
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                TextButton(
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).menuButtonTheme.style,
                  child: const Text(
                    "cgu",
                  ),
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
                TextButton(
                  onPressed: () {
                    print("ok");
                  },
                  style: Theme.of(context).menuButtonTheme.style,
                  child: const Text(
                    "aide",
                  ),
                ),
                const SizedBox(
                  width: DefaultTheme.smallSpacer,
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

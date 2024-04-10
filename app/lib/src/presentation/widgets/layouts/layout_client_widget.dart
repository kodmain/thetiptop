import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_simple_widget.dart';

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
      backgroundColor: AppColor.whiteCream.color,
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
                            SvgPicture.asset(
                              'assets/images/parts_logo.svg',
                              height: 150,
                              semanticsLabel: "Logo ThéTipTop",
                              colorFilter: ColorFilter.mode(
                                AppColor.blackGreen.color,
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
            decoration: BoxDecoration(
              color: AppColor.greyDark.color,
              shape: BoxShape.rectangle,
            ),
            child: Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                BtnSimpleWidget(
                  onPressed: () {
                    print("ok");
                  },
                  text: "boutique ThéTipTop",
                ),
                BtnSimpleWidget(
                  onPressed: () {
                    print("ok");
                  },
                  text: "règlement",
                ),
                BtnSimpleWidget(
                  onPressed: () {
                    print("ok");
                  },
                  text: "cgu",
                ),
                BtnSimpleWidget(
                  onPressed: () {
                    print("ok");
                  },
                  text: "aide",
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

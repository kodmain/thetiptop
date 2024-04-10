import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';
import 'package:thetiptop_client/src/presentation/widgets/btn/btn_simple_widget.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_admin_widget.dart';

class LayoutAdminWidget extends StatelessWidget {
  final Widget child;

  const LayoutAdminWidget({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    double screenWidth = MediaQuery.of(context).size.width;

    return Scaffold(
      resizeToAvoidBottomInset: true,
      backgroundColor: AppColor.whiteCream.color,
      body: CustomScrollView(
        slivers: [
          SliverFillRemaining(
            hasScrollBody: false,
            child: Column(
              children: [
                Expanded(
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      SizedBox(
                        width: 200,
                        child: Container(
                          decoration: BoxDecoration(
                            color: AppColor.blackGreen.color,
                          ),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.start,
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              const SizedBox(
                                height: 25,
                              ),
                              SvgPicture.asset(
                                'assets/images/parts_logo.svg',
                                height: 60,
                                semanticsLabel: "Logo ThéTipTop",
                                colorFilter: ColorFilter.mode(
                                  AppColor.white.color,
                                  BlendMode.srcIn,
                                ),
                              ),
                              const SizedBox(
                                height: 25,
                              ),
                              MenuAdminWidget(),
                            ],
                          ),
                        ),
                      ),
                      Expanded(
                        flex: 16,
                        child: SizedBox(
                          width: double.infinity,
                          //child: SingleChildScrollView(
                          child: Padding(
                            padding: const EdgeInsets.fromLTRB(25, 25, 25, 25),
                            child: Column(
                              children: [
                                ConstrainedBox(
                                  constraints: BoxConstraints.loose(
                                    const Size.fromWidth(475),
                                  ),
                                  child: child,
                                ),
                              ],
                            ),
                          ),
                          //),
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }
}

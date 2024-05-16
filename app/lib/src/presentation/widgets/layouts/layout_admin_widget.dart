import 'package:flutter/material.dart';
import 'package:flutter_svg/svg.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'package:thetiptop_client/src/presentation/widgets/layouts/snack_bar_exception.dart';
import 'package:thetiptop_client/src/presentation/widgets/menu/menu_admin_widget.dart';

class LayoutAdminWidget extends StatelessWidget {
  final Widget child;

  const LayoutAdminWidget({
    super.key,
    required this.child,
  });

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      resizeToAvoidBottomInset: true,
      body: CustomScrollView(
        slivers: [
          SliverFillRemaining(
            hasScrollBody: false,
            child: Column(
              children: [
                const SnackBarException(),
                Expanded(
                  child: Row(
                    crossAxisAlignment: CrossAxisAlignment.stretch,
                    children: [
                      SizedBox(
                        width: 200,
                        child: Container(
                          decoration: const BoxDecoration(
                            color: DefaultTheme.blackGreen,
                          ),
                          child: Column(
                            mainAxisAlignment: MainAxisAlignment.start,
                            crossAxisAlignment: CrossAxisAlignment.center,
                            children: [
                              const SizedBox(
                                height: DefaultTheme.spacer,
                              ),
                              SvgPicture.asset(
                                'assets/images/parts_logo.svg',
                                height: 60,
                                semanticsLabel: "Logo ThéTipTop",
                                colorFilter: const ColorFilter.mode(
                                  DefaultTheme.white,
                                  BlendMode.srcIn,
                                ),
                              ),
                              const SizedBox(
                                height: 25,
                              ),
                              const MenuAdminWidget(),
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

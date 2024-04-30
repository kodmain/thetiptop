import 'package:flutter/material.dart';

/// Dans FLutter 3 les Emuns peuvent avoir des propriété,
/// et des méthodes comme n'importe quelle classe
/// https://blog.logrocket.com/deep-dive-enhanced-enums-flutter-3-0/
///
enum AppColor {
  black(Color.fromRGBO(0, 0, 0, 1)),
  blackGreen(Color.fromRGBO(51, 70, 70, 1)),
  blackGreen65(Color.fromRGBO(51, 70, 70, 0.65)),
  greyDark(Color.fromRGBO(34, 34, 34, 1)),
  greyCancel(Color.fromRGBO(92, 92, 92, 1)),
  greyLight(Color.fromRGBO(217, 217, 217, 1)),
  white(Color.fromRGBO(255, 255, 255, 1)),
  whiteCream(Color.fromRGBO(248, 238, 230, 1)),
  red(Color.fromRGBO(203, 49, 21, 1)),
  inputFill(Color.fromRGBO(217, 217, 217, 1)),
  inputBorder(Color.fromRGBO(84, 84, 84, 1)),
  textLink(Color.fromRGBO(51, 0, 255, 1)),
  blueFC(Color(0xFF2D56E7)),
  blueGC(Color(0xFF54638A));

  const AppColor(this.color);
  final Color color;
}

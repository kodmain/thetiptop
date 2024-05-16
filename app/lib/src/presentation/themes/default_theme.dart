import 'package:flutter/material.dart';

class DefaultTheme {
  static const String fontFamily = "Raleway";
  static const String fontFamilySemiBold = "Raleway-SemiBold";
  static const String fontFamilyBold = "Raleway-Bold";

  static const double fontSize = 16;
  static const double fontSizeMedium = 18;
  static const double fontSizeLarge = 22;

  static const double spacer = 30;
  static const double smallSpacer = spacer / 2;
  static const double bigSpacer = spacer * 2.5;

  static const Color white = Color(0xFFFFFFFF);
  static const Color whiteCream = Color(0xFFF8EEE6);
  static const Color greyCancel = Color(0xFF5C5C5C);
  static const Color greyDark = Color(0xFF222222);
  static const Color black = Color(0xFF000000);
  static const Color blackGreen = Color(0xFF334646);
  static const Color blackGreenTrans = Color(0xA3334646);
  static const Color inputFill = Color(0xFFD9D9D9);
  static const Color textLink = Color(0xFF3300FF);
  static const Color blueFC = Color(0xFF2D56E7);
  static const Color blueGC = Color(0xFF54638A);
  static const Color red = Color(0xFFCB3115);

  static BorderRadius standardRadius = BorderRadius.circular(10);

  static final ThemeData theme = _buildTheme();
  static ThemeData _buildTheme() {
    final ThemeData base = ThemeData.from(
      colorScheme: const ColorScheme(
        brightness: Brightness.light,
        primary: blackGreen,
        onPrimary: whiteCream,
        secondary: red,
        onSecondary: whiteCream,
        tertiary: greyCancel,
        onTertiary: whiteCream,
        error: red,
        onError: whiteCream,
        background: whiteCream,
        onBackground: black,
        surface: whiteCream,
        onSurface: black,
      ),
      useMaterial3: false,
    );

    return base.copyWith(
      // Button d'action green, basé sur elevatedButtonTheme
      elevatedButtonTheme: ElevatedButtonThemeData(
        style: ButtonStyle(
          minimumSize: const MaterialStatePropertyAll(Size(0, 65)),
          shape: MaterialStatePropertyAll(
            RoundedRectangleBorder(
              borderRadius: standardRadius,
            ),
          ),
          textStyle: const MaterialStatePropertyAll(
            TextStyle(
              fontSize: fontSizeMedium,
              fontFamily: fontFamily,
            ),
          ),
        ),
      ),

      // Button d'action rouge, basé sur elevatedButtonTheme
      outlinedButtonTheme: const OutlinedButtonThemeData(
        style: ButtonStyle(
          backgroundColor: MaterialStatePropertyAll(red),
        ),
      ),

      // Button d'action gris, basé sur elevatedButtonTheme
      filledButtonTheme: const FilledButtonThemeData(
        style: ButtonStyle(
          backgroundColor: MaterialStatePropertyAll(greyCancel),
        ),
      ),

      // Button link
      textButtonTheme: const TextButtonThemeData(
        style: ButtonStyle(
          overlayColor: MaterialStatePropertyAll(
            Colors.transparent,
          ),
          tapTargetSize: MaterialTapTargetSize.shrinkWrap,
          visualDensity: VisualDensity.compact,
          minimumSize: MaterialStatePropertyAll<Size>(
            Size(0, 0),
          ),
          padding: MaterialStatePropertyAll<EdgeInsetsGeometry>(
            EdgeInsets.zero,
          ),
          shape: MaterialStatePropertyAll<OutlinedBorder>(
            RoundedRectangleBorder(
              borderRadius: BorderRadius.zero,
            ),
          ),
        ),
      ),

      menuButtonTheme: const MenuButtonThemeData(
        style: ButtonStyle(
          foregroundColor: MaterialStatePropertyAll(whiteCream),
          textStyle: MaterialStatePropertyAll(
            TextStyle(
              fontSize: fontSize,
              fontFamily: fontFamily,
              decoration: TextDecoration.underline,
              decorationStyle: TextDecorationStyle.solid,
              decorationColor: DefaultTheme.whiteCream,
            ),
          ),
        ),
      ),

      // Input Text
      inputDecorationTheme: InputDecorationTheme(
        filled: true,
        fillColor: inputFill,
        border: OutlineInputBorder(
          borderRadius: standardRadius,
        ),
        focusedBorder: OutlineInputBorder(
          borderRadius: standardRadius,
        ),
        labelStyle: const TextStyle(
          fontWeight: FontWeight.bold,
          fontFamily: fontFamily,
        ),
        floatingLabelStyle: const TextStyle(
          fontWeight: FontWeight.bold,
          fontFamily: fontFamilyBold,
          fontSize: fontSizeMedium,
        ),
        errorMaxLines: 8,
      ),

      iconTheme: const IconThemeData(
        size: 40,
      ),

      // Text du theme
      textTheme: const TextTheme(
        // Text par default
        bodyMedium: TextStyle(
          fontSize: fontSize,
          fontFamily: fontFamily,
        ),

        displaySmall: TextStyle(
          color: whiteCream,
          fontSize: fontSize,
          fontFamily: fontFamily,
        ),

        // text gras dans des pages
        bodyLarge: TextStyle(
          color: blackGreen,
          fontSize: fontSize,
          fontFamily: fontFamilyBold,
        ),

        // title H1
        titleLarge: TextStyle(
          fontFamily: fontFamily,
          fontSize: fontSizeLarge,
        ),
      ),
    );
  }
}

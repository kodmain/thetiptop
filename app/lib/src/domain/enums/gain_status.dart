import 'package:flutter/material.dart';

/// Dans FLutter 3 les Emuns peuvent avoir des propriété,
/// et des méthodes comme n'importe quelle classe
/// https://blog.logrocket.com/deep-dive-enhanced-enums-flutter-3-0/
///
enum GainStatus {
  wait(
    {
      "color": Color.fromRGBO(189, 92, 7, 1),
      "label": "En Attente",
    },
  ),
  sent(
    {
      "color": Color.fromRGBO(203, 49, 21, 1),
      "label": "Envoyé",
    },
  ),
  recover(
    {
      "color": Color.fromRGBO(51, 70, 70, 1),
      "label": "Récupéré",
    },
  );

  const GainStatus(this.gain);
  final Map<String, dynamic> gain;

  String getLabel() {
    return gain["label"];
  }

  Color getColor() {
    return gain["color"];
  }
}

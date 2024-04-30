import 'package:flutter/material.dart';
import 'package:flutter/semantics.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/app.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';

Future<void> main() async {
  /// Flutter utilise des "bindings" pour connecter le code Dart avec le moteur de rendu sous-jacent
  /// (généralement basé sur le moteur de rendu graphique Skia). Ces bindings sont essentiels pour
  /// la communication entre le code Dart de votre application et le moteur de rendu qui exécute
  /// votre interface utilisateur.
  ///
  /// La fonction ensureInitialized() est spécifiquement utilisée pour initialiser ces bindings avant
  /// le lancement effectif de votre application. Cela peut être particulièrement important dans certaines
  /// situations, par exemple, lorsque vous utilisez des plugins ou des fonctionnalités qui dépendent de ces bindings.
  ///
  /// Exemple : Firebase, l'initialisation de base de données, ou d'autres tâches qui nécessitent que l'environnement
  /// Flutter soit correctement configuré avant leur exécution
  ///
  WidgetsFlutterBinding.ensureInitialized();

  final sharedPreferences = await SharedPreferences.getInstance();

  runApp(
    ProviderScope(
      overrides: [
        sharedPreferencesProvider.overrideWithValue(sharedPreferences),
      ],
      child: const TheTipTop(),
    ),
  );

  SemanticsBinding.instance.ensureSemantics();
}

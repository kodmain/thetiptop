import 'package:flutter/widgets.dart';
import 'package:flutter_gen/gen_l10n/app_localizations.dart';

/// Peut être utilisé pour rechercher toutes les chaînes codées en dur,
/// utile pour identifier les chaînes qui doivent être localisées.
/// https://github.com/bizz84/starter_architecture_flutter_firebase/blob/master/lib/src/localization/string_hardcoded.dart
///
extension StringHardcoded on String {
  String get hardcoded => this;
}

/// https://codewithandrea.com/articles/flutter-localization-build-context-extension/
///
extension LocalizedBuildContext on BuildContext {
  AppLocalizations get loc => AppLocalizations.of(this);
}

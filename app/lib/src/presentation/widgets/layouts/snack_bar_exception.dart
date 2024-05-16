import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/infrastructure/services/exception/viewer_exception_provider.dart';
import 'package:thetiptop_client/src/infrastructure/tools/localization/localization.dart';

class SnackBarException extends HookConsumerWidget {
  /// Affichage des exception pour l'utilisateur.
  ///
  /// Ce widget n'affiche "rien" directement, il est donc positionnable n'importe où.
  ///
  /// Il affichera une [Snackbar] par intermédiaire de [ScaffoldMessenger] à chaque
  /// modification du state de [viewerExceptionServiceProvider].
  ///
  const SnackBarException({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    ref.listen(viewerExceptionServiceProvider, (prev, next) {
      SnackBar snackBar = SnackBar(
        content: Text(context.loc.error(next!.error)),
        backgroundColor: (const Color.fromARGB(199, 255, 0, 0)),
        action: SnackBarAction(
          label: 'dismiss',
          onPressed: () {},
        ),
      );
      ScaffoldMessenger.of(context).showSnackBar(snackBar);
    });

    return const SizedBox.square(dimension: 0);
  }
}

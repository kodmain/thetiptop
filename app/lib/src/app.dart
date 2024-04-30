import 'package:flutter/material.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:thetiptop_client/src/app_router.dart';

class TheTipTop extends ConsumerWidget {
  const TheTipTop({super.key});

  @override
  Widget build(BuildContext context, WidgetRef ref) {
    return MaterialApp.router(
      showSemanticsDebugger: false,
      debugShowCheckedModeBanner: false,
      debugShowMaterialGrid: false,

      // https://api.flutter.dev/flutter/rendering/PerformanceOverlayLayer/checkerboardRasterCacheImages.html
      //
      // Le compositeur peut parfois décider de mettre en cache certaines parties
      // de la hiérarchie des widgets. De telles parties ne changent généralement
      // pas souvent d'une image à l'autre et sont coûteuses à rendre.
      // Cela peut accélérer le rendu global.
      //
      // Cependant, il y a un certain coût initial pour construire ces entrées de cache.
      // Et, si les entrées du cache ne sont pas utilisées très souvent,
      // ce coût peut ne pas valoir l'accélération du rendu des images suivantes.
      // Si le développeur veut être certain que le remplissage du cache raster
      // ne provoque pas de saccades, cette option peut être définie.
      checkerboardRasterCacheImages: false,

      // https://api.flutter.dev/flutter/widgets/Navigator/restorationScopeId.html
      //
      // ID de restauration pour restaurer l'état de l'app, y compris son historique.
      // Si un ID de restauration est fourni, le navigateur conservera son état interne
      // et le restaurera lors de la restauration de l'état. Utile lorsqu'un utilisateur
      // quitte et retourne à l'application après qu'elle a été tuée
      restorationScopeId: 'thetiptop',

      // Configuration de ThemeMode.
      theme: ThemeData(),

      // Gestion des routes, a partir de la déclaration faites dans app_router.
      //
      routerConfig: AppRouter.routes(),
    );
  }
}

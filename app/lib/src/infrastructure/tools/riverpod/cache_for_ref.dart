import 'dart:async';
import 'dart:developer';
import 'package:hooks_riverpod/hooks_riverpod.dart';

extension CacheForRef on AutoDisposeRef {
  /// Le fournisseur en vie pendant [duration] depuis sa création initiale
  /// (même si tous les auditeurs sont supprimés avant cette date)
  ///
  void cacheFor(Duration duration) {
    log('cacheFor | add $hashCode');

    /// Préserver l'état afin que la requête ne se déclenche plus
    /// si l'utilisateur quitte et rentre dans le même écran
    ///
    final link = keepAlive();

    /// Démarre un timer de temps [duration] après lequel nous invalidons le cache
    ///
    final timer = Timer(duration, () {
      log('cacheFor | invalidation du cache $hashCode');
      link.close();
    });

    /// Nous annulons le timer lorsque l'état du fournisseur est supprimé
    ///
    onDispose(() {
      log('cacheFor | dispose $hashCode');
      timer.cancel();
    });
  }
}

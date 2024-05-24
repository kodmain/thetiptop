import 'package:flutter_test/flutter_test.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:mocktail/mocktail.dart';
import 'package:thetiptop_client/src/domain/controllers/client_controller.dart';
import 'package:thetiptop_client/src/domain/models/client.dart';
import 'package:thetiptop_client/src/domain/repositories/client_repository.dart';

class MockClientRepository extends Mock implements ClientRepository {}

/// classe Listener générique,
///
/// Utilisée pour suivre le moment où un fournisseur
/// notifie ses auditeurs
///
class Listener<T> extends Mock {
  void call(T? previous, T next);
}

void main() {
  /// Méthode pour créer un ProviderContainer qui remplace ClientRepositoryProvider
  ///
  ProviderContainer makeProviderContainer(MockClientRepository clientRepository) {
    final container = ProviderContainer(
      overrides: [
        clientRepositoryProvider.overrideWithValue(clientRepository),
      ],
    );
    addTearDown(container.dispose);
    return container;
  }

  setUpAll(() {
    /// Permet d'utiliser [any] et [captureAny] sur des paramètres de type [value]
    ///
    /// Permet de fournir des valeurs de substitution par défaut pour les types génériques
    /// lors de la vérification des appels à des fonctions.
    ///
    /// Indique à mocktail qu'une instance de AsyncLoading<Client?> peut être utilisée
    /// comme valeur par défaut chaque fois que AsyncLoading<Client?> est attendu.
    ///
    /// Nous permet d'utiliser any avec un matcher (any(that: isA<AsyncLoading>()))
    /// puisque AsyncLoading != AsyncLoading avec data et que nous ne pouvons pas
    /// directement déclarer un AsyncLoading avec data
    ///
    registerFallbackValue(const AsyncLoading<Client?>());
  });

  group('initialization', () {
    test('initial state : AsyncData', () {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenAnswer((_) => Future.value());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, const AsyncData<Client?>(null)));

      /// Vérifie qu’aucun appel redondant ne se produit.
      verifyNoMoreInteractions(listener);

      /// Vérifie qu’aucune autre méthode ne sois appeler.
      verifyNever(() => clientRepository.signUp(any(), any()));
    });
  });

  group('signUp', () {
    test('success', () async {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenAnswer((_) => Future.value());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      const asyncData = AsyncData<Client?>(null);

      /// Vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncData));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(clientControllerProvider.notifier);
      await controller.signUp({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncData, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncData),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie que la méthode du repo ne sois appeler qu'une seule fois.
      ///
      verify(() => clientRepository.signUp(any(), any())).called(1);
    });

    test('error : repository', () async {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenThrow(Exception());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      const asyncData = AsyncData<Client?>(null);

      /// Vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncData));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(clientControllerProvider.notifier);
      await controller.signUp({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement a l'état d'exception (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncData, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie que la méthode du repo ne sois appeler qu'une seule fois.
      ///
      verify(() => clientRepository.signUp(any(), any())).called(1);
    });

    test('error : no password', () async {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenAnswer((_) => Future.value());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      const asyncData = AsyncData<Client?>(null);

      /// Vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncData));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(clientControllerProvider.notifier);
      await controller.signUp({'email': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement a l'état d'exception (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncData, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie que la méthode du repo ne sois jamais appeler.
      ///
      verifyNever(() => clientRepository.signUp(any(), any()));
    });

    test('error : no email', () async {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenAnswer((_) => Future.value());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      const asyncData = AsyncData<Client?>(null);

      /// Vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncData));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(clientControllerProvider.notifier);
      await controller.signUp({'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement a l'état d'exception (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncData, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie que la méthode du repo ne sois jamais appeler.
      ///
      verifyNever(() => clientRepository.signUp(any(), any()));
    });

    test('error : empty', () async {
      /// Instantiation du mock
      ///
      final clientRepository = MockClientRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => clientRepository.signUp(any(), any())).thenAnswer((_) => Future.value());

      /// Override du clientRepositoryProvider
      ///
      final container = makeProviderContainer(clientRepository);

      /// Écoute le provider a tester
      ///
      final listener = Listener<AsyncValue<Client?>>();
      container.listen(clientControllerProvider, listener, fireImmediately: true);

      const asyncData = AsyncData<Client?>(null);

      /// Vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncData));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(clientControllerProvider.notifier);
      await controller.signUp({});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement a l'état d'exception (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncData, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie que la méthode du repo ne sois jamais appeler.
      ///
      verifyNever(() => clientRepository.signUp(any(), any()));
    });
  });
}

import 'package:flutter_test/flutter_test.dart';
import 'package:hooks_riverpod/hooks_riverpod.dart';
import 'package:mocktail/mocktail.dart';
import 'package:thetiptop_client/src/domain/controllers/auth_controller.dart';
import 'package:thetiptop_client/src/domain/enums/auth_status.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository.dart';

class MockTokenRepository extends Mock implements TokenRepository {}

/// classe Listener générique,
///
/// Utilisée pour suivre le moment où un fournisseur
/// notifie ses auditeurs
///
class Listener<T> extends Mock {
  void call(T? previous, T next);
}

void main() {
  /// Méthode pour créer un ProviderContainer qui remplace tokenRepositoryProvider
  ///
  ProviderContainer makeProviderContainer(MockTokenRepository tokenRepository) {
    final container = ProviderContainer(
      overrides: [
        tokenRepositoryProvider.overrideWithValue(tokenRepository),
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
    registerFallbackValue(const Token(jwt: 'stringJwt'));
    registerFallbackValue(const AsyncLoading<AuthStatus?>());
    registerFallbackValue(const AsyncData<AuthStatus?>(AuthStatus.disconnected));
  });

  group('initialization', () {
    test('initial state : AsyncData', () {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, const AsyncData<AuthStatus?>(AuthStatus.disconnected)));

      /// Vérifie qu’aucun appel redondant ne se produit.
      verifyNoMoreInteractions(listener);

      /// Vérifie qu’aucune autre méthode ne sois appeler.
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verifyNever(() => tokenRepository.removeLocalToken());
    });
  });

  group('signIn', () {
    test('success : connected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);
      const asyncDataConnected = AsyncData<AuthStatus?>(AuthStatus.connected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataConnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verify(() => tokenRepository.getToken(any(), any())).called(1);
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verify(() => tokenRepository.saveLocalToken(any())).called(1);
      verifyNever(() => tokenRepository.removeLocalToken());
    });

    test('success : disconnected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataDisconnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verify(() => tokenRepository.getToken(any(), any())).called(1);
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verify(() => tokenRepository.removeLocalToken()).called(1);
    });

    test('error : repository', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenThrow(Exception());
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verify(() => tokenRepository.getToken(any(), any())).called(1);
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verifyNever(() => tokenRepository.removeLocalToken());
    });

    test('error : no password', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verifyNever(() => tokenRepository.removeLocalToken());
    });

    test('error : no email', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verifyNever(() => tokenRepository.removeLocalToken());
    });

    test('error : empty', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verifyNever(() => tokenRepository.removeLocalToken());
    });
  });

  group('signOut', () {
    test('success : disconnected to disconnected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);
      const asyncDataConnected = AsyncData<AuthStatus?>(AuthStatus.connected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signout();

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataDisconnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verify(() => tokenRepository.removeLocalToken()).called(1);
    });

    test('success : connected to disconnected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenAnswer((_) async => true);

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);
      const asyncDataConnected = AsyncData<AuthStatus?>(AuthStatus.connected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataConnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verify(() => tokenRepository.getToken(any(), any())).called(1);
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verify(() => tokenRepository.saveLocalToken(any())).called(1);
      verifyNever(() => tokenRepository.removeLocalToken());

      /// Execute la méthode a tester
      ///
      await controller.signout();

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataConnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataDisconnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verify(() => tokenRepository.removeLocalToken()).called(1);
    });

    test('error : repository : disconnected to disconnected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => null);
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenThrow(Exception());

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signout();

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verify(() => tokenRepository.removeLocalToken()).called(1);
    });

    test('error : repository : connected to disconnected', () async {
      /// Instantiation du mock
      ///
      final tokenRepository = MockTokenRepository();

      /// Simulation des méthodes du mocks
      ///
      when(() => tokenRepository.getToken(any(), any())).thenAnswer((_) async => const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.renewToken(any(that: isA<Token>()))).thenAnswer((_) async => null);
      when(() => tokenRepository.getLocalToken()).thenReturn(const Token(jwt: 'stringJwt'));
      when(() => tokenRepository.saveLocalToken(any(that: isA<Token>()))).thenAnswer((_) async => true);
      when(() => tokenRepository.removeLocalToken()).thenThrow(Exception());

      /// Override du tokenRepositoryProvider
      ///
      final container = makeProviderContainer(tokenRepository);

      /// Écoute le provider a tester et appelle l'écouteur [listener]
      /// chaque fois que sa valeur change
      ///
      final listener = Listener<AsyncValue<AuthStatus?>>();
      container.listen(authControllerProvider, listener, fireImmediately: true);

      const asyncDataDisconnected = AsyncData<AuthStatus?>(AuthStatus.disconnected);
      const asyncDataConnected = AsyncData<AuthStatus?>(AuthStatus.connected);

      /// vérifie la valeur initiale de la méthode de build
      /// build renvoie une valeur immédiatement, nous attendons donc AsyncData
      ///
      verify(() => listener(null, asyncDataDisconnected));

      /// Récupère le provider de notre controller et execute la méthode a tester
      ///
      final controller = container.read(authControllerProvider.notifier);
      await controller.signin({'email': 'dataString', 'password': 'dataString'});

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataDisconnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), asyncDataConnected),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verify(() => tokenRepository.getToken(any(), any())).called(1);
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verify(() => tokenRepository.saveLocalToken(any())).called(1);
      verifyNever(() => tokenRepository.removeLocalToken());

      /// Execute la méthode a tester
      ///
      await controller.signout();

      /// Vérifie dans l'ordre le déroulement des diffèrent états la méthode
      /// - transition des données à l'état de chargement (état de chargement)
      /// - transition de l'état de chargement aux données (état final une fois la connexion terminée)
      ///
      verifyInOrder([
        () => listener(asyncDataConnected, any(that: isA<AsyncLoading>())),
        () => listener(any(that: isA<AsyncLoading>()), any(that: isA<AsyncError>())),
      ]);

      /// Vérifie qu’aucun appel redondant ne se produit.
      ///
      verifyNoMoreInteractions(listener);

      /// Vérifie les différents appels de méthode du repo.
      ///
      verifyNever(() => tokenRepository.getToken(any(), any()));
      verifyNever(() => tokenRepository.renewToken(any()));
      verifyNever(() => tokenRepository.getLocalToken());
      verifyNever(() => tokenRepository.saveLocalToken(any()));
      verify(() => tokenRepository.removeLocalToken()).called(1);
    });
  });
}

import 'dart:convert';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:dio/dio.dart';
import 'package:shared_preferences/shared_preferences.dart';
import 'package:thetiptop_client/src/domain/env/env.dart';
import 'package:thetiptop_client/src/domain/models/token.dart';
import 'package:thetiptop_client/src/domain/repositories/token_repository.dart';
import 'package:thetiptop_client/src/infrastructure/providers/shared_preferences_provider.dart';
import 'package:thetiptop_client/src/infrastructure/services/dio/dio_service_provider.dart';
import 'package:thetiptop_client/src/infrastructure/services/exception/viewer_exception_provider.dart';

class MockDio extends Mock implements Dio {}

class MockViewerExceptionService extends Mock implements ViewerExceptionService {}

class MockSharedPreferences extends Mock implements SharedPreferences {}

void main() {
  late MockDio mockDio;
  late MockViewerExceptionService mockViewerExceptionService;
  late MockSharedPreferences mockSharedPreferences;
  late TokenRepository tokenRepository;

  setUpAll(() {
    registerFallbackValue(RequestOptions(path: ''));
    registerFallbackValue(Response(requestOptions: RequestOptions(path: '')));
    registerFallbackValue(DioException(requestOptions: RequestOptions(path: '')));
    registerFallbackValue(FormData());
    registerFallbackValue(const Token(jwt: ''));
  });

  setUp(() {
    mockDio = MockDio();
    mockViewerExceptionService = MockViewerExceptionService();
    mockSharedPreferences = MockSharedPreferences();
    tokenRepository = TokenRepository(
      dioService: mockDio,
      viewerExceptionService: mockViewerExceptionService,
      sharedPreferences: mockSharedPreferences,
    );
  });

  group('getToken', () {
    test('success', () async {
      when(() => mockDio.post(any(), data: any(named: 'data'))).thenAnswer((_) async {
        return Response(
          data: '{"jwt":"jwtValue"}',
          statusCode: 200,
          requestOptions: RequestOptions(path: ''),
        );
      });

      final result = await tokenRepository.getToken('email', 'password');

      /// Vérication de la reponse :
      /// - qu'elle ne soit pas vide
      /// - qu'elle soit bien transformée en object Token
      /// - que la valeur du JWT sois exate
      ///
      expect(result, isNotNull);
      expect(result, isA<Token>());
      expect(result?.jwt, equals("jwtValue"));
    });

    test('error : dio', () async {
      when(() => mockDio.post(any(), data: any(named: 'data'))).thenThrow(DioException(
        requestOptions: RequestOptions(path: ''),
        response: Response(
          requestOptions: RequestOptions(path: ''),
          statusCode: 400,
          data: 'Error',
        ),
      ));

      /// vérification d'une exception soit bien émise
      ///
      expect(
        () async => await tokenRepository.getToken('email', 'password'),
        throwsA(isA<Exception>()),
      );

      /// vérification que la methode showDioError est appeler une fois
      ///
      verify(() => mockViewerExceptionService.showDioError(any())).called(1);
    });
  });

/*
  group('renewToken', () {
    test('success', () async {
      when(() => mockDio.post(any(), data: any(named: 'data'))).thenAnswer((_) async {
        return Response(
          data: '{"jwt":"jwtValue"}',
          statusCode: 200,
          requestOptions: RequestOptions(path: ''),
        );
      });

      final result = await tokenRepository.getToken('email', 'password');

      /// Vérication de la reponse :
      /// - qu'elle ne soit pas vide
      /// - qu'elle soit bien transformée en object Token
      /// - que la valeur du JWT sois exate
      ///
      expect(result, isNotNull);
      expect(result, isA<Token>());
      expect(result?.jwt, equals("jwtValue"));
    });

    test('error : dio', () async {
      when(() => mockDio.post(any(), data: any(named: 'data'))).thenThrow(DioException(
        requestOptions: RequestOptions(path: ''),
        response: Response(
          requestOptions: RequestOptions(path: ''),
          statusCode: 400,
          data: 'Error',
        ),
      ));

      /// vérification d'une exception soit bien émise
      ///
      expect(
        () async => await tokenRepository.getToken('email', 'password'),
        throwsA(isA<Exception>()),
      );

      /// vérification que la methode showDioError est appeler une fois
      ///
      verify(() => mockViewerExceptionService.showDioError(any())).called(1);
    });
  });
  */

  group('getLocalToken', () {
    test('success : token exist', () {
      when(() => mockSharedPreferences.getString('token')).thenReturn('{"jwt":"jwtValue"}');

      final result = tokenRepository.getLocalToken();

      /// Vérication de la reponse :
      /// - qu'elle ne soit pas vide
      /// - qu'elle soit bien transformée en object Token
      /// - que la valeur du JWT sois exate
      ///
      expect(result, isNotNull);
      expect(result, isA<Token>());
      expect(result?.jwt, equals("jwtValue"));
    });

    test('success : token empty', () {
      when(() => mockSharedPreferences.getString('token')).thenReturn('{}');

      final result = tokenRepository.getLocalToken();

      /// Vérication de la reponse :
      /// - qu'elle ne soit pas vide
      /// - qu'elle soit bien transformée en object Token
      /// - que la valeur du JWT sois exate
      ///
      expect(result, isNotNull);
      expect(result, isA<Token>());
      expect(result?.jwt, equals(""));
    });

    test('success : token not exist', () {
      when(() => mockSharedPreferences.getString('token')).thenReturn(null);

      final result = tokenRepository.getLocalToken();

      /// Vérication de la reponse :
      /// - qu'elle ne soit pas vide
      /// - qu'elle soit bien transformée en object Token
      /// - que la valeur du JWT sois exate
      ///
      expect(result, isNotNull);
      expect(result, isA<Token>());
      expect(result?.jwt, equals(""));
    });
  });

  group('saveLocalToken', () {
    test('success', () async {
      when(() => mockSharedPreferences.setString(any(), any())).thenAnswer((_) async => true);

      const token = Token(jwt: 'jwtValue');
      await tokenRepository.saveLocalToken(token);

      /// vérification que la methode setString est appeler une fois
      ///
      verify(() => mockSharedPreferences.setString('token', jsonEncode(token.toJson()))).called(1);
    });
  });

  group('removeLocalToken', () {
    test('success', () {
      when(() => mockSharedPreferences.remove(any())).thenAnswer((_) async => true);
      tokenRepository.removeLocalToken();

      /// vérification que la methode setString est appeler une fois
      ///
      verify(() => mockSharedPreferences.remove('token')).called(1);
    });
  });
}

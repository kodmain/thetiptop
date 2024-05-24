import 'package:flutter_test/flutter_test.dart';
import 'package:thetiptop_client/src/infrastructure/tools/form/validator.dart';

void main() {
  group('notEmpty', () {
    test('error : value null', () {
      final result = Validator().notEmpty(value: null);
      expect(result, 'Ce champs est obligatoire');
    });

    test('error : value empty', () {
      final result = Validator().notEmpty(value: '');
      expect(result, 'Ce champs est obligatoire');
    });

    test('error with custom message', () {
      final result = Validator().notEmpty(value: '', message: 'Custom message');
      expect(result, 'Custom message');
    });

    test('success', () {
      final result = Validator().notEmpty(value: 'test');
      expect(result, null);
    });
  });

  group('isEmail', () {
    test('error : value not email', () {
      final result = Validator().isEmail(value: 'not-an-email');
      expect(result, 'Veuillez entrer votre email');
    });

    test('error : value null', () {
      final result = Validator().isEmail(value: null);
      expect(result, 'Veuillez entrer votre email');
    });

    test('error : value  empty', () {
      final result = Validator().isEmail(value: '');
      expect(result, 'Veuillez entrer votre email');
    });

    test('error with custom message', () {
      final result = Validator().isEmail(value: '', message: 'Custom message');
      expect(result, 'Custom message');
    });

    test('error : list bad emails', () {
      List<String> invalidEmails = [
        'user@@example.com', // Double @ symbol
        'user@.example.com', // Leading dot in domain
        'user@example..com', // Double dot in domain
        'user@example,com', // Comma in domain
        'user@-example.com', // Leading hyphen in domain
        'user@example-.com', // Trailing hyphen in domain
        'user@.com', // Domain starts with a dot
        'user@example.c', // Single character TLD
        'user@example.abcdefghijk', // Too long TLD
        'user@ex..ample.com', // Double dot in domain
        'user@exam_ple.com', // Underscore in domain
        'user@exam!ple.com', // Exclamation mark in domain
        'user@ex&ample.com', // Ampersand in domain
        'user@exam#ple.com', // Hash in domain
        'user@ex*ample.com', // Asterisk in domain
        'user@exam(ple.com', // Parenthesis in domain
        'user@exam)ple.com', // Parenthesis in domain
        'user@exam=ple.com', // Equals sign in domain
        'user@exam+ple.com', // Plus sign in domain
        'user@exam{ple.com', // Curly brace in domain
        'user@exam}ple.com', // Curly brace in domain
        'user@exam^ple.com', // Caret in domain
        'user@example.com.', // Trailing dot in domain
        'user.@example.com', // Trailing dot in local part
        '.user@example.com', // Leading dot in local part
        'user..name@example.com', // Double dot in local part
        'user@example.com (Joe Smith)', // Text after address
        'user@example..com', // Double dot in domain
        'user@example.c', // Single character TLD
        'user@example,com', // Comma in domain
        'user@.example.com', // Leading dot in domain
        'user@example@example.com', // Multiple @ symbols
        'user@example@example@example.com', // Multiple @ symbols
        'user@example.com@example', // Multiple @ symbols
        'user@example.com,', // Comma after domain
        'user@example..com', // Double dot in domain
        'user@example.com (comment)', // Comment after domain
        'user@ example.com', // Space in domain
        'user@example .com', // Space in domain
        'user@exa mple.com', // Space in domain
        'user@.example.com', // Leading dot in domain
        'user@example.c', // Single character TLD
        'user@ex_mple.com', // Underscore in domain
        'user@exam!ple.com', // Exclamation mark in domain
        'user@exam#ple.com', // Hash in domain
        'user@example-.com', // Trailing hyphen in domain
        'user@-example.com', // Leading hyphen in domain
        'user@localhost', // Localhost domain - techniquement valide
        'user@[127.0.0.1]', // IP address domain - techniquement valide
        'user@[IPv6:2001:0db8:85a3:0000:0000:8a2e:0370:7334]', // IPv6 address domain - techniquement valide
      ];

      for (var email in invalidEmails) {
        final result = Validator().isEmail(value: email);
        expect(result, 'Veuillez entrer votre email');
      }
    });

    test('success', () {
      final result = Validator().notEmpty(value: 'test@example.com');
      expect(result, null);
    });

    test('success : list good emails', () {
      final validEmails = [
        'user-123@example.com',
        'user123@example.com',
        'user@example.com',
        'user@example.com',
        'user.name@example.com',
        'user-name@example.com',
        'user+name@example.com',
        'user_name@example.com',
        'user.name+tag+sorting@example.com',
        'user.name@sub.example.com',
        'user@example.co.uk',
        'user@example.org',
        'user@example.net',
        'user@example.io',
      ];

      for (var email in validEmails) {
        final result = Validator().isEmail(value: email);
        expect(result, null);
      }
    });
  });

  group('isTrue', () {
    test('error : value false', () {
      final result = Validator().isTrue(value: false);
      expect(result, 'Cette confirmation est obligatoire');
    });

    test('error : value null', () {
      final result = Validator().isTrue(value: null);
      expect(result, 'Cette confirmation est obligatoire');
    });

    test('error : value empty', () {
      final result = Validator().isTrue(value: '');
      expect(result, 'Cette confirmation est obligatoire');
    });

    test('error : value int', () {
      final result = Validator().isTrue(value: 1);
      expect(result, 'Cette confirmation est obligatoire');
    });

    test('error : value string', () {
      final result = Validator().isTrue(value: 'notEmpty');
      expect(result, 'Cette confirmation est obligatoire');
    });

    test('error with custom message', () {
      final result = Validator().isTrue(value: false, message: 'Custom message');
      expect(result, 'Custom message');
    });

    test('success', () {
      final result = Validator().isTrue(value: true);
      expect(result, null);
    });
  });

  group('isComplex', () {
    test('error : only minuscule', () {
      final result = Validator().isComplex(value: 'password');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only majuscule', () {
      final result = Validator().isComplex(value: 'PASSWORD');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only int', () {
      final result = Validator().isComplex(value: '12345678');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only special char', () {
      final result = Validator().isComplex(value: '/-*.?;:=');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only string', () {
      final result = Validator().isComplex(value: 'passWORD');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only string and int', () {
      final result = Validator().isComplex(value: 'passW0RD');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : only string and special char', () {
      final result = Validator().isComplex(value: 'pass-WORD');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error : too short', () {
      final result = Validator().isComplex(value: 'ps-0R');
      expect(result, "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres");
    });

    test('error with custom message', () {
      final result = Validator().isComplex(value: 'password', message: 'Custom message');
      expect(result, "Custom message");
    });

    test('success', () {
      final result = Validator().isComplex(value: 'pass-W0RD');
      expect(result, null);
    });
  });

  group('isEqual', () {
    test('error : string values are not equal', () {
      final result = Validator().isEqual(firstValue: 'value1', secondValue: 'value2');
      expect(result, 'Valeur inattendue');
    });

    test('error : string values are not equal', () {
      final result = Validator().isEqual(firstValue: null, secondValue: '');
      expect(result, 'Valeur inattendue');
    });

    test('error : int values are not equal', () {
      final result = Validator().isEqual(firstValue: 1, secondValue: 2);
      expect(result, 'Valeur inattendue');
    });

    test('error with custom message', () {
      final result = Validator().isEqual(firstValue: 'value1', secondValue: 'value2', message: 'Custom message');
      expect(result, 'Custom message');
    });

    test('success', () {
      final result = Validator().isEqual(firstValue: 'value', secondValue: 'value');
      expect(result, null);
    });
  });
}

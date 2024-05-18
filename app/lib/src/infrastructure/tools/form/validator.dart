class Validator {
  static const String regexrComplexFormat = r"^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\d\s:])([^\s]){8,}$";
  static const String regexrEmailFormat = r"^[a-zA-Z0-9!#$%&*+/=?^_|~-]+(\.[a-zA-Z0-9!#$%&*+/=?^_|~-]+)*@([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$";

  /// Vérifie qu'une valeur ne soit pas null ou vide.
  ///
  String? notEmpty({value, message}) {
    if (value == null || value.isEmpty) {
      return message ?? "Ce champs est obligatoire";
    }
    return null;
  }

  /// Vérifie qu'une valeur soit un email.
  ///
  String? isEmail({value, message, regexp = regexrEmailFormat}) {
    if (!RegExp(regexp).hasMatch(value ?? "")) {
      return message ?? "Veuillez entrer votre email";
    }

    return null;
  }

  /// Vérifie qu'une valeur soit égale à true.
  ///
  String? isTrue({value, message}) {
    if (value != true) {
      return message ?? "Cette confirmation est obligatoire";
    }

    return null;
  }

  /// Vérifie la complexité d'une valeur.
  ///
  String? isComplex({value, message, regexp = regexrComplexFormat}) {
    if (!RegExp(regexp).hasMatch(value ?? "")) {
      return message ?? "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres";
    }

    return null;
  }

  /// Vérifie l'égalité entre deux valeur.
  ///
  String? isEqual({firstValue, secondValue, message}) {
    if (firstValue != secondValue) {
      return message ?? "Valeur inattendue";
    }

    return null;
  }
}

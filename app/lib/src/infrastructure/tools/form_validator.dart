class FormValidator {
  static const String regexrComplexFormat = r"^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^\w\d\s:])([^\s]){8,}$";
  static const String regexrEmailFormat = r"^[a-zA-Z0-9!#$%&*+/=?^_|~-]+(\.[a-zA-Z0-9!#$%&*+/=?^_|~-]+)*@([a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?\.)+[a-zA-Z0-9]([a-zA-Z0-9-]*[a-zA-Z0-9])?$";

  String? notEmpty({value, message}) {
    if (value == null || value.isEmpty) {
      return message ?? "Ce champs est obligatoire";
    }
    return null;
  }

  String? isEmail({value, message}) {
    if (!RegExp(regexrEmailFormat).hasMatch(value ?? "")) {
      return message ?? "Veuillez entrer votre email";
    }

    return null;
  }

  String? isTrue({value, message}) {
    print(value);

    if (value != true) {
      return message ?? "Cette confirmation est obligatoire";
    }

    return null;
  }

  String? isComplex({value, message}) {
    if (!RegExp(regexrComplexFormat).hasMatch(value ?? "")) {
      return message ?? "Votre mot de passe doit faire au minimum 8 caractères et être constitué : \r - De lettres minuscules \r - De lettres majuscules \r - De caractères spéciaux  \r - De chiffres";
    }

    return null;
  }

  String? isEqual({firstValue, secondValue, message}) {
    if (firstValue != secondValue) {
      return message ?? "Valeur inattendue";
    }

    return null;
  }
}

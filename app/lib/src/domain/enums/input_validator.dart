/// Dans FLutter 3 les Emuns peuvent avoir des propriété,
/// et des méthodes comme n'importe quelle classe
/// https://blog.logrocket.com/deep-dive-enhanced-enums-flutter-3-0/
///
enum InputValidator {
  email("email"),
  notEmpty("notEmpty"),
  equal("equal");

  const InputValidator(this.typeValidator);
  final String typeValidator;
}

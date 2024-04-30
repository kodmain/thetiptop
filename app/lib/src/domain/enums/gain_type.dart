/// Dans FLutter 3 les Emuns peuvent avoir des propriété,
/// et des méthodes comme n'importe quelle classe
/// https://blog.logrocket.com/deep-dive-enhanced-enums-flutter-3-0/
///
enum GainType {
  signature(
    {
      "assets": "assets/images/parts_lot-the-signature.svg",
      "label": "Thé Signature",
    },
  ),
  detox(
    {
      "assets": "assets/images/parts_lot-the-detox.svg",
      "label": "Thé Détox",
    },
  ),
  infuseur(
    {
      "assets": "assets/images/parts_lot-infuseur.svg",
      "label": "Infuseur à Thé",
    },
  ),
  smallCoffret(
    {
      "assets": "assets/images/parts_lot-small-coffret.svg",
      "label": "Coffret Découverte",
    },
  ),
  bigCoffret(
    {
      "assets": "assets/images/parts_lot-big-coffret.svg",
      "label": "Coffret Découverte+",
    },
  );

  const GainType(this.gain);
  final Map<String, dynamic> gain;

  String getLabel() {
    return gain["label"];
  }

  String getAssets() {
    return gain["assets"];
  }
}

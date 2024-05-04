import 'package:flutter/gestures.dart';
import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/presentation/themes/default_theme.dart';
import 'dart:html' as html;

class CheckboxFormField extends FormField<bool> {
  CheckboxFormField({
    Key? key,
    required String textLabel,
    TextStyle? textStyle,
    String? linkUrl,
    String? linkLabel,
    ValueChanged<bool>? onChanged,
    FormFieldSetter<bool>? onSaved,
    FormFieldValidator<bool>? validator,
    AutovalidateMode autovalidateMode = AutovalidateMode.disabled,
    bool? initialValue,
    bool value = false,
    bool tristate = false,
    String? errorText,
  }) : super(
          key: key,
          initialValue: initialValue ?? false,
          onSaved: onSaved,
          validator: validator,
          autovalidateMode: autovalidateMode,
          builder: (FormFieldState<bool> field) {
            final ThemeData themeData = Theme.of(field.context);

            return Column(
              children: [
                Row(
                  children: [
                    Checkbox(
                      side: (field.errorText != null && !field.isValid)
                          ? BorderSide(
                              color: themeData.colorScheme.error,
                              width: 2,
                            )
                          : BorderSide(
                              color: themeData.colorScheme.surfaceTint,
                              width: 2,
                            ),
                      materialTapTargetSize: MaterialTapTargetSize.padded,
                      value: field.value,
                      onChanged: field.didChange,
                      tristate: tristate,
                    ),
                    Expanded(
                      child: Wrap(
                        children: [
                          RichText(
                            textAlign: TextAlign.start,
                            softWrap: true,
                            text: TextSpan(
                              style: textStyle,
                              children: [
                                TextSpan(
                                  text: textLabel,
                                ),
                                if (linkLabel != null && linkUrl != null) ...[
                                  TextSpan(
                                    text: linkLabel,
                                    style: const TextStyle(
                                      color: DefaultTheme.textLink,
                                      decoration: TextDecoration.underline,
                                      decorationStyle: TextDecorationStyle.solid,
                                      decorationColor: DefaultTheme.textLink,
                                    ),
                                    recognizer: TapGestureRecognizer()
                                      ..onTapDown = (event) {
                                        html.window.open(linkUrl!, 'new tab');
                                      },
                                  ),
                                ]
                              ],
                            ),
                          ),
                        ],
                      ),
                    ),
                  ],
                ),
                if (field.errorText != null && !field.isValid) ...[
                  Row(
                    children: [
                      Padding(
                        padding: const EdgeInsets.fromLTRB(10, 0, 0, 15),
                        child: Text(
                          field.errorText ?? "Champ requis",
                          style: themeData.textTheme.bodySmall!.copyWith(color: themeData.colorScheme.error),
                        ),
                      ),
                    ],
                  ),
                ],
              ],
            );
          },
        );
}

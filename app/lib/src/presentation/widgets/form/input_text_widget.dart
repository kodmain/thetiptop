import 'package:flutter/material.dart';
import 'package:thetiptop_client/src/domain/enums/app_color.dart';

class InputText extends StatelessWidget {
  final FormFieldValidator<String>? validator;
  final bool obscureText;
  final String labelText;
  final String? hintText;
  final String? initialValue;

  const InputText({
    super.key,
    this.validator,
    this.obscureText = false,
    required this.labelText,
    this.hintText,
    this.initialValue,
  });

  @override
  Widget build(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(0, 15, 0, 15),
      child: TextFormField(
        initialValue: initialValue,
        autovalidateMode: AutovalidateMode.onUserInteraction,
        obscureText: obscureText,
        cursorColor: AppColor.blackGreen.color,
        decoration: InputDecoration(
          filled: true,
          fillColor: AppColor.inputFill.color,
          focusColor: AppColor.inputFill.color,
          focusedBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColor.inputBorder.color),
            borderRadius: BorderRadius.circular(10.0),
          ),
          enabledBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColor.inputBorder.color),
            borderRadius: BorderRadius.circular(10.0),
          ),
          errorBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColor.red.color),
            borderRadius: BorderRadius.circular(10.0),
          ),
          focusedErrorBorder: OutlineInputBorder(
            borderSide: BorderSide(color: AppColor.red.color),
            borderRadius: BorderRadius.circular(10.0),
          ),
          hintText: hintText,
          labelText: labelText,
          labelStyle: TextStyle(
            color: AppColor.blackGreen.color,
            fontWeight: FontWeight.bold,
            fontFamily: "Raleway",
          ),
          floatingLabelStyle: TextStyle(
            color: AppColor.blackGreen.color,
            fontWeight: FontWeight.bold,
            fontFamily: "Raleway-Bold",
            fontSize: 18,
          ),
          errorMaxLines: 8,
        ),
        validator: validator,
      ),
    );
  }
}

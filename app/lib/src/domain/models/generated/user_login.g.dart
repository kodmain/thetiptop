// GENERATED CODE - DO NOT MODIFY BY HAND

part of '../user_login.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$UserLoginImpl _$$UserLoginImplFromJson(Map<String, dynamic> json) =>
    _$UserLoginImpl(
      email: Input.fromJson(json['email'] as Map<String, dynamic>),
      password: Input.fromJson(json['password'] as Map<String, dynamic>),
    );

Map<String, dynamic> _$$UserLoginImplToJson(_$UserLoginImpl instance) =>
    <String, dynamic>{
      'email': instance.email.toJson(),
      'password': instance.password.toJson(),
    };

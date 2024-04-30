// GENERATED CODE - DO NOT MODIFY BY HAND

part of '../input.dart';

// **************************************************************************
// JsonSerializableGenerator
// **************************************************************************

_$InputImpl _$$InputImplFromJson(Map<String, dynamic> json) => _$InputImpl(
      value: json['value'] as String,
      error: json['error'] as String? ?? "",
      isValid: json['is_valid'] as bool? ?? false,
    );

Map<String, dynamic> _$$InputImplToJson(_$InputImpl instance) =>
    <String, dynamic>{
      'value': instance.value,
      'error': instance.error,
      'is_valid': instance.isValid,
    };

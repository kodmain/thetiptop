// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of '../input.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

Input _$InputFromJson(Map<String, dynamic> json) {
  return _Input.fromJson(json);
}

/// @nodoc
mixin _$Input {
  String get value => throw _privateConstructorUsedError;
  String get error => throw _privateConstructorUsedError;
  bool get isValid => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $InputCopyWith<Input> get copyWith => throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $InputCopyWith<$Res> {
  factory $InputCopyWith(Input value, $Res Function(Input) then) =
      _$InputCopyWithImpl<$Res, Input>;
  @useResult
  $Res call({String value, String error, bool isValid});
}

/// @nodoc
class _$InputCopyWithImpl<$Res, $Val extends Input>
    implements $InputCopyWith<$Res> {
  _$InputCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? value = null,
    Object? error = null,
    Object? isValid = null,
  }) {
    return _then(_value.copyWith(
      value: null == value
          ? _value.value
          : value // ignore: cast_nullable_to_non_nullable
              as String,
      error: null == error
          ? _value.error
          : error // ignore: cast_nullable_to_non_nullable
              as String,
      isValid: null == isValid
          ? _value.isValid
          : isValid // ignore: cast_nullable_to_non_nullable
              as bool,
    ) as $Val);
  }
}

/// @nodoc
abstract class _$$InputImplCopyWith<$Res> implements $InputCopyWith<$Res> {
  factory _$$InputImplCopyWith(
          _$InputImpl value, $Res Function(_$InputImpl) then) =
      __$$InputImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({String value, String error, bool isValid});
}

/// @nodoc
class __$$InputImplCopyWithImpl<$Res>
    extends _$InputCopyWithImpl<$Res, _$InputImpl>
    implements _$$InputImplCopyWith<$Res> {
  __$$InputImplCopyWithImpl(
      _$InputImpl _value, $Res Function(_$InputImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? value = null,
    Object? error = null,
    Object? isValid = null,
  }) {
    return _then(_$InputImpl(
      value: null == value
          ? _value.value
          : value // ignore: cast_nullable_to_non_nullable
              as String,
      error: null == error
          ? _value.error
          : error // ignore: cast_nullable_to_non_nullable
              as String,
      isValid: null == isValid
          ? _value.isValid
          : isValid // ignore: cast_nullable_to_non_nullable
              as bool,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$InputImpl implements _Input {
  const _$InputImpl(
      {required this.value, this.error = "", this.isValid = false});

  factory _$InputImpl.fromJson(Map<String, dynamic> json) =>
      _$$InputImplFromJson(json);

  @override
  final String value;
  @override
  @JsonKey()
  final String error;
  @override
  @JsonKey()
  final bool isValid;

  @override
  String toString() {
    return 'Input(value: $value, error: $error, isValid: $isValid)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$InputImpl &&
            (identical(other.value, value) || other.value == value) &&
            (identical(other.error, error) || other.error == error) &&
            (identical(other.isValid, isValid) || other.isValid == isValid));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, value, error, isValid);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$InputImplCopyWith<_$InputImpl> get copyWith =>
      __$$InputImplCopyWithImpl<_$InputImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$InputImplToJson(
      this,
    );
  }
}

abstract class _Input implements Input {
  const factory _Input(
      {required final String value,
      final String error,
      final bool isValid}) = _$InputImpl;

  factory _Input.fromJson(Map<String, dynamic> json) = _$InputImpl.fromJson;

  @override
  String get value;
  @override
  String get error;
  @override
  bool get isValid;
  @override
  @JsonKey(ignore: true)
  _$$InputImplCopyWith<_$InputImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

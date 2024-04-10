// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of '../user_login.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

UserLogin _$UserLoginFromJson(Map<String, dynamic> json) {
  return _UserLogin.fromJson(json);
}

/// @nodoc
mixin _$UserLogin {
  Input get email => throw _privateConstructorUsedError;
  Input get password => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $UserLoginCopyWith<UserLogin> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UserLoginCopyWith<$Res> {
  factory $UserLoginCopyWith(UserLogin value, $Res Function(UserLogin) then) =
      _$UserLoginCopyWithImpl<$Res, UserLogin>;
  @useResult
  $Res call({Input email, Input password});

  $InputCopyWith<$Res> get email;
  $InputCopyWith<$Res> get password;
}

/// @nodoc
class _$UserLoginCopyWithImpl<$Res, $Val extends UserLogin>
    implements $UserLoginCopyWith<$Res> {
  _$UserLoginCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? email = null,
    Object? password = null,
  }) {
    return _then(_value.copyWith(
      email: null == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as Input,
      password: null == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as Input,
    ) as $Val);
  }

  @override
  @pragma('vm:prefer-inline')
  $InputCopyWith<$Res> get email {
    return $InputCopyWith<$Res>(_value.email, (value) {
      return _then(_value.copyWith(email: value) as $Val);
    });
  }

  @override
  @pragma('vm:prefer-inline')
  $InputCopyWith<$Res> get password {
    return $InputCopyWith<$Res>(_value.password, (value) {
      return _then(_value.copyWith(password: value) as $Val);
    });
  }
}

/// @nodoc
abstract class _$$UserLoginImplCopyWith<$Res>
    implements $UserLoginCopyWith<$Res> {
  factory _$$UserLoginImplCopyWith(
          _$UserLoginImpl value, $Res Function(_$UserLoginImpl) then) =
      __$$UserLoginImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({Input email, Input password});

  @override
  $InputCopyWith<$Res> get email;
  @override
  $InputCopyWith<$Res> get password;
}

/// @nodoc
class __$$UserLoginImplCopyWithImpl<$Res>
    extends _$UserLoginCopyWithImpl<$Res, _$UserLoginImpl>
    implements _$$UserLoginImplCopyWith<$Res> {
  __$$UserLoginImplCopyWithImpl(
      _$UserLoginImpl _value, $Res Function(_$UserLoginImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? email = null,
    Object? password = null,
  }) {
    return _then(_$UserLoginImpl(
      email: null == email
          ? _value.email
          : email // ignore: cast_nullable_to_non_nullable
              as Input,
      password: null == password
          ? _value.password
          : password // ignore: cast_nullable_to_non_nullable
              as Input,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UserLoginImpl extends _UserLogin {
  const _$UserLoginImpl({required this.email, required this.password})
      : super._();

  factory _$UserLoginImpl.fromJson(Map<String, dynamic> json) =>
      _$$UserLoginImplFromJson(json);

  @override
  final Input email;
  @override
  final Input password;

  @override
  String toString() {
    return 'UserLogin(email: $email, password: $password)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UserLoginImpl &&
            (identical(other.email, email) || other.email == email) &&
            (identical(other.password, password) ||
                other.password == password));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, email, password);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$UserLoginImplCopyWith<_$UserLoginImpl> get copyWith =>
      __$$UserLoginImplCopyWithImpl<_$UserLoginImpl>(this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UserLoginImplToJson(
      this,
    );
  }
}

abstract class _UserLogin extends UserLogin {
  const factory _UserLogin(
      {required final Input email,
      required final Input password}) = _$UserLoginImpl;
  const _UserLogin._() : super._();

  factory _UserLogin.fromJson(Map<String, dynamic> json) =
      _$UserLoginImpl.fromJson;

  @override
  Input get email;
  @override
  Input get password;
  @override
  @JsonKey(ignore: true)
  _$$UserLoginImplCopyWith<_$UserLoginImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

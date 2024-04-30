// coverage:ignore-file
// GENERATED CODE - DO NOT MODIFY BY HAND
// ignore_for_file: type=lint
// ignore_for_file: unused_element, deprecated_member_use, deprecated_member_use_from_same_package, use_function_type_syntax_for_parameters, unnecessary_const, avoid_init_to_null, invalid_override_different_default_values_named, prefer_expression_function_bodies, annotate_overrides, invalid_annotation_target, unnecessary_question_mark

part of '../user_login_form_state.dart';

// **************************************************************************
// FreezedGenerator
// **************************************************************************

T _$identity<T>(T value) => value;

final _privateConstructorUsedError = UnsupportedError(
    'It seems like you constructed your class using `MyClass._()`. This constructor is only meant to be used by freezed and you are not supposed to need it nor use it.\nPlease check the documentation here for more information: https://github.com/rrousselGit/freezed#adding-getters-and-methods-to-our-models');

UserLoginFormState _$UserLoginFormStateFromJson(Map<String, dynamic> json) {
  return _UserLoginFormState.fromJson(json);
}

/// @nodoc
mixin _$UserLoginFormState {
  UserLogin get form => throw _privateConstructorUsedError;

  Map<String, dynamic> toJson() => throw _privateConstructorUsedError;
  @JsonKey(ignore: true)
  $UserLoginFormStateCopyWith<UserLoginFormState> get copyWith =>
      throw _privateConstructorUsedError;
}

/// @nodoc
abstract class $UserLoginFormStateCopyWith<$Res> {
  factory $UserLoginFormStateCopyWith(
          UserLoginFormState value, $Res Function(UserLoginFormState) then) =
      _$UserLoginFormStateCopyWithImpl<$Res, UserLoginFormState>;
  @useResult
  $Res call({UserLogin form});

  $UserLoginCopyWith<$Res> get form;
}

/// @nodoc
class _$UserLoginFormStateCopyWithImpl<$Res, $Val extends UserLoginFormState>
    implements $UserLoginFormStateCopyWith<$Res> {
  _$UserLoginFormStateCopyWithImpl(this._value, this._then);

  // ignore: unused_field
  final $Val _value;
  // ignore: unused_field
  final $Res Function($Val) _then;

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? form = null,
  }) {
    return _then(_value.copyWith(
      form: null == form
          ? _value.form
          : form // ignore: cast_nullable_to_non_nullable
              as UserLogin,
    ) as $Val);
  }

  @override
  @pragma('vm:prefer-inline')
  $UserLoginCopyWith<$Res> get form {
    return $UserLoginCopyWith<$Res>(_value.form, (value) {
      return _then(_value.copyWith(form: value) as $Val);
    });
  }
}

/// @nodoc
abstract class _$$UserLoginFormStateImplCopyWith<$Res>
    implements $UserLoginFormStateCopyWith<$Res> {
  factory _$$UserLoginFormStateImplCopyWith(_$UserLoginFormStateImpl value,
          $Res Function(_$UserLoginFormStateImpl) then) =
      __$$UserLoginFormStateImplCopyWithImpl<$Res>;
  @override
  @useResult
  $Res call({UserLogin form});

  @override
  $UserLoginCopyWith<$Res> get form;
}

/// @nodoc
class __$$UserLoginFormStateImplCopyWithImpl<$Res>
    extends _$UserLoginFormStateCopyWithImpl<$Res, _$UserLoginFormStateImpl>
    implements _$$UserLoginFormStateImplCopyWith<$Res> {
  __$$UserLoginFormStateImplCopyWithImpl(_$UserLoginFormStateImpl _value,
      $Res Function(_$UserLoginFormStateImpl) _then)
      : super(_value, _then);

  @pragma('vm:prefer-inline')
  @override
  $Res call({
    Object? form = null,
  }) {
    return _then(_$UserLoginFormStateImpl(
      null == form
          ? _value.form
          : form // ignore: cast_nullable_to_non_nullable
              as UserLogin,
    ));
  }
}

/// @nodoc
@JsonSerializable()
class _$UserLoginFormStateImpl implements _UserLoginFormState {
  const _$UserLoginFormStateImpl(this.form);

  factory _$UserLoginFormStateImpl.fromJson(Map<String, dynamic> json) =>
      _$$UserLoginFormStateImplFromJson(json);

  @override
  final UserLogin form;

  @override
  String toString() {
    return 'UserLoginFormState(form: $form)';
  }

  @override
  bool operator ==(Object other) {
    return identical(this, other) ||
        (other.runtimeType == runtimeType &&
            other is _$UserLoginFormStateImpl &&
            (identical(other.form, form) || other.form == form));
  }

  @JsonKey(ignore: true)
  @override
  int get hashCode => Object.hash(runtimeType, form);

  @JsonKey(ignore: true)
  @override
  @pragma('vm:prefer-inline')
  _$$UserLoginFormStateImplCopyWith<_$UserLoginFormStateImpl> get copyWith =>
      __$$UserLoginFormStateImplCopyWithImpl<_$UserLoginFormStateImpl>(
          this, _$identity);

  @override
  Map<String, dynamic> toJson() {
    return _$$UserLoginFormStateImplToJson(
      this,
    );
  }
}

abstract class _UserLoginFormState implements UserLoginFormState {
  const factory _UserLoginFormState(final UserLogin form) =
      _$UserLoginFormStateImpl;

  factory _UserLoginFormState.fromJson(Map<String, dynamic> json) =
      _$UserLoginFormStateImpl.fromJson;

  @override
  UserLogin get form;
  @override
  @JsonKey(ignore: true)
  _$$UserLoginFormStateImplCopyWith<_$UserLoginFormStateImpl> get copyWith =>
      throw _privateConstructorUsedError;
}

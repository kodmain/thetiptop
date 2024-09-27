package services

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"github.com/kodmain/thetiptop/api/internal/domain/user/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
	serializer "github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

func UserAuth(service services.UserServiceInterface, credentialDTO *transfert.Credential) (int, any) {
	if err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	}); err != nil {
		return err.Code(), err
	}

	credentialID, role, err := service.UserAuth(credentialDTO)
	if err != nil {
		return err.Code(), err
	}

	accessToken, refreshToken, err := serializer.FromID(*credentialID, map[string]any{
		"role": role,
	})

	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func UserAuthRenew(refresh *serializer.Token) (int, any) {
	var err errors.ErrorInterface = errors.ErrBadRequest
	if refresh == nil {
		err.WithData("missing token")
		return err.Code(), err
	}

	err = errors.ErrUnauthorized
	if refresh.Type != serializer.REFRESH {
		err.WithData("invalid token type")
		return err.Code(), err
	}

	if refresh.HasExpired() {
		err.WithData("refresh token has expired")
		return err.Code(), err
	}

	accessToken, refreshToken, err := serializer.FromID(refresh.ID, refresh.Data)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, fiber.Map{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	}
}

func CredentialUpdate(service services.UserServiceInterface, validationDTO *transfert.Validation, credentialDTO *transfert.Credential) (int, any) {
	if err := validationDTO.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	}); err != nil {
		return err.Code(), err
	}

	if err := credentialDTO.Check(data.Validator{
		"email":    {validator.Required, validator.Email},
		"password": {validator.Required, validator.Password},
	}); err != nil {
		return err.Code(), err
	}

	validation, err := service.PasswordValidation(validationDTO, &transfert.Credential{
		Email: credentialDTO.Email,
	})

	if err != nil {
		return err.Code(), err
	}

	if err := service.PasswordUpdate(credentialDTO); err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, validation
}

func MailValidation(service services.UserServiceInterface, dtoValidation *transfert.Validation, dtoCredential *transfert.Credential) (int, any) {
	if err := dtoValidation.Check(data.Validator{
		"token": {validator.Required, validator.Luhn},
	}); err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	if err := dtoCredential.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	}); err != nil {
		return fiber.StatusBadRequest, err.Error()
	}

	validation, err := service.MailValidation(dtoValidation, dtoCredential)
	if err != nil {
		return err.Code(), err
	}

	return fiber.StatusOK, validation

}

func ValidationRecover(service services.UserServiceInterface, dtoCredential *transfert.Credential, dtoValidation *transfert.Validation) (int, any) {
	if err := dtoCredential.Check(data.Validator{
		"email": {validator.Required, validator.Email},
	}); err != nil {
		return err.Code(), err
	}

	if err := dtoValidation.Check(data.Validator{
		"type": {validator.Required},
	}); err != nil {
		return err.Code(), err
	}

	if err := service.ValidationRecover(dtoValidation, dtoCredential); err != nil {
		return err.Code(), err
	}

	return fiber.StatusNoContent, nil
}

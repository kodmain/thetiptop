package services

import (
	"fmt"

	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
)

func (s *UserService) RegisterEmployee(dtoCredential *transfert.Credential, dtoEmployee *transfert.Employee) (*entities.Employee, error) {
	if dtoCredential == nil || dtoEmployee == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	_, err := s.repo.ReadCredential(dtoCredential)
	if err == nil {
		return nil, fmt.Errorf(errors.ErrEmployeeAlreadyExists)
	}

	credential, err := s.repo.CreateCredential(dtoCredential)
	if err != nil {
		return nil, err
	}

	employee, err := s.repo.CreateEmployee(dtoEmployee)
	if err != nil {
		return nil, err
	}

	employee.CredentialID = &credential.ID

	employee.Validations = append(employee.Validations, &entities.Validation{
		EmployeeID: &employee.ID,
		Type:       entities.MailValidation,
	})

	if err := s.repo.UpdateEmployee(employee); err != nil {
		return nil, err
	}

	if err := s.repo.UpdateCredential(credential); err != nil {
		return nil, err
	}

	go s.sendValidationMail(credential, employee.Validations[0])

	return employee, nil
}

func (s *UserService) GetEmployee(dtoEmployee *transfert.Employee) (*entities.Employee, error) {
	// Vérifier si le DTO de l'employé est valide
	if dtoEmployee == nil {
		return nil, fmt.Errorf(errors.ErrNoDto)
	}

	// Lire l'employé depuis le repository
	employee, err := s.repo.ReadEmployee(dtoEmployee)
	if err != nil {
		return nil, fmt.Errorf(errors.ErrEmployeeNotFound)
	}

	return employee, nil
}

func (s *UserService) DeleteEmployee(dtoEmployee *transfert.Employee) error {
	// Vérifier si le DTO de l'employé est valide
	if dtoEmployee == nil || dtoEmployee.ID == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	// Supprimer l'employé du repository
	if err := s.repo.DeleteEmployee(dtoEmployee); err != nil {
		return err
	}

	return nil
}

func (s *UserService) UpdateEmployee(dtoEmployee *transfert.Employee) error {
	// Lire l'employé depuis le repository pour vérifier s'il existe
	employee, err := s.repo.ReadEmployee(dtoEmployee)
	if err != nil {
		return fmt.Errorf(errors.ErrEmployeeNotFound)
	}

	// Mettre à jour l'employé dans le repository
	if err := s.repo.UpdateEmployee(employee); err != nil {
		return err
	}

	return nil
}

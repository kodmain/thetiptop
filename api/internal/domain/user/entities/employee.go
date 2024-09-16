package entities

import (
	"fmt"
	"reflect"
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/errors"
	"gorm.io/gorm"
)

type Employee struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	//Credential  *Credential `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CredentialID *string     `gorm:"type:varchar(36);index;" json:"credential_id"` // Foreign key to Credential
	Validations  Validations `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
}

func (employee *Employee) UpdateWith(dto *transfert.Employee) error {
	// Vérifie que le DTO n'est pas nul
	if dto == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	// Récupère les valeurs des structs employee et DTO
	employeeVal := reflect.ValueOf(employee).Elem()
	dtoVal := reflect.ValueOf(dto).Elem()

	// Parcourt les champs du DTO pour mettre à jour les champs correspondants dans le employee
	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Field(i)
		employeeField := employeeVal.FieldByName(dtoVal.Type().Field(i).Name)

		// Vérifie que le champ du employee existe et est assignable
		if employeeField.IsValid() && employeeField.CanSet() {
			// Si le champ du DTO est un pointeur et non nil
			if dtoField.Kind() == reflect.Ptr && !dtoField.IsNil() {
				// Si le champ correspondant dans le employee est aussi un pointeur, on assigne la valeur
				if employeeField.Kind() == reflect.Ptr {
					employeeField.Set(dtoField) // Assigner directement le pointeur
				} else {
					employeeField.Set(dtoField.Elem()) // Assigner la valeur pointée
				}
			} else if dtoField.Kind() != reflect.Ptr { // Si ce n'est pas un pointeur, on compare les valeurs directement
				if !reflect.DeepEqual(employeeField.Interface(), dtoField.Interface()) {
					employeeField.Set(dtoField) // Assigner la valeur si elle est différente
				}
			}
		}
	}

	return nil
}

func (employee *Employee) HasSuccessValidation(validationType ValidationType) *Validation {
	for _, validation := range employee.Validations {
		if validation.Type == validationType && validation.Validated {
			return validation
		}
	}

	return nil
}

func (employee *Employee) HasNotExpiredValidation(validationType ValidationType) *Validation {
	for i := len(employee.Validations) - 1; i >= 0; i-- {
		validation := employee.Validations[i]
		if validation.Type == validationType && !validation.HasExpired() && !validation.Validated {
			return validation
		}
	}

	return nil
}

func (employee *Employee) BeforeUpdate(tx *gorm.DB) error {
	employee.UpdatedAt = time.Now()
	return nil
}

func (employee *Employee) AfterFind(tx *gorm.DB) error {
	tx.Find(&employee.Validations, "employee_id = ?", employee.ID)
	//tx.Find(&employee.Credential, "employee_id = ?", employee.ID)
	return nil
}

func (employee *Employee) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	employee.ID = id.String()

	/*
		if employee.Credential != nil {
			employee.Credential.EmployeeID = &employee.ID
		}
	*/

	for _, validation := range employee.Validations {
		validation.EmployeeID = &employee.ID
	}

	return nil
}

func CreateEmployee(obj *transfert.Employee) *Employee {
	return &Employee{
		Validations: make(Validations, 0),
	}
}

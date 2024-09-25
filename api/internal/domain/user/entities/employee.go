package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
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
	CredentialID *string     `gorm:"type:varchar(36);index;" json:"-"` // Foreign key to Credential
	Validations  Validations `gorm:"foreignKey:EmployeeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
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

	for _, validation := range employee.Validations {
		validation.EmployeeID = &employee.ID
	}

	return nil
}

func (e *Employee) IsPublic() bool {
	return false
}

func (e *Employee) GetOwnerID() string {
	if e.CredentialID == nil {
		return ""
	}

	return *e.CredentialID
}

func CreateEmployee(obj *transfert.Employee) *Employee {
	e := &Employee{
		Validations:  make(Validations, 0),
		CredentialID: obj.CredentialID,
	}

	if obj.ID != nil {
		e.ID = *obj.ID
	}

	return e
}

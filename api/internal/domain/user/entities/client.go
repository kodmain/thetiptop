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

type Client struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	//Credential  *Credential `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CredentialID *string     `gorm:"type:varchar(36);index;" json:"-"` // Foreign key to Credential
	Validations  Validations `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	// Additional fields
	CGU        *bool `gorm:"type:boolean;default:false" json:"cgu"`
	Newsletter *bool `gorm:"type:boolean;default:false" json:"newsletter"`
}

func (client *Client) UpdateWith(dto *transfert.Client) error {
	// Vérifie que le DTO n'est pas nul
	if dto == nil {
		return fmt.Errorf(errors.ErrNoDto)
	}

	// Récupère les valeurs des structs Client et DTO
	clientVal := reflect.ValueOf(client).Elem()
	dtoVal := reflect.ValueOf(dto).Elem()

	// Parcourt les champs du DTO pour mettre à jour les champs correspondants dans le client
	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Field(i)
		clientField := clientVal.FieldByName(dtoVal.Type().Field(i).Name)

		// Vérifie que le champ du client existe et est assignable
		if clientField.IsValid() && clientField.CanSet() {
			// Si le champ du DTO est un pointeur et non nil
			if dtoField.Kind() == reflect.Ptr && !dtoField.IsNil() {
				// Si le champ correspondant dans le client est aussi un pointeur, on assigne la valeur
				if clientField.Kind() == reflect.Ptr {
					clientField.Set(dtoField) // Assigner directement le pointeur
				} else {
					clientField.Set(dtoField.Elem()) // Assigner la valeur pointée
				}
			} else if dtoField.Kind() != reflect.Ptr { // Si ce n'est pas un pointeur, on compare les valeurs directement
				if !reflect.DeepEqual(clientField.Interface(), dtoField.Interface()) {
					clientField.Set(dtoField) // Assigner la valeur si elle est différente
				}
			}
		}
	}

	return nil
}

func (client *Client) HasSuccessValidation(validationType ValidationType) *Validation {
	for _, validation := range client.Validations {
		if validation.Type == validationType && validation.Validated {
			return validation
		}
	}

	return nil
}

func (client *Client) HasNotExpiredValidation(validationType ValidationType) *Validation {
	for i := len(client.Validations) - 1; i >= 0; i-- {
		validation := client.Validations[i]
		if validation.Type == validationType && !validation.HasExpired() && !validation.Validated {
			return validation
		}
	}

	return nil
}

func (client *Client) BeforeUpdate(tx *gorm.DB) error {
	client.UpdatedAt = time.Now()
	return nil
}

func (client *Client) AfterFind(tx *gorm.DB) error {
	tx.Find(&client.Validations, "client_id = ?", client.ID)
	//tx.Find(&client.Credential, "client_id = ?", client.ID)
	return nil
}

func (client *Client) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	client.ID = id.String()

	/*
		if client.Credential != nil {
			client.Credential.ClientID = &client.ID
		}
	*/

	for _, validation := range client.Validations {
		validation.ClientID = &client.ID
	}

	return nil
}

func CreateClient(obj *transfert.Client) *Client {
	c := &Client{
		Validations: make(Validations, 0),
		CGU:         obj.CGU,
		Newsletter:  obj.Newsletter,
	}

	if obj.ID != nil {
		c.ID = *obj.ID
	}

	return c
}

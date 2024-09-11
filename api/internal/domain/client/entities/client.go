package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"gorm.io/gorm"
)

type Client struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// Relations
	Credential  *Credential `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Validations Validations `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`

	// Additional fields
	CGU        *bool `gorm:"type:boolean;default:false" json:"cgu"`
	Newsletter *bool `gorm:"type:boolean;default:false" json:"newsletter"`
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
	tx.Find(&client.Credential, "client_id = ?", client.ID)
	return nil
}

func (client *Client) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	client.ID = id.String()

	if client.Credential != nil {
		client.Credential.ClientID = &client.ID
	}

	for _, validation := range client.Validations {
		validation.ClientID = &client.ID
	}

	return nil
}

func CreateClient(obj *transfert.Client) *Client {
	return &Client{
		Validations: make(Validations, 0),
		CGU:         obj.CGU,
		Newsletter:  obj.Newsletter,
	}
}

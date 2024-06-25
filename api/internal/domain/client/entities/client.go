package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"gorm.io/gorm"
)

type Client struct {
	// gorm model
	ID        string         `gorm:"type:varchar(36);primaryKey;"`
	CreatedAt time.Time      `json:"-"`
	UpdatedAt time.Time      `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	// entity
	Email       string      `gorm:"type:varchar(320);uniqueIndex"`
	Password    string      `gorm:"type:varchar(255)" json:"-"` // private field
	Validations Validations `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	CGU         bool        `gorm:"type:boolean;default:false"`
	Newsletter  bool        `gorm:"type:boolean;default:false"`
}

func (client *Client) CompareHash(password string) bool {
	return hash.CompareHash(client.Password, client.Email+":"+password, hash.BCRYPT) == nil
}

func (client *Client) BeforeUpdate(tx *gorm.DB) error {
	logger.Info("BeforeUpdate")
	client.UpdatedAt = time.Now()
	return nil
}

func (client *Client) AfterFind(tx *gorm.DB) error {
	tx.Find(&client.Validations, "client_id = ?", client.ID)
	return nil
}

func (client *Client) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	password, err := hash.Hash(client.Email+":"+client.Password, hash.BCRYPT)
	if err != nil {
		return err
	}

	client.ID = id.String()
	client.Password = password

	for _, validation := range client.Validations {
		validation.ClientID = client.ID
	}

	return nil
}

func CreateClient(obj *transfert.Client) *Client {
	return &Client{
		Email:       obj.Email,
		Password:    obj.Password,
		Validations: []*Validation{},
	}
}

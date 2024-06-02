package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/hash"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	ID       string `gorm:"type:varchar(36);primaryKey;"`
	Email    string `gorm:"type:varchar(320);uniqueIndex"`
	Password string `gorm:"type:varchar(255)"` // private field

	ValidationEmail bool `gorm:"type:boolean;default:false"`
	CGU             bool `gorm:"type:boolean;default:false"`
	Newsletter      bool `gorm:"type:boolean;default:false"`
}

func (client *Client) CompareHash(password string) bool {
	return hash.CompareHash(client.Password, client.Email+":"+password, hash.BCRYPT) == nil
}

func (client *Client) BeforeUpdate(tx *gorm.DB) error {
	client.UpdatedAt = time.Now()
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

	return nil
}

func CreateClient(obj *transfert.Client) *Client {
	return &Client{
		Email:    obj.Email,
		Password: obj.Password,
	}
}

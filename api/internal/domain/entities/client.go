package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/application/dto"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security"
	"gorm.io/gorm"
)

/*
type Client struct {
	gorm.Model
	ID        uuid.UUID  `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email     string     `gorm:"type:varchar(100);uniqueIndex"`
	Password  string     `gorm:"size:255"`
	CreatedAt time.Time  `gorm:"type:timestamp"`
	UpdatedAt time.Time  `gorm:"type:timestamp"`
	DeletedAt *time.Time `gorm:"type:timestamp;index"`
}
*/

type Client struct {
	gorm.Model
	ID       uuid.UUID `gorm:"type:uuid;primaryKey;"`
	Email    string    `gorm:"type:varchar(100);uniqueIndex"`
	Password string    `gorm:"size:255"`
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

	client.ID = id
	return nil
}

func CreateClient(obj *dto.Client) (*Client, error) {
	password, err := security.Hash(obj.Email+":"+obj.Password, security.BCRYPT)

	if err != nil {
		return nil, nil
	}

	return &Client{
		ID:       uuid.New(),
		Email:    obj.Email,
		Password: password,
	}, nil
}

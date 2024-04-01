package entities

import (
	"time"

	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/architecture/events"
	"github.com/kodmain/thetiptop/api/internal/architecture/persistence"
	"gorm.io/gorm"
)

// Client defines the model for the 'clients' table.
// This model includes ID, Email, Password, CreatedAt, UpdatedAt, and DeletedAt.
//
// Parameters:
// - gorm.Model: Provides ID, CreatedAt, UpdatedAt, DeletedAt.
// - ID: uuid.UUID Unique identifier for the client, serves as the primary key.
// - Email: string Client's email, must be unique.
// - Password: string Client's password.
// - CreatedAt: time.Time Timestamp when the client record was created.
// - UpdatedAt: time.Time Timestamp when the client record was last updated.
// - DeletedAt: *time.Time Timestamp when the client record was deleted, can be null.
//
// Returns:
// - Client: Struct representing a client.

func init() {
	events.Subscribe(events.MIGRATE, func(a ...any) {
		persistence.Migrate(&Client{})
	})
}

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
	ID        uuid.UUID
	Email     string     `gorm:"type:varchar(100);uniqueIndex"`
	Password  string     `gorm:"size:255"`
	CreatedAt time.Time  `gorm:"type:timestamp"`
	UpdatedAt time.Time  `gorm:"type:timestamp"`
	DeletedAt *time.Time `gorm:"type:timestamp;index"`
}

func NewClient(email, password string) *Client {
	return &Client{
		ID:        uuid.New(),
		Email:     email,
		Password:  password,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

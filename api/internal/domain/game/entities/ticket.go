package entities

import (
	"time"

	"github.com/google/uuid"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/game"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/security/token"
	"gorm.io/gorm"
)

type Ticket struct {
	// Gorm model
	ID        string          `gorm:"type:varchar(36);primaryKey;" json:"id"`
	CreatedAt time.Time       `json:"-"`
	UpdatedAt time.Time       `json:"-"`
	DeletedAt *gorm.DeletedAt `gorm:"index" json:"-"`

	// Additional fields
	CredentialID *string    `gorm:"type:varchar(36);index" json:"credential_id"`
	Token        token.Luhn `gorm:"type:varchar(16);uniqueIndex" json:"token"`
	Prize        *string    `gorm:"type:varchar(36);index" json:"prize"`
}

func CreateTicket(obj *transfert.Ticket) *Ticket {
	t := &Ticket{
		CredentialID: obj.CredentialID,
		Prize:        obj.Prize,
		Token:        token.NewLuhnP(obj.Token),
	}

	if obj.ID != nil {
		t.ID = *obj.ID
	}

	return t
}

func (ticket *Ticket) IsPublic() bool {
	return false
}

func (ticket *Ticket) GetOwnerID() string {
	if ticket.CredentialID == nil {
		return ""
	}

	return *ticket.CredentialID
}

func (ticket *Ticket) BeforeUpdate(tx *gorm.DB) error {
	ticket.UpdatedAt = time.Now()
	return nil
}

func (ticket *Ticket) BeforeCreate(tx *gorm.DB) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return err
	}

	ticket.ID = id.String()

	return nil
}

package security

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

type PermissionInterface interface {
	IsGranted(roles ...string) bool
	GetCredentialID() *string
	CanRead(ressource database.Entity, rules ...Rule) bool
	CanCreate(ressource database.Entity, rules ...Rule) bool
	CanUpdate(ressource database.Entity, rules ...Rule) bool
	CanDelete(ressource database.Entity, rules ...Rule) bool
}

type UserAccess struct {
	CredentialID string
	Role         string
}

type Rule func(p *UserAccess, entity database.Entity) bool

func (p *UserAccess) GetCredentialID() *string {
	if p.CredentialID == "" {
		return nil
	}

	return &p.CredentialID
}

func (p *UserAccess) IsGranted(roles ...string) bool {
	for _, role := range roles {
		if p.Role == role {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanRead(ressource database.Entity, rules ...Rule) bool {
	if p.CredentialID == ressource.GetOwnerID() && p.CredentialID != "" {
		return true
	}

	if ressource.IsPublic() {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanCreate(ressource database.Entity, rules ...Rule) bool {
	if p.CredentialID == ressource.GetOwnerID() && p.CredentialID != "" {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanUpdate(ressource database.Entity, rules ...Rule) bool {
	if p.CredentialID == ressource.GetOwnerID() && p.CredentialID != "" {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanDelete(ressource database.Entity, rules ...Rule) bool {
	if p.CredentialID == ressource.GetOwnerID() && p.CredentialID != "" {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func NewUserAccess(token any) *UserAccess {
	p := &UserAccess{}
	if token != nil {
		if token, ok := token.(*jwt.Token); ok {
			p.CredentialID = token.ID
			if role, exists := token.Data["role"]; exists {
				if roleStr, ok := role.(string); ok {
					p.Role = roleStr
				}
			}
		}
	}

	return p
}

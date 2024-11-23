package security

import (
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

type PermissionInterface interface {
	IsAuthenticated() bool
	IsGrantedByRoles(roles ...Role) bool
	IsGrantedByRules(rules ...Rule) bool
	GetCredentialID() *string
	CanRead(ressource database.Entity, rules ...Rule) bool
	CanCreate(ressource database.Entity, rules ...Rule) bool
	CanUpdate(ressource database.Entity, rules ...Rule) bool
	CanDelete(ressource database.Entity, rules ...Rule) bool
}

type UserAccess struct {
	CredentialID string
	Role         Role
}

type Role string
type Rule func(p *UserAccess) bool

const (
	ROLE_ADMIN     Role = "admin"
	ROLE_ANONYMOUS Role = "anonymous"
	ROLE_CONNECTED Role = "connected"
)

func (p *UserAccess) IsAuthenticated() bool {
	return p.CredentialID != ""
}

func (p *UserAccess) GetCredentialID() *string {
	if p.CredentialID == "" {
		return nil
	}

	return &p.CredentialID
}

func (p *UserAccess) IsGrantedByRules(rules ...Rule) bool {
	for _, rule := range rules {
		if rule(p) {
			return true
		}
	}

	return false
}

func (p *UserAccess) IsGrantedByRoles(roles ...Role) bool {
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
		if rule(p) {
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
		if rule(p) {
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
		if rule(p) {
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
		if rule(p) {
			return true
		}
	}

	return false
}

func NewUserAccess(token any) *UserAccess {
	p := &UserAccess{
		Role: ROLE_ANONYMOUS,
	}
	if token != nil {
		if token, ok := token.(*jwt.Token); ok {
			p.CredentialID = token.ID
			if role, exists := token.Data["role"]; exists {
				if roleStr, ok := role.(string); ok {
					p.Role = Role(roleStr)
				}
			}
		}
	}

	return p
}

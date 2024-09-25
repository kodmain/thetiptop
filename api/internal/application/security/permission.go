package security

import (
	"github.com/kodmain/thetiptop/api/internal/domain/user/entities"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/serializers/jwt"
)

type PermissionInterface interface {
	IsGranted(roles ...string) bool
	CanRead(ressource entities.Entity, rules ...Rule) bool
	CanCreate(ressource entities.Entity, rules ...Rule) bool
	CanUpdate(ressource entities.Entity, rules ...Rule) bool
	CanDelete(ressource entities.Entity, rules ...Rule) bool
}

type UserAccess struct {
	CredentialID string
	Role         []string
}

type Rule func(p *UserAccess, r entities.Entity) bool

func (p *UserAccess) IsGranted(roles ...string) bool {
	for _, role := range roles {
		for _, userRole := range p.Role {
			if userRole == role {
				return true
			}
		}
	}
	return false
}

func (p *UserAccess) CanRead(ressource entities.Entity, rules ...Rule) bool {
	if p.IsGranted("admin") || ressource.IsPublic() {
		return true
	}

	if p.CredentialID == ressource.GetOwnerID() {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanCreate(ressource entities.Entity, rules ...Rule) bool {
	if p.IsGranted("admin") || p.CredentialID == ressource.GetOwnerID() {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanUpdate(ressource entities.Entity, rules ...Rule) bool {
	if p.IsGranted("admin") || p.CredentialID == ressource.GetOwnerID() {
		return true
	}

	for _, rule := range rules {
		if rule(p, ressource) {
			return true
		}
	}

	return false
}

func (p *UserAccess) CanDelete(ressource entities.Entity, rules ...Rule) bool {
	if p.IsGranted("admin") || p.CredentialID == ressource.GetOwnerID() {
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
		}
	}

	return p
}